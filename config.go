package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

//
// file contains all config related functions
//

// Path to the config file
const configFileName = "./data/config.yml"

// Configuration structure to old various bits of config data
type Configuration struct {
	CityIDs       []int  `yaml:"CityIDs"`
	APIKey        string `yaml:"APIKey", envconfig:"APIKEY"`
	WeatherUpdate string `yaml:"WeatherUpdate", envconfig:"WEATHER_UPDATE_SCHEDULE"`
}

// defaultConfiguration returns a default Config structure
// in case no config file was found
func defaultConfiguration(cfg *Configuration) {
	log.Info().Msg("Generating default Configuration file, No valid API key set")

	*cfg = Configuration{CityIDs: make([]int, 0), APIKey: "--INVALID API KEY--"}
}

func loadConfiguration() (Configuration, error) {
	log.Info().Msg("Preparing to load configuration file")

	var cfg Configuration

	f, err := os.Open(configFileName)
	if err != nil {
		// unable to open the file, use a default config
		log.Error().Err(err).Msg("unable to read config file, using defaults")
		defaultConfiguration(&cfg)
	} else {
		defer f.Close()

		decode := yaml.NewDecoder(f)
		err = decode.Decode(&cfg)
		if err != nil {
			// unable to decode, use default
			log.Error().Err(err).Msg("unable to parse config file, using defaults")
			defaultConfiguration(&cfg)
		}
	}

	// read in any enviroment udpates
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Error().Err(err).Msg("Unable to process enviroment")
		return Configuration{}, err
	}

	for i, n := range cfg.CityIDs {
		log.Debug().Msgf("[%v] %v", i, n)
	}

	return cfg, nil
}

func saveConfiguration(cfg *Configuration) error {
	log.Info().Msg("Preparing to save configuration")

	f, err := os.OpenFile(configFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Error().Err(err).Msg("Unable to open config file")
		return err
	}

	defer f.Close()

	encode := yaml.NewEncoder(f)
	err = encode.Encode(&Config)
	log.Info().Msg("Configuration Encoding complete")
	return err
}
