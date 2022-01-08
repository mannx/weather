package main

import (
	"encoding/json"
	"fmt"
	"io"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"time"

	"html/template"
	"net/http"

	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/robfig/cron/v3"
)

// DBPath points to the database file. /app/db.db in container, ./data/db.db while developing
//const DBPath = "/app/db.db"
const DBPath = "./data/db.db"

const cityID = 6138517

// Version of the current build
const Version = 0.02

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

func viewWeatherHandler(c echo.Context) error {
	// open the db and show the latest weather report
	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		return err
	}

	defer db.Close() // dont forget to close the db

	wd, err := GetWeatherData(db)
	if err != nil {
		return err
	}

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
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var res map[string]interface{}
	json.Unmarshal([]byte(body), &res)

	fmt.Printf("\n****\n")
	displayJSON(res, 0)

	return c.String(http.StatusOK, "see log")
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // can chagne to zerolog.DebugLevel for more info, or ErrorLevel for just errors

	log.Info().Msgf("Weather version: %v", Version)

	// get required data from environment variables
	apiKey = os.Getenv("APIKEY")
	if apiKey == "" {
		log.Fatal().Msg("APIKEY environment variable not defined.")
	}

	log.Info().Msg("Creating database if required...")

	// make sure the database is created before we start trying to use it
	// open the db
	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to open database")
	}

	defer db.Close()

	log.Info().Msg("Creating tables if required...")

	err = createDB(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create required tables")
	}

	log.Info().Msg("Initializing all HTML templates")

	t := &Template{
		templates: template.Must(template.ParseGlob("./static/*.html")),
	}

	log.Info().Msg("Initializing echo and middleware")

	e := echo.New()

	// middle ware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		//AllowOrigins: []string{"http://localhost"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(middleware.Static("./static"))

	e.Renderer = t

	// routes
	e.GET("/view", viewWeatherHandler)
	e.GET("/api/days", dayViewHandler) // handle viewing of several days of data
	e.GET("/raw", rawViewHandler)

	e.GET("/api/json/chart", chartViewHandler) // returns the data we need to display a simple chart
	e.GET("/api/test", testViewHandler)

	log.Info().Msg("Setting up cron job for updates")

	c := cron.New()
	expr := os.Getenv("WEATHER_UPDATE_SCHEDULE")
	if expr == "" {
		// update not set, default to hourly
		log.Info().Msg("WEATHER_UPDATE_SCHEDULE not set, defaulting to @hourly")
		expr = "@hourly"
	}

	c.AddFunc(expr, updateWeatherFunc)
	c.Start() // make sure to start the jobs

	log.Info().Msg("Starting server...")

	// start the server
	e.Logger.Fatal(e.Start(":8080"))
}

// updateWeatherFunc is called to retrieve and store the current weather
// this should only be called from a scheduled job
func updateWeatherFunc() {
	err := getCurrentWeather()
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve weather")
	} else {
		log.Info().Msg("Weather updated successfully")
	}
}

// DayViewData temp thing
type DayViewData struct {
	Weather []WeatherData
	IPAddr  int
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
		return err
	}

	defer db.Close() // dont forget to close the db

	str := fmt.Sprintf("SELECT * FROM (SELECT %s FROM weather ORDER BY StoreTime DESC LIMIT %v) t1 ORDER BY t1.StoreTime", SQLItems, numDays)
	wd, err := getWeatherDataCustom(db, str)
	if err != nil {
		log.Error().Err(err).Str("handler", "dayViewHandler").Msg("Unable to retrieve weather from db")
	}

	// convert the time into a user friendly time string
	for i, n := range wd {
		t := time.Unix(int64(n.StoreTime), 0)
		wd[i].TimeString = t.String()
	}

	log.Info().Str("handler", "dayViewHandler").Msgf("Returning %v days to view", numDays)
	return c.Render(http.StatusOK, "days.html", DayViewData{wd, 12354})
}

type chartItem struct {
	Temp      float64 `json: "temp"`
	FeelsLike float64 `json: "feelsLike"`
}

type chartJSON struct {
	Data []chartItem `json: "data"`
}

//
// returns a json data for the previous :num entries containing temp and feelslike temp
func chartViewHandler(c echo.Context) error {
	var days int64

	err := echo.QueryParamsBinder(c).Int64("num", &days).BindError()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		return err
	}

	defer db.Close() // dont forget to close the db

	str := fmt.Sprintf("SELECT * FROM (SELECT %s FROM weather ORDER BY StoreTime DESC LIMIT %v) t1 ORDER BY t1.StoreTime", SQLItems, days)
	wd, err := getWeatherDataCustom(db, str)
	if err != nil {
		log.Error().Err(err).Str("handler", "dayViewHandler").Msg("Unable to retrieve weather from db")
	}

	// convert the time into a user friendly time string
	for i, n := range wd {
		t := time.Unix(int64(n.StoreTime), 0)
		wd[i].TimeString = t.String()
	}

	log.Debug().Msgf("building json object (count: %v)", len(wd))

	dta := chartJSON{}
	dta.Data = make([]chartItem, len(wd))

	for i, n := range wd {
		dta.Data[i].Temp = n.Temp
		log.Debug().Msgf(" [%v] Temp: %v (%v)", i, n.Temp, dta.Data[i].Temp)
	}

	return c.JSON(http.StatusOK, &dta)

}

func testViewHandler(c echo.Context) error {

	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		return err
	}

	defer db.Close() // dont forget to close the db

	// get the lastest weather entry
	wd, err := GetWeatherData(db)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &wd)
}
