package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	models "github.com/mannx/weather/models"
	"github.com/rs/zerolog/log"
	//models "github.com/mannx/weather/models"
)

// WeatherChartData contains selected items used for charting some data
type WeatherChartData struct {
	DateUnix   int64
	TimeString string // just the time, no date portion

	Temp      float64
	FeelsLike float64
	WindSpeed float64
	Rain      float64
	Snow      float64
}

// WeatherDataView is used to return the weather data along with several
// combined values for the given range
type WeatherDataView struct {
	Data []models.WeatherData

	Low  float64
	High float64
	Snow float64
	Rain float64

	ChartData []WeatherChartData
}

func handle24hrView(c echo.Context) error {
	// retrieve all entries in the last 24 hours
	now := time.Now()
	prev := now.Add(-time.Hour * 24)

	// retrieve the data
	var wd []models.WeatherData

	res := DB.Find(&wd, "store_time BETWEEN ? AND ?", prev.Unix(), now.Unix())
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("Unable to retrieve data")
		return res.Error
	}

	view := computeWeatherDataView(wd)
	return c.JSON(http.StatusOK, &view)
}

func computeWeatherDataView(data []models.WeatherData) WeatherDataView {
	v := WeatherDataView{Data: data}
	v.ChartData = make([]WeatherChartData, 0)

	// compute the highest
	f := false

	for _, e := range data {
		if f == false {
			// init the high temp (otherwise high's less than 0 dont register)
			v.High = e.Temp
			f = true
		}

		if e.Temp > v.High {
			v.High = e.Temp
		}
		if e.Temp < v.Low {
			v.Low = e.Temp
		}

		v.Snow += e.Snow1h
		v.Rain += e.Rain1h

		cd := WeatherChartData{Temp: e.Temp, FeelsLike: e.FeelsLike, Snow: e.Snow1h, Rain: e.Rain1h, WindSpeed: e.WindSpeed}
		cd.DateUnix = e.StoreTime
		cd.TimeString = time.Unix(e.StoreTime, 0).Format(time.Kitchen)

		v.ChartData = append(v.ChartData, cd)
	}

	return v
}
