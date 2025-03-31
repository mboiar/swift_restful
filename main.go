package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	routes "swift-restful/api/v1"
	"swift-restful/controllers"
	"swift-restful/repository"
	dbCon "swift-restful/repository/sqlc"

	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine
	q      *dbCon.Queries
	ctx    context.Context

	SwiftController controllers.SwiftController
	SwiftRoutes     routes.SwiftRoutes
)

func init() {
	ctx = context.TODO()
	var err error
	_, q, err = repository.SetupDB("env/development.env")
	if err != nil {
		log.Fatalf("could not setup database: %v", err)
	}
	SwiftController = *controllers.NewSwiftController(q, ctx)
	SwiftRoutes = routes.NewRouteSwift(SwiftController)

	server := gin.Default()
	server.SetTrustedProxies(nil)
}

func main() {
	router := server.Group("/api")
	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "API is working"})
	})
	SwiftRoutes.SwiftRoute(router)
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	})
	log.Fatal(server.Run(":8000"))
}
