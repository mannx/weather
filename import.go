package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	models "github.com/mannx/weather/models"
	"github.com/rs/zerolog/log"
)

// weatherToMap converts a JSON type mapping into a singular flat map with
//  sub levels starting with the tld string
//	ie.
//		"weather":
//			"code": 102
//	converts to map["weather.code"]=102
//
func weatherToMap(data map[string]interface{}, output *map[string]interface{}, tld string) {
	for k, v := range data {
		switch n := v.(type) {
		case map[string]interface{}: // data deeper in
			t := fmt.Sprintf("%s.%s", tld, k)
			weatherToMap(n, output, t)
		case []interface{}:
			// array of map[]'s
			for _, x := range n {
				switch m := x.(type) {
				case map[string]interface{}:
					t := fmt.Sprintf("%s.%s", tld, k)
					weatherToMap(m, output, t)
				}
			}
		default:
			// get the name of this item using the tld, and current key
			// and store the value
			t := fmt.Sprintf("%s.%s", tld, k)
			(*output)[t] = v
		}
	}
}

func getFloat64(in map[string]interface{}, key string) float64 {
	v := in[key]
	if v == nil {
		return 0.0
	}
	return v.(float64)
}

func getWeatherData(input map[string]interface{}) models.WeatherData {
	wd := models.WeatherData{}

	// copy over the float64's values
	wd.Code = input[".weather.id"].(float64)
	wd.Temp = input[".main.temp"].(float64)
	wd.FeelsLike = input[".main.feels_like"].(float64)
	wd.Pressure = input[".main.pressure"].(float64)
	wd.Humidity = input[".main.humidity"].(float64)
	wd.WindSpeed = input[".wind.speed"].(float64)
	wd.WindDir = input[".wind.deg"].(float64)
	wd.WindGust = getFloat64(input, ".wind.gust")
	wd.Rain1h = getFloat64(input, ".rain.1h")
	wd.Rain3h = getFloat64(input, ".rain.3h")
	wd.Snow1h = getFloat64(input, ".snow.1h")
	wd.Snow3h = getFloat64(input, ".snow.3h")

	wd.Icon = input[".weather.icon"].(string)
	wd.Sunrise = getFloat64(input, ".sys.sunrise")
	wd.Sunset = getFloat64(input, ".sys.runset")
	wd.Name = input[".name"].(string)

	return wd
}

func getCurrentWeather(cityID int) error {
	log.Debug().Msgf("getCurrentWeather(%v) ==>", cityID)

	// retrieve the current weather in json format
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?id=%v&appid=%v&units=metric", cityID, Environment.ApiKey)
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var res map[string]interface{}
	json.Unmarshal([]byte(body), &res)

	output := make(map[string]interface{})
	weatherToMap(res, &output, "")

	wd := getWeatherData(output)

	// get the current time and store it
	wd.StoreTime = time.Now().Unix()
	wd.CityID = cityID

	// save to the db
	DB.Save(&wd)
	return nil
}
