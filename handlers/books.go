package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Pages      int    `json:"pages"`
	CheckedOut bool   `json:"checkedOut"`
	Patron     string `json:"patron"`
}

// preload a set of books
var Books = []Book{
	{ID: "2", Title: "Chronicles of Narnia", Author: "CS Lewis", Pages: 562, CheckedOut: false},
	{ID: "1", Title: "Lord of the Rings", Author: "JRR Tolkien", Pages: 1077, CheckedOut: false},
	{ID: "3", Title: "Cat in the Hat", Author: "Dr. Seuss", Pages: 32, CheckedOut: false},
	{ID: "4", Title: "Super Powereds", Author: "Drew Hayes", Pages: 398, CheckedOut: false},
	{ID: "5", Title: "Tale of Two Cities", Author: "Charles Dickens", Pages: 1200, CheckedOut: false},
}

func GetBooks(c *gin.Context) {
	if len(Books) > 0 {
		c.JSON(http.StatusOK, Books)
		return
	}
	err := errors.New("there are no books in the library")
	log.Error().Err(err).Msg("")
	c.JSON(http.StatusNotFound, gin.H{"message": "There are no books in the library"})

}

func GetBook(c *gin.Context) {
	id := c.Param("id")
	message := "ID: " + id
	log.Debug().Msg(message)
	b, err := GetBookByID(id)
	if err != nil {
		log.Error().Err(err).Msg("")
		c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}
	c.JSON(http.StatusOK, b)

}

func CreateBook(c *gin.Context) {
	var req Book
	if err := c.BindJSON(&req); err == nil {
		req.ID = strconv.Itoa(len(Books) + 1)
		Books = append(Books, req)
		c.JSON(http.StatusCreated, req)
		return
	}
	err := errors.New("seems we have an error here")
	log.Error().Err(err).Msg("")
	c.JSON(http.StatusBadRequest, gin.H{"message": "Something went wrong"})

}

func UpdateBook(c *gin.Context) {
	var req Book
	id := c.Param("id")
	b, err := GetBookByID(id)
	if err != nil {
		log.Error().Err(err).Msg("")
		c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}
	if err := c.BindJSON(&req); err == nil {
		//lets user update a single field or multiple fields
		if req.ID != "" {
			b.ID = req.ID
		}
		if req.Title != "" {
			b.Title = req.Title
		}
		if req.Author != "" {
			b.Author = req.Author
		}
		if req.Pages != 0 {
			b.Pages = req.Pages
		}
		c.JSON(http.StatusOK, b)

	} else {
		log.Error().Err(err).Msg("")
		c.JSON(http.StatusBadRequest, err)
	}

}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := deleteBookByID(id); err != nil {
		log.Error().Err(err).Msg("")
		c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"message": "deleted"})
	}

}

func deleteBookByID(id string) error {
	var index int = -1
	for i, b := range Books {
		if b.ID == id {
			index = i
		}
	}
	if index != -1 {
		Books = append(Books[:index], Books[index+1:]...)
		return nil
	} else {
		return errors.New("book not found")
	}

}

func GetBookByID(id string) (*Book, error) {

	for i, b := range Books {
		if b.ID == id {
			return &Books[i], nil
		}
	}
	return nil, errors.New("book not found")
}
