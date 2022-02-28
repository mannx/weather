package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	models "github.com/mannx/weather/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func GetDailyWeatherView(c echo.Context, db *gorm.DB) error {
	// we have the current day to view in parameters, along with the city id
	var month, year, day int
	var city int

	err := echo.QueryParamsBinder(c).
		Int("month", &month).
		Int("year", &year).
		Int("day", &day).
		Int("city", &city).
		BindError()

	if err != nil {
		log.Error().Err(err).Msg("[/api/daily] Unable to bind parameters")

		sr := models.ServerResponse{
			Message: "Unable to bind parameters",
			Error:   true,
		}
		return c.JSON(http.StatusOK, &sr)
	}

	if city == 0 {
		// invalid city name, return an empty result
		return c.JSON(http.StatusOK, &models.ServerResponse{Message: "Select City to view Data for", Error: true})
	}

	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, time.Month(month), day+1, 0, 0, 0, 0, time.UTC)

	log.Debug().Msgf("[/api/daily] Start: %v", start)
	log.Debug().Msgf("[/api/daily] End: %v", end)

	var data []models.WeatherData

	res := db.Where("city_id = ?", city).Where("store_time BETWEEN ? AND ?", start.Unix(), end.Unix()).Find(&data)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("Unable to retrieve data")
		return c.JSON(http.StatusOK, &models.ServerResponse{Message: "unable to get data", Error: true})
	}

	log.Debug().Msgf("Retrieved %v weather records...", res.RowsAffected)
	if res.RowsAffected == 0 {
		log.Warn().Msgf("0 weather records retrieved for city %v", city)
		return c.JSON(http.StatusOK, &models.ServerResponse{Message: "No records available", Error: true})
	}

	// return structure containing the stats we want to display
	type dailyStats struct {
		MinWindSpeed float64
		MaxWindSpeed float64
		AverageWind  float64
		MaxWindGust  float64

		MinTemp     float64
		MaxTemp     float64
		AverageTemp float64

		AverageRain float64
		AverageSnow float64
	}

	stats := dailyStats{}

	// go through data and compute our stats
	for _, wd := range data {
		// make sure mins start with a possible value != 0
		// test to see if this works as intended
		if stats.MinWindSpeed == 0 {
			stats.MinWindSpeed = wd.WindSpeed
		}
		if stats.MinTemp == 0 {
			stats.MinTemp = wd.Temp
		}
		if stats.MaxTemp == 0 {
			stats.MaxTemp = wd.Temp
		}

		if stats.MinWindSpeed > wd.WindSpeed {
			stats.MinWindSpeed = wd.WindSpeed
		}

		if stats.MaxWindSpeed < wd.WindSpeed {
			stats.MaxWindSpeed = wd.WindSpeed
		}

		if stats.MaxWindGust < wd.WindGust {
			stats.MaxWindGust = wd.WindGust
		}

		stats.AverageWind += wd.WindSpeed

		if stats.MinTemp > wd.Temp {
			stats.MinTemp = wd.Temp
		}

		if stats.MaxTemp < wd.Temp {
			stats.MaxTemp = wd.Temp
		}

		stats.AverageTemp += wd.Temp
		stats.AverageRain += wd.Rain1h
		stats.AverageSnow += wd.Snow1h
	}

	stats.AverageWind = stats.AverageWind / float64(len(data))
	stats.AverageRain = stats.AverageRain / float64(len(data))
	stats.AverageSnow = stats.AverageSnow / float64(len(data))
	stats.AverageTemp = stats.AverageTemp / float64(len(data))

	log.Debug().Msg("Returning daily data...")
	return c.JSON(http.StatusOK, &stats)
}
