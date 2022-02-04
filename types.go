package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	models "github.com/mannx/weather/models"
	"github.com/rs/zerolog/log"
)

// SQLItems contains the database items we are pulling out during a retrieve
const SQLItems string = "StoreTime, temp, feelsLike, humidity, windSpeed, windDir, rain1h, snow1h, icon"

// GetWeatherData builds a WeatherData object from the database using the last entry
func GetWeatherData(db *sql.DB) (models.WeatherData, error) {
	q := "SELECT %s FROM weather ORDER BY id DESC LIMIT 1"
	query := fmt.Sprintf(q, SQLItems)

	wd, err := getWeatherDataCustom(db, query)
	if err != nil {
		return models.WeatherData{}, err
	}

	if len(wd) > 1 {
		log.Error().Str("function", "GetWeatherData").Msg("Returned more than 1 result. Expected only 1")
	}

	return wd[0], nil
}

func getWeatherDataCustom(db *sql.DB, query string) ([]models.WeatherData, error) {
	r, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	wda := make([]models.WeatherData, 0)

	for r.Next() {
		wd := models.WeatherData{}
		if err := r.Scan(&wd.StoreTime, &wd.Temp, &wd.FeelsLike, &wd.Humidity, &wd.WindSpeed, &wd.WindDir, &wd.Rain1h, &wd.Snow1h, &wd.Icon); err != nil {
			return nil, err
		}

		wda = append(wda, wd)
	}

	return wda, nil
}

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

// openDB opens a connection to the database, caller is responsible for closing the connection.
func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// createDB is used to create the database if it doesnt exist
// should be called on startup only
func createDB(db *sql.DB) error {
	sql := "CREATE TABLE IF NOT EXISTS weather (id INTEGER PRIMARY KEY, Time REAL, code REAL, temp REAL, feelsLike REAL, pressure REAL, humidity REAL, windSpeed REAL, windDir REAL, "
	sql = fmt.Sprintf("%swindGust REAL, rain1h REAL, rain3h REAL, snow1h REAL, snow3h REAL, icon TEXT, city INTEGER, Timezone REAL, StoreTime INTEGER)", sql)

	state, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	state.Exec()
	return nil
}

func getCurrentWeather(cityID int) error {
	// retrieve the current weather in json format
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?id=%v&appid=%v&units=metric", cityID, Config.APIKey)
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

	// open the db
	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		return err
	}

	defer db.Close()

	// # of values to insert: 13
	sql := "INSERT INTO weather (code, temp, feelsLike, pressure, humidity, windSpeed, windDir, windGust, rain1h, rain3h, snow1h, snow3h, icon, city, StoreTime) VALUES (?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?)"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}

	defer stmt.Close() // make sure to free resources

	_, err = stmt.Exec(wd.Code, wd.Temp, wd.FeelsLike, wd.Pressure, wd.Humidity, wd.WindSpeed, wd.WindDir, wd.WindGust, wd.Rain1h, wd.Rain3h, wd.Snow1h, wd.Snow3h, wd.Icon, cityID, wd.StoreTime)
	return err
}
