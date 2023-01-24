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
)

type Movie struct {
	ID       string    `json:"id"`
	ISBN     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "send 200 to tell u it's fine")
	})
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	err := e.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
