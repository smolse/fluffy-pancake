package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/smolse/fluffy-pancake/internal/config"
	"github.com/smolse/fluffy-pancake/internal/datastores"
	"github.com/smolse/fluffy-pancake/internal/router"
	"github.com/smolse/fluffy-pancake/internal/service"
)

func main() {
	// Load the application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize and connect to the data store
	db, err := datastores.NewDataStore(&cfg.DataStore)
	if err != nil {
		log.Fatalf("Failed to create data store: %v", err)
	}
	err = db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to data store: %v", err)
	}
	defer db.Close()

	// Initialize the risk service
	svc := service.NewRiskService(db)

	// Update the Gin mode
	gin.SetMode(cfg.Gin.Mode)

	// Initialize the HTTP server and launch it in a goroutine so that it won't block the graceful shutdown handling
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router.NewRouter(svc),
	}
	go func() {
		log.Printf("Starting server on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Gracefully handle SIGINT and SIGTERM signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server is forced to shutdown due to an error: ", err)
	}

	log.Println("Server has been shutdown gracefully")
}
