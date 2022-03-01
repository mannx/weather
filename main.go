package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite"

	models "github.com/mannx/weather/models"
)

// DBPath points to the database file.
var dbPath string

// Version of the current build
const Version = 0.07

// Config -> Global configuration that is loaded for this instance
var Config models.Configuration

// DB -> Global connection to the database
var DB *gorm.DB

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // can chagne to zerolog.DebugLevel for more info, or ErrorLevel for just errors

	log.Info().Msgf("Weather version: %v", Version)

	log.Info().Msg("Initializing environment...")
	Environment.Init()

	updateLogLevel() // sets the log level based on environment variable, defaults to INFO

	// load configuration file
	err := loadConfiguration(&Config)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse configuration files, exiting...")
	}
	log.Debug().Msgf("APIKey: %v", Config.APIKey)

	dbPath = Environment.Path("db.db")

	log.Info().Msg("Initializing database...")
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database!")
	}

	log.Info().Msg("Migrating database...")
	migrateDB()

	// this is used to initialize a already in use database.
	// once all db's that are in use have been successully updated, this code can be removed for good
	// since this should only affect early versions
	cityadd := flag.Bool("city", false, "add first city id to all current entries")
	flag.Parse()
	if *cityadd {
		updateCityIDdb()
	}
	// END CODE REMOVAL SECTION

	log.Info().Msg("Initializing echo and middleware")
	e := initServer()

	log.Info().Msg("Setting up cron job for updates")

	c := cron.New()
	c.AddFunc(Environment.Schedule, updateWeatherFunc)
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
	// instead of config file, we use the ConfigEntry table for cities
	// only use entries noted as Active
	var data []models.ConfigEntry
	res := DB.Find(&data)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("Unable to retrieve cites for weather update!")
		return
	}

	for _, i := range data {
		log.Debug().Msgf("   => Logging weather for city: %v [%v]", i.Name, i.CityID)
		err := getCurrentWeather(i.CityID)
		if err != nil {
			log.Error().Err(err).Msg("Unable to retrieve weather")
		} else {
			log.Info().Msg("Weather updated successfully")
		}
	}
}

func migrateDB() {
	DB.AutoMigrate(&models.WeatherData{})
	DB.AutoMigrate(&models.CityData{})
	DB.AutoMigrate(&models.ConfigEntry{})
}

func updateLogLevel() {
	lookup := map[string]zerolog.Level{
		"DEBUG":    zerolog.DebugLevel,
		"INFO":     zerolog.InfoLevel,
		"WARN":     zerolog.WarnLevel,
		"ERROR":    zerolog.ErrorLevel,
		"FATAL":    zerolog.FatalLevel,
		"PANIC":    zerolog.PanicLevel,
		"NOLEVEL":  zerolog.NoLevel,
		"DISABLED": zerolog.Disabled,
	}

	key := strings.ToUpper(Environment.LogLevel)
	log.Debug().Msgf("Setting log level to: %v [%v]", key, lookup[key])
	zerolog.SetGlobalLevel(lookup[key])
}

// START CODE REMOVAL SECTION
// temp function to add the test city id to the current database
// once all early version db's have been updated, this can get removed
func updateCityIDdb() {
	if len(Config.CityIDs) > 1 {
		log.Error().Msg("Too many city ids found, using first in list")
	}

	var data []models.WeatherData
	cityID := Config.CityIDs[0]

	res := DB.Find(&data)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("Unable to retrieve data for update")
		return
	}

	log.Debug().Msg("Updating city id's in database...")
	for _, d := range data {
		d.CityID = cityID
		DB.Save(&d)
	}
	log.Debug().Msg("Update finished")
}

// END CODE REMOVAL SECTION
