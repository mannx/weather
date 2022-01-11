package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "modernc.org/sqlite"

	"github.com/robfig/cron/v3"
)

// DBPath points to the database file. /app/db.db in container, ./data/db.db while developing
const DBPath = "./data/db.db"

// Version of the current build
const Version = 0.04

// Config -> Global configuration that is loaded for this instance
var Config Configuration

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // can chagne to zerolog.DebugLevel for more info, or ErrorLevel for just errors

	log.Info().Msgf("Weather version: %v", Version)

	// load configuration file
	Config, err := loadConfiguration()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse configuration files, exiting...")
	}
	log.Debug().Msgf("APIKey: %v", Config.APIKey)

	log.Info().Msg("Creating database if required...")

	// make sure the database is created before we start trying to use it
	// open the db
	db, err := openDB()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open/create database")
	}
	defer db.Close()

	log.Info().Msg("Creating tables if required...")

	err = createDB(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create required tables")
	}

	log.Info().Msg("Initializing echo and middleware")
	e := initServer()

	log.Info().Msg("Setting up cron job for updates")

	c := cron.New()
	expr := os.Getenv("WEATHER_UPDATE_SCHEDULE")
	if expr == "" {
		// update not set, default to hourly
		log.Info().Msg("WEATHER_UPDATE_SCHEDULE not set, defaulting to @hourly")
		expr = "@hourly"
	}

	c.AddFunc(expr, updateWeatherFunc)
	c.Start() // make sure to start the jobs

	log.Info().Msg("Starting server...")

	// start the server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down server")
		}
	}()

	// wait for interrupt signal to shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// updateWeatherFunc is called to retrieve and store the current weather
// this should only be called from a scheduled job
func updateWeatherFunc() {
	for _, i := range Config.CityIDs {
		err := getCurrentWeather(i)
		if err != nil {
			log.Error().Err(err).Msg("Unable to retrieve weather")
		} else {
			log.Info().Msg("Weather updated successfully")
		}
	}
}
