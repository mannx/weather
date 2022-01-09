package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func initServer() *echo.Echo {
	e := echo.New()

	// middle ware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(middleware.Static("./static"))

	// routes
	e.GET("/api/test", testViewHandler)
	e.GET("/api/24hr", handle24hrView)

	return e
}

func testViewHandler(c echo.Context) error {
	db, err := openDB()
	if err != nil {
		return err
	}

	defer db.Close()

	// get the lastest weather entry
	wd, err := GetWeatherData(db)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &wd)
}

func handle24hrView(c echo.Context) error {
	db, err := openDB()
	if err != nil {
		return err
	}

	defer db.Close()

	// retrieve all entries in the last 24 hours
	now := time.Now()
	prev := now.Add(-time.Hour * 24)

	// retrieve the data
	sql := fmt.Sprintf("SELECT %s FROM weather WHERE StoreTime BETWEEN %v AND %v ORDER BY id DESC", SQLItems, prev.Unix(), now.Unix())

	wd, err := getWeatherDataCustom(db, sql)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &wd)
}
