package movie

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mallowww/crud/model"
)

func Create(c echo.Context) error {
	var newMovie []model.Movie
	err := c.Bind(&newMovie)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	movies = append(movies, newMovie...)
	return c.JSON(http.StatusCreated, newMovie)
}
