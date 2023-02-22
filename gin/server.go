package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []Book{
	{ID: "1", Title: "title1", Author: "author1"},
	{ID: "2", Title: "title2", Author: "author2"},
	{ID: "3", Title: "title3", Author: "author3"},
}

func main() {
	r := gin.Default()
	r.GET("/", handleHome)
	r.GET("/books", handleGetBooks)
	r.GET("/books/:id", handleGetBookByID)
	r.POST("/books", handleSaveBook)
	r.DELETE("/books/:id", handleDeleteBook)
	r.Run()
}

func handleHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the book store!",
	})
}

func handleGetBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

func handleGetBookByID(c *gin.Context) {
	bookId := c.Param("id")
	for _, book := range books {
		if book.ID == bookId {
			c.JSON(http.StatusOK, book)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Book not found",
	})
}

func handleSaveBook(c *gin.Context) {
	var book Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book data",
		})
		return
	}

	for _, existingBook := range books {
		if existingBook.ID == book.ID {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Book ID already exists",
			})
			return
		}
	}

	books = append(books, book)
	c.JSON(http.StatusCreated, book)
}

func handleDeleteBook(c *gin.Context) {
	id := c.Param("id")

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
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
