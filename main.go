package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Ping test
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	return r
}

func main() {
	// db, queries, err := repository.SetupDB()
	// if err != nil {
	// 	log.Fatal(err)
	// 	os.Exit(1)
	// }
	r := setupRouter()

	r.Run(":8000")
}
