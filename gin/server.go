package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// avoid error shadow casting
	var err error
	db, err = gorm.Open(sqlite.Open("sample-library.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	db.AutoMigrate(&Book{})

	// initBooks()
	r := gin.Default()
	r.GET("/", handleHomepage)
	r.GET("/allBooks", handleGetAllBooks)
	r.GET("/allBooks/:id", handleGetBook)
	r.POST("/allBooks", handleCreateBook)
	r.DELETE("/allBooks/:id", handleDeleteBookById)
	r.Run()
}

type Book struct {
	gorm.Model
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var (
	allBooks []Book
)

func initBooks() {
	allBooks = []Book{
		{ID: "1", Title: "title1", Author: "author1"},
		{ID: "2", Title: "title2", Author: "author2"},
		{ID: "3", Title: "title3", Author: "author3"},
	}
}

func handleHomepage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the book store!",
	})
}

func handleGetAllBooks(c *gin.Context) {
	var books []Book
	if err := db.Find(&books).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, allBooks)
}

func handleGetBook(c *gin.Context) {
	bookId := c.Param("id")
	var book Book
	if err := db.Where("id = ?", bookId).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Book not found",
		})
		return
	}

	c.JSON(http.StatusOK, book)
}

func handleCreateBook(c *gin.Context) {
	var book Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book data",
		})
		return
	}

	for _, existingBook := range allBooks {
		if existingBook.ID == book.ID {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Book ID already exists",
			})
			return
		}
	}

	allBooks = append(allBooks, book)
	c.JSON(http.StatusCreated, book)
}

func handleDeleteBookById(c *gin.Context) {
	id := c.Param("id")

	for i, book := range allBooks {
		if book.ID == id {
			allBooks = append(allBooks[:i], allBooks[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "Book deleted successfully",
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Book not found",
	})
}
