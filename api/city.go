package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	models "github.com/mannx/weather/models"
	"gorm.io/gorm"
)

func GetCityList(c echo.Context, db *gorm.DB, cfg *models.Configuration) error {
	type Result struct {
		Cities []int
	}

	res := Result{Cities: cfg.CityIDs}
	return c.JSON(http.StatusOK, &res)
}
