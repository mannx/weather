package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	models "github.com/mannx/weather/models"
)

// serverError returns a basic error response message
func serverError(c echo.Context, msg string) error {
	sr := models.ServerResponse{
		Error:   true,
		Message: msg,
	}

	return c.JSON(http.StatusOK, &sr)
}
