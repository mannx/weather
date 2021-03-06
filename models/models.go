package models

import "gorm.io/gorm"

// WeatherData structure holds data for a given weather pull
type WeatherData struct {
	gorm.Model

	CityID    int // id for the city that this data is for
	Code      float64
	Temp      float64
	FeelsLike float64
	Pressure  float64
	Humidity  float64
	WindSpeed float64
	WindDir   float64
	WindGust  float64
	Rain1h    float64
	Rain3h    float64
	Snow1h    float64
	Snow3h    float64

	Icon    string // icon name for the current weather
	Sunrise float64
	Sunset  float64
	Name    string // name of the city we are storing data for

	TimeString string // used only before send off to the template engine to display a user friendly Time string
	StoreTime  int64  // system time when this entry was added to the db
}

// ServerResponse is used to hold an error response to send back
type ServerResponse struct {
	Message string
	Error   bool
}

// Configuration stores configuration data set by the user in a config file
type Configuration struct {
	CityIDs       []int  `yaml:"CityIDs"`
	APIKey        string `yaml:"APIKey", envconfig:"WEATHER_API_KEY"`
	WeatherUpdate string `yaml:"WeatherUpdate", envconfig:"WEATHER_UPDATE_SCHEDULE"`
}

// CityData holds a city ID and name for quick look up instead of parseing the json every time
type CityData struct {
	gorm.Model

	CityID int
	Name   string
}
