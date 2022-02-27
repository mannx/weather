package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
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
const Version = 0.06

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

	// load configuration file
	err := loadConfiguration(&Config)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse configuration files, exiting...")
	}
	log.Debug().Msgf("APIKey: %v", Config.APIKey)

	dbPath = filepath.Join(Environment.DataPath, "db.db")

	log.Info().Msg("Initializing database...")
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database!")
	}

	log.Info().Msg("Migrating database...")
	migrateDB()

	// uncomment run and rebuild, or use a flag to call if required (or web option?)
	/*log.Debug().Msg("Adding city data to database...")
	updateCityIDdb()*/

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
	log.Debug().Msgf("updateWeatherFunc(%v)", len(Config.CityIDs))

	for _, i := range Config.CityIDs {
		log.Debug().Msgf("   => Logging weather for city: %v", i)
		err := getCurrentWeather(i)
		if err != nil {
			log.Error().Err(err).Msg("Unable to retrieve weather")
		} else {
			log.Info().Msg("Weather updated successfully")
		}
	}
}

func migrateDB() {
	DB.AutoMigrate(&models.WeatherData{})
}

// temp function to add the test city id to the current database
func updateCityIDdb() {
	if len(Config.CityIDs) > 1 {
		log.Error().Msg("Too many city ids found")
		return
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
