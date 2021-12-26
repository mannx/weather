package main

type WeatherData struct {
	time      float64
	code      float64
	temp      float64
	feelsLike float64
	pressure  float64
	humidity  float64
	windSpeed float64
	windDir   float64
	windGust  float64
	rain1h    float64
	rain3h    float64
	snow1h    float64
	snow3h    float64

	icon    string // icon name for the current weather
	sunrise float64
	sunset  float64
	name    string // name of the city we are storing data for
}
