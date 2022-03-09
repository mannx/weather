package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	models "github.com/mannx/weather/models"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

//
// file contains all config related functions
//

// Path to the config file
const configFileName = "config.yml"

// defaultConfiguration returns a default Config structure
// in case no config file was found
func defaultConfiguration(cfg *models.Configuration) {
	log.Info().Msg("Generating default Configuration file, No valid API key set")

	*cfg = models.Configuration{CityIDs: make([]int, 0), APIKey: "--INVALID API KEY--"}
}

func loadConfiguration(cfg *models.Configuration) error {
	log.Info().Msg("Preparing to load configuration file")

	cfn := Environment.Path(configFileName) // use the supplied path
	f, err := os.Open(cfn)
	if err != nil {
		// unable to open the file, use a default config
		log.Error().Err(err).Msg("unable to read config file, using defaults")
		defaultConfiguration(cfg)
	} else {
		defer f.Close()

		decode := yaml.NewDecoder(f)
		err = decode.Decode(&cfg)
		if err != nil {
			// unable to decode, use default
			log.Error().Err(err).Msg("unable to parse config file, using defaults")
			defaultConfiguration(cfg)
		}
	}

	// read in any enviroment udpates
	err = envconfig.Process("", cfg)
	if err != nil {
		log.Error().Err(err).Msg("Unable to process enviroment")
		return err
	}

	for i, n := range cfg.CityIDs {
		log.Debug().Msgf("[%v] %v", i, n)
	}

	return nil
}
