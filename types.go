package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// WeatherData structure holds data for a given weather pull
type WeatherData struct {
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

	TimeString string // used only before send off to the template engine to display a user friendly Time string
	StoreTime  int64  // system time when this entry was added to the db
}

// GetWeatherData builds a WeatherData object from the database using the last entry
func GetWeatherData(db *sql.DB) WeatherData {
	s := "SELECT StoreTime, temp, feelsLike, humidity, windSpeed, windDir, rain1h, snow1h, icon FROM weather ORDER BY id DESC LIMIT 1"
	r, err := db.Query(s)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Close()

	wd := WeatherData{}

	for r.Next() {
		if err := r.Scan(&wd.StoreTime, &wd.Temp, &wd.FeelsLike, &wd.Humidity, &wd.WindSpeed, &wd.WindDir, &wd.Rain1h, &wd.Snow1h, &wd.Icon); err != nil {
			log.Fatal(err)
		}
	}

	return wd
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

func getWeatherData(input map[string]interface{}) WeatherData {
	wd := WeatherData{}

	// copy over the float64's values
	wd.code = input[".weather.id"].(float64)
	wd.Temp = input[".main.temp"].(float64)
	wd.FeelsLike = input[".main.feels_like"].(float64)
	wd.Pressure = input[".main.pressure"].(float64)
	wd.Humidity = input[".main.humidity"].(float64)
	wd.WindSpeed = input[".wind.speed"].(float64)
	wd.WindDir = input[".wind.deg"].(float64)
	wd.WindGust = getFloat64(input, ".wind.gust")
	wd.Rain1h = getFloat64(input, ".rain.rain.1h")
	wd.Rain3h = getFloat64(input, ".rain.rain.3h")
	wd.Snow1h = getFloat64(input, ".snow.snow.1h")
	wd.Snow3h = getFloat64(input, ".snow.snow.3h")

	wd.Icon = input[".weather.icon"].(string)
	wd.sunrise = getFloat64(input, ".sys.sunrise")
	wd.sunset = getFloat64(input, ".sys.runset")
	wd.name = input[".name"].(string)

	return wd
}

// createDB is used to create the database if it doesnt exist
// should be called on startup only
func createDB(db *sql.DB) error {
	sql := "CREATE TABLE IF NOT EXISTS weather (id INTEGER PRIMARY KEY, Time REAL, code REAL, temp REAL, feelsLike REAL, pressure REAL, humidity REAL, windSpeed REAL, windDir REAL, "
	sql = fmt.Sprintf("%swindGust REAL, rain1h REAL, rain3h REAL, snow1h REAL, snow3h REAL, icon TEXT, city INTEGER, Timezone REAL, StoreTime INTEGER)", sql)

	state, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	state.Exec()
	return nil
}

func getCurrentWeather() {
	// retrieve the current weather in json format
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?id=%v&appid=%v&units=metric", cityID, apiKey)
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}
	json.Unmarshal([]byte(body), &res)

	output := make(map[string]interface{})
	weatherToMap(res, &output, "")

	wd := getWeatherData(output)

	// get the current time and store it
	wd.StoreTime = time.Now().Unix()

	// open the db
	db, err := sql.Open("sqlite", "./db.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// # of values to insert: 13
	sql := "INSERT INTO weather (code, temp, feelsLike, pressure, humidity, windSpeed, windDir, windGust, rain1h, rain3h, snow1h, snow3h, icon, city, StoreTime) VALUES (?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?)"
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close() // make sure to free resources

	_, err = stmt.Exec(wd.code, wd.Temp, wd.FeelsLike, wd.Pressure, wd.Humidity, wd.WindSpeed, wd.WindDir, wd.WindGust, wd.Rain1h, wd.Rain3h, wd.Snow1h, wd.Snow3h, wd.Icon, cityID, wd.StoreTime)
	if err != nil {
		log.Fatal(err)
	}
}
