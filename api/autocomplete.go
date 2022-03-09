package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// AutoCompleteCity returns the list of all possible cities we can gather data from
func AutoCompleteCity(c echo.Context) error {
	type returnData struct {
		ID   int
		Name string
	}

	// read in the city data JSONCityData defined in city.go
	file, err := ioutil.ReadFile("./city.json")
	if err != nil {
		log.Error().Err(err).Msg("Unable to read city.json to get names")
		return serverError(c, "Unable to retrieve city names")
	}

	var data []JSONCityData
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse city.json to get city names")
		return serverError(c, "Unable to retrieve city list")
	}

	var ret []returnData
	for _, d := range data {
		ret = append(ret, returnData{ID: int(d.ID), Name: d.Name})
	}

	return c.JSON(http.StatusOK, &ret)
}
