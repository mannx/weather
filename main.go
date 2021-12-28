package main

import (
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

// set the path to the database when used to access it
// if testing locally, set to ./data/db.db
const DBPath = "/app/db.db"
const cityID = 6138517

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

func main() {
	// get required data from environment variables
	apiKey = os.Getenv("APIKEY")
	if apiKey == "" {
		fmt.Printf("API KEY NOT FOUND!\n")
		log.Fatal()
	} else {
		fmt.Printf("API Key: %v\n", apiKey)
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

	e.Renderer = t

	// routes
	e.GET("/view", viewWeatherHandler)
	e.GET("/new", newWeatherHandler)

	fmt.Printf("Starting server...\n")

	// start the server
	e.Logger.Fatal(e.Start(":8080"))
}
