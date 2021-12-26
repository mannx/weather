package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"

	"encoding/json"
	"net/http"

	//_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

const cityID = 6138517
const apiKey = "8500043bc3c464bdc0a90c69333c50b9"

type weatherData struct {
	time      float64
	code      float64
	temp      float64
	feelsLike float64
	pressure  float64
	humidity  float64
	windSpeed float64
	windDir   float64
	windGust  float64
	rain1h    float64
	rain3h    float64
	snow1h    float64
	snow3h    float64

	icon    string // icon name for the current weather
	sunrise float64
	sunset  float64
	name    string // name of the city we are storing data for
}

// displayJson outputs the json structure to stdout with indentations
//  indent determines the current indentation level in spaces
func displayJSON(d map[string]interface{}, indent int) {
	var istr string

	for i := 0; i < indent; i++ {
		istr = fmt.Sprintf(" %s", istr)
	}

	for k, v := range d {
		switch n := v.(type) {
		case map[string]interface{}: // data deeper in, only display the key
			fmt.Printf("%s%s:\n", istr, k)
			displayJSON(n, indent+5) // display the next leve
		case []interface{}: // slices of things
			fmt.Printf("* %s%s:\n", istr, k)
			for _, x := range n {
				switch m := x.(type) {
				case map[string]interface{}:
					displayJSON(m, indent+5)
				}
			}
		default:
			fmt.Printf("%s%s: %v (%T)\n", istr, k, v, v)
		}
	}
}

// weatherToMap converts a JSON type mapping into a singular flat map with
//  sub levels starting with the tld string
//	ie.
//		"weather":
//			"code": 102
//	converts to map["weather.code"]=102
//
func weatherToMap(data map[string]interface{}, output *map[string]interface{}, tld string) {
	fmt.Println("weatherToMap() starting...")

	for k, v := range data {
		switch n := v.(type) {
		case map[string]interface{}: // data deeper in
			t := fmt.Sprintf("%s.%s", tld, k)
			fmt.Printf("- %s", t)
			weatherToMap(n, output, t)
		case []interface{}:
			// array of map[]'s
			fmt.Printf("** %s.%s", tld, k)
			for _, x := range n {
				switch m := x.(type) {
				case map[string]interface{}:
					t := fmt.Sprintf("%s.%s", tld, k)
					fmt.Printf("++ %s\n", t)
					weatherToMap(m, output, t)
				}
			}
		default:
			fmt.Printf("%s.%s: %v (%T)\n", tld, k, v, v)
			// get the name of this item using the tld, and current key
			// and store the value
			t := fmt.Sprintf("%s.%s", tld, k)
			//(*output)[t] = v.(float64)
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

func getWeatherData(input map[string]interface{}) weatherData {
	wd := weatherData{}

	// copy over the float64's values
	wd.code = input[".weather.id"].(float64)
	wd.temp = input[".main.temp"].(float64)
	wd.feelsLike = input[".main.feels_like"].(float64)
	wd.pressure = input[".main.pressure"].(float64)
	wd.humidity = input[".main.humidity"].(float64)
	wd.windSpeed = input[".wind.speed"].(float64)
	wd.windDir = input[".wind.deg"].(float64)
	wd.windGust = getFloat64(input, ".wind.gust")
	wd.rain1h = getFloat64(input, ".rain.rain.1h")
	wd.rain3h = getFloat64(input, ".rain.rain.3h")
	wd.snow1h = getFloat64(input, ".snow.snow.1h")
	wd.snow3h = getFloat64(input, ".snow.snow.3h")

	wd.icon = input[".weather.icon"].(string)
	wd.sunrise = getFloat64(input, ".sys.sunrise")
	wd.sunset = getFloat64(input, ".sys.runset")
	wd.name = input[".name"].(string)

	wd.time = getFloat64(input, ".dt")

	return wd
}

func main() {
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

	displayJSON(res, 0)

	output := make(map[string]interface{})
	weatherToMap(res, &output, "")

	wd := getWeatherData(output)
	fmt.Printf("code: %v (%T)\n", wd.code, wd.code)

	//	fmt.Println("\n\n*********\n\n")
	//	fmt.Printf("Weather code: %v\n", wd.code)

	// open the db
	//db, err := sql.Open("sqlite3", "./db.db")
	db, err := sql.Open("sqlite", "./db.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	/*

		type weatherData struct {
			time      float64
			code      float64
			temp      float64
			feelsLike float64
			pressure  float64
			humidity  float64
			windSpeed float64
			windDir   float64
			windGust  float64
			rain1h    float64
			rain3h    float64
			snow1h    float64
			snow3h    float64

			icon    string // icon name for the current weather
			sunrise float64
			sunset  float64
			name    string // name of the city we are storing data for
		}
	*/

	sql := "CREATE TABLE IF NOT EXISTS weather (id INTEGER PRIMARY KEY, time TEXT, code REAL, temp REAL, feelsLike REAL, pressure REAL, humidity REAL, windSpeed REAL, windDir REAL, "
	sql = fmt.Sprintf("%swindGust REAL, rain1h REAL, rain3h REAL, snow1h REAL, snow3h REAL, icon TEXT, city INTEGER)", sql)

	//state, _ := db.Prepare("CREATE TABLE IF NOT EXISTS weather (id INTEGER PRIMARY KEY, time TEXT)")
	state, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	state.Exec()

	/*state, _ = db.Prepare("INSERT INTO weather (time) VALUES (?)")
	state.Exec(fmt.Sprintf("%s", now.Format(time.RFC1123)))
	fmt.Println(cityID)*/

	// # of values to insert: 13
	sql = "INSERT INTO weather (time, code, temp, feelsLike, pressure, humidity, windSpeed, windDir, windGust, rain1h, rain3h, snow1h, snow3h, icon, city) VALUES (?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?)"
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close() // make sure to free resources

	_, err = stmt.Exec(wd.time, wd.code, wd.temp, wd.feelsLike, wd.pressure, wd.humidity, wd.windSpeed, wd.windDir, wd.windGust, wd.rain1h, wd.rain3h, wd.snow1h, wd.snow3h, wd.icon, cityID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Current weather logged sucessfully!")
}
