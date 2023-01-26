package movie

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mallowww/crud/model"
)

var movies = []model.Movie{
	{ID: "1", ISBN: "438227", Title: "Movie One", Director: &model.Director{Firstname: "Number", Lastname: "One"}},
	{ID: "2", ISBN: "45455", Title: "Movie Two", Director: &model.Director{Firstname: "Number", Lastname: "Two"}},
}

func GetAll(c echo.Context) error {
	return c.JSON(http.StatusOK, movies)
}
