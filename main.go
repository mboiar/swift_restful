package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	_, q, err = repository.SetupDB(nil)
	if err != nil {
		log.Fatalf("could not setup database: %v", err)
	}
	SwiftController = *controllers.NewSwiftController(q, ctx)
	SwiftRoutes = routes.NewRouteSwift(SwiftController)

	server = gin.Default()
	server.SetTrustedProxies(nil)
}

func main() {

	slog.Info(fmt.Sprintf("WTFFF %v", server))
	router := server.Group("/")
	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "API is working"})
	})
	SwiftRoutes.SwiftRoute(router)
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	})

	srv := &http.Server{
		Addr:    ":8000",
		Handler: server.Handler(),
	}

	go func() {
		slog.Info("Starting server on :8000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	<-ctx.Done()
	log.Println("timeout of 5 seconds.")
	log.Println("Server exiting")
}
