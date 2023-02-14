package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	_ "github.com/gorilla/mux"
)

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	app := fiber.New()
	app.Get("/rentals", GetAllRentals)
	app.Get("/rentals/:name/:surname", GetRentalByName)
	app.Get("/rentals/:id", GetRentalById)

	app.Get("/query", GetNameFromQuery)
	app.Get("/query2", GetNameFromQueryParser)

	app.Get("/wildcards/*", Wildcards)
	app.Post("/rentals", CreateRental)

	app.Listen(":8080")
}

func GetAllRentals(c *fiber.Ctx) error {
	return c.SendString("get rentals")
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
