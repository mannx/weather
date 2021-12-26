package main

import (
	"database/sql"
	"log"
)

type WeatherData struct {
	time      float64
	code      float64
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
	sunrise float64
	sunset  float64
	name    string // name of the city we are storing data for
}

func GetWeatherData(db *sql.DB) WeatherData {
	s := "SELECT temp, feelsLike, humidity, windSpeed, windDir, rain1h, snow1h, icon FROM weather ORDER BY id DESC LIMIT 1"
	r, err := db.Query(s)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Close()

	wd := WeatherData{}

	for r.Next() {
		if err := r.Scan(&wd.Temp, &wd.FeelsLike, &wd.Humidity, &wd.WindSpeed, &wd.WindDir, &wd.Rain1h, &wd.Snow1h, &wd.Icon); err != nil {
			log.Fatal(err)
		}
	}

	return wd
}
