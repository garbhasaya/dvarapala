package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"dvarapala/internal/db"
	"dvarapala/internal/platform/auth"
	platformhttp "dvarapala/internal/platform/http"
	"dvarapala/internal/user"
)

// @title Dvarapala API
// @version 1.0
// @description This is a microservice for user management.
// @host localhost:8080
// @BasePath /

func main() {
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "logs"
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("failed to create log directory: %v\n", err)
		os.Exit(1)
	}

	logFile, err := os.OpenFile(filepath.Join(logDir, "api.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	logger := slog.New(slog.NewJSONHandler(mw, nil))
	slog.SetDefault(logger)

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "dvarapala.db"
	}

	client, err := db.NewSQLiteClient(dbPath)
	if err != nil {
		slog.Error("failed to open sqlite client", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	// Auth setup
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "very-secret-key" // Should be changed in production
	}
	jwtManager := auth.NewJWTManager(jwtSecret, 24*time.Hour)

	// Initialize components
	userRepo := user.NewRepository(client)
	userSvc := user.NewService(userRepo, jwtManager)
	userHandler := user.NewHandler(userSvc)

	router := platformhttp.NewRouter(userHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to listen and serve", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited gracefully")
}
