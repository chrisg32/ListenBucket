package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	listenbucket "github.com/chrisg32/ListenBucket"
	"github.com/chrisg32/ListenBucket/internal/api"
	"github.com/chrisg32/ListenBucket/internal/config"
	"github.com/chrisg32/ListenBucket/internal/database"
	"github.com/chrisg32/ListenBucket/internal/downloader"
)

func main() {
	cfg := config.Load()

	// Ensure data directory exists
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Ensure media directory exists
	if err := os.MkdirAll(cfg.MediaDir, 0755); err != nil {
		log.Fatalf("Failed to create media directory: %v", err)
	}

	// Initialize database
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize downloader
	dl := downloader.New(db, cfg.MediaDir, cfg.BaseURL)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start background workers
	go dl.Start(ctx)
	go dl.StartSourceChecker(ctx)

	// Initialize API server
	server := api.NewServer(db, dl, cfg)
	handler := server.SetupRoutes(listenbucket.WebFS)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		fmt.Printf("🎧 ListenBucket is running at http://localhost:%s\n", cfg.Port)
		fmt.Printf("   Data directory: %s\n", cfg.DataDir)
		fmt.Printf("   Base URL: %s\n", cfg.BaseURL)
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
	cancel()

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	fmt.Println("Goodbye!")
}
