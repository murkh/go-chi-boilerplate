package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/murkh/gig-mobile-backend/internal/config"
	"github.com/murkh/gig-mobile-backend/internal/core/domain"
	"github.com/murkh/gig-mobile-backend/internal/handlers"
	"github.com/murkh/gig-mobile-backend/internal/platform/database"
	"github.com/murkh/gig-mobile-backend/internal/platform/logger"
	"github.com/murkh/gig-mobile-backend/internal/repository/postgres"
	"github.com/murkh/gig-mobile-backend/internal/service"
	"go.uber.org/zap"
)

func main() {
	// Load Config
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize Logger
	log, err := logger.New(cfg.Logger.Level, cfg.Server.Mode)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Database Setup
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto Migration
	if err := db.AutoMigrate(&domain.User{}, &domain.Post{}); err != nil {
		log.Fatal("Failed to migrate database", zap.Error(err))
	}

	// Dependency Injection
	userRepo := postgres.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Router Setup
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Mount Routes
	userHandler.RegisterRoutes(r)

	// Start Server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		log.Info("Starting server", zap.String("port", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// TODO: Add context with timeout for shutdown
}
