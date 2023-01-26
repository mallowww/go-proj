package movie

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mallowww/crud/model"
)

func Update(c echo.Context) error {
	movie := model.Movie{}
	err := c.Bind(&movies)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Can't bind the request body into movies")
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Can't convert id as string to int")
	}

	var currentMovie *model.Movie
	for i, m := range movies {
		if m.ID == strconv.Itoa(id) {
			currentMovie = &movies[i]
			break
		}
	}
	if currentMovie == nil {
		return c.JSON(http.StatusNotFound, "Not found that movie by id")
	}

	*currentMovie = movie
	return c.JSON(http.StatusOK, movies)
}
