package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type EnvironmentDefinition struct {
	DataPath string `envconfig:"WEATHER_DATA_PATH"`
	Schedule string `envconfig:"WEATHER_UPDATE_SCHEDULE"`
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
}
