package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	api "github.com/mannx/weather/api"
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
	/*e.GET("/api/24hr", handle24hrView)
	e.GET("/api/latest", getLatestWeatherView)
	e.GET("/api/daily", getDailyWeatherView)*/

	e.GET("/api/24hr", func(c echo.Context) error { return api.Handle24hrView(c, DB) })
	e.GET("/api/latest", func(c echo.Context) error { return api.GetLatestWeatherView(c, DB) })
	e.GET("api/daily", func(c echo.Context) error { return api.GetDailyWeatherView(c, DB) })

	return e
}
