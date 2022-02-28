package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	models "github.com/mannx/weather/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func GetCityList(c echo.Context, db *gorm.DB, cfg *models.Configuration) error {
	// we use a new struct to avoid sending unnessecary information back
	type returnData struct {
		ID   int
		Name string
	}

	// retrieve the name of each city in the list of cities
	// check the CityData table for a cached name, otherwise use the json data
	var data []returnData

	for _, n := range cfg.CityIDs {
		rd := returnData{ID: n}

		// do we have a cached name?
		var cd models.CityData
		res := db.Find(&cd, "city_id = ?", n)
		if res.Error != nil {
			log.Error().Err(res.Error).Msg("Unable to lookup city name")
			return serverError(c, fmt.Sprintf("Unable to lookup city id: %v", n))
		}
		if res.RowsAffected == 0 {
			// not in db yet, find in json data
			rd.Name = lookupCity(n, db)
		} else {
			rd.Name = cd.Name
		}

		data = append(data, rd)
	}

	return c.JSON(http.StatusOK, &data)
}

func serverError(c echo.Context, msg string) error {
	sr := models.ServerResponse{
		Error:   true,
		Message: msg,
	}

	return c.JSON(http.StatusOK, &sr)
}

// lookup a city by id from the json data and return the string and cache in the db
func lookupCity(id int, db *gorm.DB) string {
	type jsonData struct {
		ID   float64 `json: "id"` // several entries havea trailing .0 to the id for some reason
		Name string  `json: "name"`
	}

	file, err := ioutil.ReadFile("./city.json")
	if err != nil {
		log.Error().Err(err).Msg("Unable to read city.json to get names")
		return "--NO CITY NAME LIST FOUND--"
	}

	var data []jsonData
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse city.json to get city names")
		return "--NO CITY NAME LIST FOUND--"
	}

	// find the name
	for _, i := range data {
		if int(i.ID) == id {
			// found it, save it db and return
			cd := models.CityData{CityID: id, Name: i.Name}
			go db.Save(&cd)
			return i.Name
		}
	}

	// not found
	log.Warn().Msgf("Unable to find city id [%v] in city.json", id)
	return "--INVALID CITY ID--"
}
