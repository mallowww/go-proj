package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	_ "github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	jwtware "github.com/gofiber/jwt/v3"
)

var db *sqlx.DB

const jwtSecret = "authen"

type User struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username json:"username"`
	Password string `db:"password" json:"password"`
}

type SignupRequest struct {
	Username string `json:"username`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username`
	Password string `json:"password"`
}

func main() {
	var err error
	url := os.Getenv("DATABASE_URL")
	db, err = sqlx.Open("mysql", url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	app := fiber.New()

	app.Use("home", jwtware.New(jwtware.Config{
		SigningMethod: "HS256",
		SigningKey: []byte(jwtSecret),
		SuccessHandler:func(c *fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.ErrUnauthorized
		},
	}))

	app.Post("/signup", Signup)
	app.Post("/login", Login)
	app.Get("/home", Home)

	app.Listen(":8080")
}

func Signup(c *fiber.Ctx) error {
	request := SignupRequest{}
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	if request.Username == "" || request.Password == "" {
		return fiber.ErrUnprocessableEntity
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	query := "insert user (username, password) values (?, ?)"
	result, err := db.Exec(query, request.Username, string(password))
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	user := User{
		Id:       int(id),
		Username: request.Username,
		Password: string(password),
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func Login(c *fiber.Ctx) error {
	request := LoginRequest{}
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	if request.Username == "" || request.Password == "" {
		return fiber.ErrUnprocessableEntity
	}

	user := User{}
	query := "select id, username, password from user where username=?"
	err = db.Get(&user, query, request.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Incorrect username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Incorrect username or password")
	}

	expiresAt := jwt.TimeFunc().Add(time.Hour * 24)
	cliams := jwt.StandardClaims{
		Issuer:    strconv.Itoa(user.Id),
		ExpiresAt: jwt.At(expiresAt),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	token, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"jwtToken": token,
	})
}

func Home(c *fiber.Ctx) error {
	return c.SendString("Home")
}

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func Fiber() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: false,
		StrictRouting: false,
	})
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

	// Group
	v1 := app.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("version", "v1")
		return c.Next()
	})
	v1.Get("/rentals", func(c *fiber.Ctx) error {
		return c.SendString("v1")
	})

	v2 := app.Group("/v2", func(c *fiber.Ctx) error {
		c.Set("version", "v2")
		return c.Next()
	})
	v2.Get("/rentals", func(c *fiber.Ctx) error {
		return c.SendString("v2")
	})

	// Mount
	userApp := fiber.New()
	userApp.Get("/login", func(c *fiber.Ctx) error {
		return c.SendString("Login")
	})
	app.Mount("/user", userApp)

	// Server
	app.Server().MaxConnsPerIP = 1
	app.Get("/server", func(c *fiber.Ctx) error {
		time.Sleep(time.Second * 30)
		return c.SendString("server")
	})

	// Environtment
	app.Get("/env", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"BaseURL":     c.BaseURL(),
			"Hostname":    c.Hostname(),
			"IP":          c.IP(),
			"IPs":         c.IPs(),
			"OriginalURL": c.OriginalURL(),
			"Path":        c.Path(),
			"Protocol":    c.Protocol(),
			"Subdomains":  c.Subdomains(),
		})
	})

	// Body
	app.Post("/body", func(c *fiber.Ctx) error {
		fmt.Printf("IsJson: %v\n", c.Is("json"))

		data := map[string]interface{}{}
		err := c.BodyParser(&data)
		if err != nil {
			return err
		}

		fmt.Println(data)
		return nil

	})

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
