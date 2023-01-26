package movie

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mallowww/crud/model"
)

func GetById(c echo.Context) error {
	var movie *model.Movie
	for _, m := range movies {
		if m.ID == c.Param("id") {
			movie = &m
			break
		}
	}
	if movie == nil {
		return c.JSON(http.StatusNotFound, "Not found that movie by id")
	}
	return c.JSON(http.StatusOK, movie)
}
