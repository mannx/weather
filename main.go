package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"text/template"

	"encoding/json"
	"net/http"

	_ "modernc.org/sqlite"
)

const cityID = 6138517
const apiKey = "8500043bc3c464bdc0a90c69333c50b9"

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

	wd.time = getFloat64(input, ".dt")

	return wd
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

	// open the db
	db, err := sql.Open("sqlite", "./db.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	sql := "CREATE TABLE IF NOT EXISTS weather (id INTEGER PRIMARY KEY, time TEXT, code REAL, temp REAL, feelsLike REAL, pressure REAL, humidity REAL, windSpeed REAL, windDir REAL, "
	sql = fmt.Sprintf("%swindGust REAL, rain1h REAL, rain3h REAL, snow1h REAL, snow3h REAL, icon TEXT, city INTEGER)", sql)

	state, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	state.Exec()

	// # of values to insert: 13
	sql = "INSERT INTO weather (time, code, temp, feelsLike, pressure, humidity, windSpeed, windDir, windGust, rain1h, rain3h, snow1h, snow3h, icon, city) VALUES (?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?)"
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close() // make sure to free resources

	_, err = stmt.Exec(wd.time, wd.code, wd.Temp, wd.FeelsLike, wd.Pressure, wd.Humidity, wd.WindSpeed, wd.WindDir, wd.WindGust, wd.Rain1h, wd.Rain3h, wd.Snow1h, wd.Snow3h, wd.Icon, cityID)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Current weather logged sucessfully!")
}

func newWeatherHandler(w http.ResponseWriter, r *http.Request) {
	getCurrentWeather()
	fmt.Fprintf(w, "Current weather logged successfully!")
}

func viewWeatherHandler(w http.ResponseWriter, req *http.Request) {
	// open the db and show the latest weather report
	db, err := sql.Open("sqlite", "db.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close() // dont forget to close the db

	wd := GetWeatherData(db)

	t, err := template.ParseFiles("latest.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, &wd)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/new", newWeatherHandler)
	http.HandleFunc("/view", viewWeatherHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
