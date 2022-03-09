package main

import (
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type EnvironmentDefinition struct {
	DataPath string `envconfig:"WEATHER_DATA_PATH"`
	Schedule string `envconfig:"WEATHER_UPDATE_SCHEDULE"`
	ApiKey   string `envconfig:"WEATHER_API_KEY"`
	LogLevel string `envconfig:"WEATHER_LOG_LEVEL"`
}

var Environment EnvironmentDefinition

func (e *EnvironmentDefinition) Init() {
	e.Default()

	err := envconfig.Process("", e)
	if err != nil {
		log.Error().Err(err).Msg("Unable to process environment variables")
	}
}

func (e *EnvironmentDefinition) Default() {
	e.DataPath = "/data"
	e.Schedule = "@hourly"
	e.ApiKey = "--INVALID API KEY--"
	e.LogLevel = "INFO"
}

// Path returns a path to a file in the data path
func (e *EnvironmentDefinition) Path(file string) string {
	return filepath.Join(e.DataPath, file)
}
