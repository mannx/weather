package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	models "github.com/mannx/weather/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

//func GetCityList(c echo.Context, db *gorm.DB, cfg *models.Configuration) error {
func GetCityList(c echo.Context, db *gorm.DB) error {
	// we use a new struct to avoid sending unnessecary information back
	type returnData struct {
		ID   int
		Name string
	}

	// retrieve the name of each city in the list of cities
	// check the CityData table for a cached name, otherwise use the json data
	var data []returnData

	var ce []models.ConfigEntry
	res := db.Find(&ce)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("[GetCityList] Unable to retrieve config table entries...")
		return serverError(c, "Unable to retrieve city list...")
	}

	for _, n := range ce {
		if n.Active {
			rd := returnData{ID: n.CityID, Name: n.Name}
			data = append(data, rd)
		}
	}

	return c.JSON(http.StatusOK, &data)
}

// JSONCityData contains the data in the format from openweathermap.org, city.list.json
type JSONCityData struct {
	ID      float64 `json: "id"` // several entries havea trailing .0 to the id for some reason
	Name    string  `json: "name"`
	State   string  `json: "state"`
	Country string  `json: "country"` // 2 letter country code
	Coords  struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
}

// lookup a city by id from the json data and return the string and cache in the db
func lookupCity(id int) string {

	file, err := ioutil.ReadFile("./city.json")
	if err != nil {
		log.Error().Err(err).Msg("Unable to read city.json to get names")
		return "--NO CITY NAME LIST FOUND--"
	}

	var data []JSONCityData
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse city.json to get city names")
		return "--NO CITY NAME LIST FOUND--"
	}

	// find the name
	for _, i := range data {
		if int(i.ID) == id {
			// found it, save it db and return
			return i.Name
		}
	}

	// not found
	log.Warn().Msgf("Unable to find city id [%v] in city.json", id)
	return "--INVALID CITY ID--"
}

// Find a city given its name and return the first entry for the city, returns an empty structure on error, ignores case
// return false if unable to find
func findCityByName(name string, out *[]JSONCityData) bool {
	file, err := ioutil.ReadFile("./city.json")
	if err != nil {
		log.Error().Err(err).Msg("[findCityByName] Unable to read city.json to find city")
		return false
	}

	var data []JSONCityData
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Error().Err(err).Msg("[findCityByName] Unable to unmarshal city list")
		return false
	}

	// find the first city that matches the name, ignore case
	for _, d := range data {
		if strings.EqualFold(name, d.Name) {
			// found a match
			log.Debug().Msgf("Lat: %v, Lon: %v", d.Coords.Lat, d.Coords.Lon)
			*out = append(*out, d)
		}
	}

	//return false
	return len(*out) > 0
}

// ValidateCity takes a city name and returns information to confirm correct city
func ValidateCity(c echo.Context, db *gorm.DB) error {
	// we are POST'd the city name in the parameter: 'city'
	// return country code, and lat/lon coordinates, and id found in hte json file
	type postData struct {
		City string `json:"city" form:"city" query:"city"`
	}

	data := new(postData)
	if err := c.Bind(&data); err != nil {
		log.Error().Err(err).Msg("[ValidateCity] Unable to bind parameters")
		return serverError(c, "Server Error -- Unable to bind parameters")
	}

	// try and locate the city information
	var city []JSONCityData
	rval := findCityByName(data.City, &city)
	if rval == false {
		// unable to location
		return serverError(c, fmt.Sprintf("Unable to find city named: %v", data.City))
	}

	return c.JSON(http.StatusOK, &city)
}

// AddCity is called after validate to add the city to the config database table
func AddCity(c echo.Context, db *gorm.DB) error {
	// city id to add is parameter: 'id'
	// parse as string for reasons
	type inbound struct {
		ID string `param: "id" query:"id" form:"id"`
	}

	inb := inbound{}
	err := c.Bind(&inb)
	if err != nil {
		log.Error().Err(err).Msg("[AddCity] Unable to bind parameters")
		return serverError(c, "Unable to retrieve city id")
	}

	city, _ := strconv.Atoi(inb.ID)

	// get the city name
	name := lookupCity(city)
	ce := models.ConfigEntry{CityID: city, Name: name, Active: true}

	log.Debug().Msg("Saving new city, and returning current list...")
	db.Save(&ce)

	return GetCityList(c, db)
}

// ConfigToDatabaseHandler is used to convert an existing configuration file into the database format we are now using
func ConfigToDatabaseHandler(c echo.Context, db *gorm.DB, cfg *models.Configuration) error {
	var ce []models.ConfigEntry

	res := db.Find(&ce)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("[ConfigToDBHandler] Unable to retrieve current database config")
		return serverError(c, "Unable to retrieve current database config")
	}

	// add the current config ids to the database but only if not already there
	log.Debug().Msg("Looping through config id's...")
	for _, id := range cfg.CityIDs {
		log.Debug().Msgf("Checking existence of id %v...", id)
		exists := false

		for _, e := range ce {
			if e.CityID == id {
				log.Debug().Msgf("%v already in config database..skipping", id)
				exists = true
				break
			}
		}

		if !exists {
			name := lookupCity(id)
			entry := models.ConfigEntry{CityID: id, Name: name, Active: true}
			log.Debug().Msgf("Adding id %v [%v] to database...", id, name)
			go db.Save(&entry)
		}
	}

	log.Debug().Msg("Config migration complete")
	return c.JSON(http.StatusOK, &models.ServerResponse{Error: false, Message: "success"})
}
