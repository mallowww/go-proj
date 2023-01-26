package movie

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Delete(c echo.Context) error {
	for i, m := range movies {
		if m.ID == c.Param("id") {
			movies = append(movies[:i], movies[i+1:]...)
			return c.JSON(http.StatusOK, "Movie deleted")
		}
	}
	return c.JSON(http.StatusNotFound, "Movie id not found")
}
