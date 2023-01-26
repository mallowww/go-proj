package movie_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mallowww/crud/model"
	"github.com/mallowww/crud/movie"
)

func TestCreate(t *testing.T) {
	e := echo.New()
	//arrange
	testMovie := model.Movie{ID: "1", ISBN: "438227", Title: "Movie One", Director: &model.Director{Firstname: "Number", Lastname: "One"}}

	req := httptest.NewRequest(http.MethodPost, "/movies", testMovie)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// act
	c := e.NewContext(req, rec)
	if err := movie.Create(c); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// assert
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	expected := `[{"id":"1","isbn":"438227","title":"Movie One","director":{"firstname":"Number","lastname":"One"}}]` + "\n"
	if rec.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, rec.Body.String())
	}
}
