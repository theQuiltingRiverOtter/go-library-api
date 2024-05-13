package router

import (
	"library-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	//Set up mode for debug, release or test
	switch mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	}
	//create a router
	r := gin.Default()

	//GET methods
	r.GET("/books", handlers.GetBooks)

	r.GET("/books/:id", handlers.GetBook)

	r.GET("/patrons", handlers.GetPatrons)

	r.GET("/patrons/:id", handlers.GetPatron)

	//POST methods
	r.POST("/patrons", handlers.CreatePatron)

	r.POST("/books", handlers.CreateBook)

	//PATCH methods

	r.PATCH("/patrons/:id/checkout", handlers.CheckoutBooks)

	r.PATCH("/patrons/:id/return", handlers.ReturnBooks)

	r.PATCH("/patrons/:id/addfine", handlers.AddFine)

	r.PATCH("/patrons/:id/reducefine", handlers.ReduceFine)

	//PUT methods

	r.PUT("/patrons/:id/update", handlers.UpdatePatron)

	r.PUT("/books/:id/update", handlers.UpdateBook)

	//DELETE methods

	r.DELETE("/patrons/:id", handlers.DeletePatron)

	r.DELETE("/books/:id", handlers.DeleteBook)

	return r
}
