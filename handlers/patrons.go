package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Patron struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	CheckedOutBooks []Book  `json:"checkedOutBooks"`
	FinesDue        float64 `json:"overdueFees"`
}

var Patrons = []*Patron{
	{ID: "1", Name: "Megan Buck"},
	{ID: "2", Name: "Riley Flynn"},
	{ID: "3", Name: "Indy Flynn"},
}

type FineReq struct {
	Fine float64 `json:"fine"`
}

func GetPatrons(c *gin.Context) {
	c.JSON(http.StatusOK, Patrons)
}

func GetPatron(c *gin.Context) {
	id := c.Param("id")
	p, err := getPatronByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, p)
}

func CreatePatron(c *gin.Context) {
	var req Patron
	//req should be an object with "name": <patron_name>
	if err := c.BindJSON(&req); err == nil {
		req.ID = strconv.Itoa(len(Patrons) + 1)
		Patrons = append(Patrons, &req)
		c.JSON(http.StatusAccepted, req)
	} else {
		c.JSON(http.StatusNotFound, err)
	}

}

type CheckoutReq struct {
	CheckedOutBooks []string `json:"checkedOutBooks"`
}

func CheckoutBooks(c *gin.Context) {
	pID := c.Param("id")
	p, err := getPatronByID(pID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	var req CheckoutReq
	// req should be an array of numerical strings that stand for the id of the books to check out
	if err := c.BindJSON(&req); err == nil {
		//loops through each id in the array, gets the Book object connected to it, and if it is not checked out yet, checks it out
		//and adds it to a slice of books on the patron's object
		for _, id := range req.CheckedOutBooks {
			book, err1 := GetBookByID(id)
			if err1 != nil {
				fmt.Println(err1)
			}
			if book.CheckedOut {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "book already checked out"})
				return

			}
			book.Patron = p.Name
			book.CheckedOut = true
			p.CheckedOutBooks = append(p.CheckedOutBooks, *book)

		}
		c.JSON(http.StatusAccepted, p)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
	}

}

func ReturnBooks(c *gin.Context) {
	pID := c.Param("id")
	p, err := getPatronByID(pID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	var req CheckoutReq
	if err := c.BindJSON(&req); err == nil {
		//goes through the array of ids passed into the request, gets the Book object connected to it
		//checks that it is checked out and to this patron, then returns it and takes it off the patron's list of books
		for _, bookID := range req.CheckedOutBooks {
			book, err1 := GetBookByID(bookID)
			if err1 != nil {
				fmt.Println(err1)
			}
			if !book.CheckedOut {
				message, _ := fmt.Scanf("%v wasn't checked out", book.Title)
				c.JSON(http.StatusConflict, gin.H{"message": message})
				return
			}
			if book.Patron != p.Name {
				message, _ := fmt.Scanf("%v didn't check out %v", p.Name, book.Title)
				c.JSON(http.StatusConflict, gin.H{"message": message})
				return
			}
			var index int = -1
			for j, b := range p.CheckedOutBooks {
				if book.ID == b.ID {
					book.CheckedOut = false
					book.Author = ""
					index = j
				}
			}
			if index != -1 {
				p.CheckedOutBooks = append(p.CheckedOutBooks[:index], p.CheckedOutBooks[index+1:]...)
			}

			c.JSON(http.StatusAccepted, p)
		}

	} else {
		c.JSON(http.StatusBadRequest, err)
	}

}

func UpdatePatron(c *gin.Context) {
	var req Patron
	//can update single field or multiple fields, but intended for just the name
	pID := c.Param("id")
	p, err := getPatronByID(pID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err})
		return
	}
	if err := c.BindJSON(&req); err == nil {
		if req.ID != "" {
			p.ID = req.ID
		}
		if req.Name != "" {
			p.Name = req.Name
		}
		if req.FinesDue != 0 {
			p.FinesDue = req.FinesDue
		}
		c.JSON(http.StatusOK, p)

	} else {
		c.JSON(http.StatusNotFound, err)
	}
}

func AddFine(c *gin.Context) {
	var req FineReq
	//float - amount of fine for patron
	pID := c.Param("id")
	p, err := getPatronByID(pID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	if err := c.BindJSON(&req); err == nil {
		//add fine to the amount the patron owes
		p.FinesDue = p.FinesDue + req.Fine
		c.JSON(http.StatusOK, p)

	} else {
		c.JSON(http.StatusBadRequest, err)
	}

}

func ReduceFine(c *gin.Context) {
	var req FineReq
	//float - amount the patron pays off of fine
	pID := c.Param("id")
	p, err := getPatronByID(pID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	if err := c.BindJSON(&req); err == nil {
		if req.Fine > p.FinesDue {
			//checks the patron is not overpaying, if they are, it 'returns' the change
			returnChange := req.Fine - p.FinesDue
			p.FinesDue = 0
			message, _ := fmt.Scanf("Return change: %v", returnChange)
			c.JSON(http.StatusAccepted, gin.H{"message": message})
		} else {
			//reduces the amount the patron owes by the amount the patron is paying off
			p.FinesDue = p.FinesDue - req.Fine
			c.JSON(http.StatusOK, p)
		}

	} else {
		c.JSON(http.StatusBadRequest, err)
	}

}

func DeletePatron(c *gin.Context) {
	id := c.Param("id")
	if err := deletePatronByID(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"message": "deleted"})
	}

}

func deletePatronByID(id string) error {
	var index int = -1
	for i, p := range Patrons {
		if p.ID == id {
			index = i
		}
	}
	if index != -1 {
		Patrons = append(Patrons[:index], Patrons[index+1:]...)
		return nil
	} else {
		return errors.New("Book Not Found")
	}

}

func getPatronByID(id string) (*Patron, error) {
	for i, patron := range Patrons {
		if patron.ID == id {
			return Patrons[i], nil
		}
	}
	return nil, errors.New("Patron not found")
}
