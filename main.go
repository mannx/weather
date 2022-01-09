package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "modernc.org/sqlite"

	"github.com/robfig/cron/v3"
)

// DBPath points to the database file. /app/db.db in container, ./data/db.db while developing
const DBPath = "./data/db.db"

// id for the city we are interested in.
//	TODO: havea configuration file/page to set this somehow
const cityID = 6138517

// Version of the current build
const Version = 0.02

// used to the store the OpenWeatherMap api key.
// passed to us through the APIKEY environment variable
var apiKey string

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // can chagne to zerolog.DebugLevel for more info, or ErrorLevel for just errors

	log.Info().Msgf("Weather version: %v", Version)

	// get required data from environment variables
	apiKey = os.Getenv("APIKEY")
	if apiKey == "" {
		log.Fatal().Msg("APIKEY environment variable not defined.")
	}

	log.Info().Msg("Creating database if required...")

	// make sure the database is created before we start trying to use it
	// open the db
	/*db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to open database")
	}

	defer db.Close()*/

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
	e.Logger.Fatal(e.Start(":8080"))
}

// updateWeatherFunc is called to retrieve and store the current weather
// this should only be called from a scheduled job
func updateWeatherFunc() {
	err := getCurrentWeather()
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve weather")
	} else {
		log.Info().Msg("Weather updated successfully")
	}
}
