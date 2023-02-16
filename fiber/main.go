package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	_ "github.com/gorilla/mux"
)

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	app := fiber.New()
	app.Use("/rentals", RentalMiddleware)
	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE",
		AllowCredentials: true,
	}))
	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Bangkok",
	}))

	app.Get("/rentals", GetAllRentals)
	app.Get("/rentals/:name/:surname", GetRentalByName)
	app.Get("/rentals/:id", GetRentalById)

	app.Get("/query", GetNameFromQuery)
	app.Get("/query2", GetNameFromQueryParser)
	app.Get("/wildcards/*", Wildcards)

	app.Post("/rentals", CreateRental)

	// static file
	app.Static("/", "./wwwroot", fiber.Static{
		Index:         "index.html",
		CacheDuration: time.Second * 7,
	})

	app.Get("/error", NewError)

	app.Listen(":8080")
}

func RentalMiddleware(c *fiber.Ctx) error {
	c.Locals("name", "var")
	fmt.Println("before")
	err := c.Next()
	fmt.Println("after")
	return err
}

func GetAllRentals(c *fiber.Ctx) error {
	name := c.Locals("name")
	return c.SendString(fmt.Sprintf("get rentals for %v", name))
}

func GetRentalByName(c *fiber.Ctx) error {
	name := c.Params("name")
	surname := c.Params("surname")
	return c.SendString("name: " + name + ", surname:" + surname)
}

func GetRentalById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	return c.SendString(fmt.Sprintf("id: %v", id))
}

func GetNameFromQuery(c *fiber.Ctx) error {
	name := c.Query("name")
	surename := c.Query("surname")
	return c.SendString("name: " + name + " surname: " + surename)
}

func GetNameFromQueryParser(c *fiber.Ctx) error {
	person := Person{}
	c.QueryParser(&person)
	return c.JSON(person)
}

func CreateRental(c *fiber.Ctx) error {
	return c.SendString("post rentals")
}

func Wildcards(c *fiber.Ctx) error {
	wildcard := c.Params("*")
	return c.SendString(wildcard)
}

func NewError(c *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusNotFound, "message not found")
}
