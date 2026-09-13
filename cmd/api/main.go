package main

import (
	"fmt"
	"log"

	"github.com/mhdalfrzh/stocklab-backend/internal/config"
	"github.com/mhdalfrzh/stocklab-backend/internal/portfolio"
	"github.com/mhdalfrzh/stocklab-backend/internal/yahoofinance"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	// Setup Database
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Dependency Injection
	portfolioRepo := portfolio.NewRepository(db)
	yahooClient := yahoofinance.NewClient()
	portfolioService := portfolio.NewService(portfolioRepo, yahooClient)
	portfolioHandler := portfolio.NewHandler(portfolioService)

	// Setup Gin Router
	router := gin.Default()

	// Register Routes
	api := router.Group("/api")
	{
		api.GET("/portfolios/dashboard", portfolioHandler.GetDashboard)
	}

	// Start Server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("server listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
