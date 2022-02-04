package models

// WeatherData structure holds data for a given weather pull
type WeatherData struct {
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
