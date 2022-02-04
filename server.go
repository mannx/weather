package main

import (
	"net/http"

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
