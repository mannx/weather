package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"time"

	"html/template"
	"net/http"

	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// DBPath points to the database file. /app/db.db in container, ./data/db.db while developing
const DBPath = "/app/db.db"
const cityID = 6138517
const Version = 0.01

//const apiKey = "8500043bc3c464bdc0a90c69333c50b9"

// used to the store the OpenWeatherMap api key.
// passed to us through the APIKEY environment variable
var apiKey string

// Template for doing template things
type Template struct {
	templates *template.Template
}

// Render for template rendering things
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newWeatherHandler(c echo.Context) error {
	getCurrentWeather()
	return c.String(http.StatusOK, "Weather Updated Successfully!")
}

func viewWeatherHandler(c echo.Context) error {
	// open the db and show the latest weather report
	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close() // dont forget to close the db

	wd := GetWeatherData(db)

	// convert the time into a user friendly time string
	t := time.Unix(int64(wd.StoreTime), 0)
	wd.TimeString = t.String()

	return c.Render(http.StatusOK, "temp.html", wd)
}

func rawViewHandler(c echo.Context) error {
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

	fmt.Printf("\n****\n")
	displayJSON(res, 0)

	return c.String(http.StatusOK, "see log")
}

func main() {
	fmt.Printf("***\nWeather version: %v\n\n", Version)

	// get required data from environment variables
	apiKey = os.Getenv("APIKEY")
	if apiKey == "" {
		fmt.Printf("API KEY NOT FOUND!\n")
		log.Fatal()
	}

	fmt.Printf("Opening database...\n")

	// make sure the database is created before we start trying to use it
	// open the db
	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Print("Creating tables if required...\n")
	err = createDB(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Initializing templates...\n")
	t := &Template{
		templates: template.Must(template.ParseGlob("./static/*.html")),
	}

	fmt.Printf("Initializing echo and middleware...\n")

	e := echo.New()

	// middle ware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.Static("./static"))

	e.Renderer = t

	// routes
	e.GET("/view", viewWeatherHandler)
	e.GET("/new", newWeatherHandler)
	e.GET("/api/days", dayViewHandler) // handle viewing of several days of data
	e.GET("/raw", rawViewHandler)

	fmt.Printf("Starting server...\n")

	// start the server
	e.Logger.Fatal(e.Start(":8080"))
}

// return json data for the previous :num days
//  /api/days?num=X
func dayViewHandler(c echo.Context) error {
	// bind the num parameter to a variable
	var numDays int64

	err := echo.QueryParamsBinder(c).
		Int64("num", &numDays).BindError()

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}

	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close() // dont forget to close the db

	//str := fmt.Sprintf("SELECT %s FROM weather ORDER BY StoreTime LIMIT %v", SQLItems, numDays)
	str := fmt.Sprintf("SELECT * FROM (SELECT %s FROM weather ORDER BY StoreTime DESC LIMIT %v) t1 ORDER BY t1.StoreTime", SQLItems, numDays)
	wd := getWeatherDataCustom(db, str)

	// convert the time into a user friendly time string
	for i, n := range wd {
		t := time.Unix(int64(n.StoreTime), 0)
		wd[i].TimeString = t.String()
	}

	fmt.Printf("[] days ot view: %v\n", numDays)
	return c.Render(http.StatusOK, "days.html", wd)
}
