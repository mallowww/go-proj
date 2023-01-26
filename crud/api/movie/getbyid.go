package movie

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetById(c echo.Context) error {
	// id := c.Param("")
	return c.JSON(http.StatusOK, movies)
}
