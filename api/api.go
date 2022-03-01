package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	models "github.com/mannx/weather/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
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

func Handle24hrView(c echo.Context, db *gorm.DB) error {
	// we have the city id provided by the 'id' parameter
	var cityid int

	err := echo.QueryParamsBinder(c).
		Int("id", &cityid).
		BindError()
	if err != nil {
		log.Error().Err(err).Msg("Unable to bind paramter: id [Handle24hrView]")
		return serverError(c, "Missing paramter 'id'")
	}

	// retrieve all entries in the last 24 hours
	now := time.Now()
	prev := now.Add(-time.Hour * 24)

	// retrieve the data
	var wd []models.WeatherData

	//res := db.Find(&wd, "store_time BETWEEN ? AND ?", prev.Unix(), now.Unix())
	res := db.Where("city_id = ?", cityid).Where("store_time BETWEEN ? AND ?", prev.Unix(), now.Unix()).Find(&wd)
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

func GetLatestWeatherView(c echo.Context, db *gorm.DB) error {
	var wd models.WeatherData

	res := db.Last(&wd)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("Unable to retrieve latest weather report")
		return res.Error
	}

	if res.RowsAffected <= 0 {
		log.Warn().Msg("Unable to retrieve latest weather...No records available")
		return c.JSON(http.StatusOK, &models.ServerResponse{Error: true, Message: "No Data Available"})
	}

	return c.JSON(http.StatusOK, &wd)
}
