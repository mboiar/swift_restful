package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Define routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Ok",
		})
	})

	// Start the server
	router.Run(":8000")
}
