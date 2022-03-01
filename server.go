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
	e.GET("/api/24hr", func(c echo.Context) error { return api.Handle24hrView(c, DB) })
	e.GET("/api/latest", func(c echo.Context) error { return api.GetLatestWeatherView(c, DB) })
	e.GET("/api/daily", func(c echo.Context) error { return api.GetDailyWeatherView(c, DB) })
	e.GET("/api/cities", func(c echo.Context) error { return api.GetCityList(c, DB) })

	e.POST("/api/city/add", func(c echo.Context) error { return api.ValidateCity(c, DB) })
	e.POST("/api/city/confirm", func(c echo.Context) error { return api.AddCity(c, DB) })

	e.GET("/api/migrate", func(c echo.Context) error { return api.ConfigToDatabaseHandler(c, DB, &Config) })
	return e
}
