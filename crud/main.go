package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	movie "github.com/mallowww/crud/api/movie"
)

func Home(c echo.Context) error {
	return c.String(http.StatusOK, "home: send 200 to tell u it's fine")
}

func main() {
	e := echo.New()
	// Routing
	e.GET("/", Home)
	e.GET("/movies", movie.GetAll)
	e.GET("/movies/{id}", getMovie)
	e.POST("/movies", movie.Create)
	// e.PUT("/movies/{id}", updateMovie)
	// e.DELTE("/movies/{id}", deleteMovie)

	addr := ":1555"
	log.Println("Server started at port", addr)
	go func() {
		err := e.Start(addr)
		if err != nil {
			fmt.Println("Signal status:", err, "\n>(o o )!! Shutting down the server gracefully. . .")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := e.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
