package main

import (
	"OpsGo/internal/application/service/devops"
	"OpsGo/internal/application/service/monitor"
	"OpsGo/internal/infrastructure/config"
	"OpsGo/internal/infrastructure/database"
	"OpsGo/internal/infrastructure/redis"
	devops_repo "OpsGo/internal/infrastructure/repository/devops"
	devopsHandler "OpsGo/internal/interfaces/http/handler/devops"
	monitorHandler "OpsGo/internal/interfaces/http/handler/monitor"
	"OpsGo/internal/interfaces/http/middleware"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Config
	// Passing empty string uses default config.yaml or APP_ENV
	if err := config.LoadConfig(""); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB()

	// 2.1 Initialize Redis
	if err := redis.InitRedis(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Monitor feature will be disabled.", err)
	} else {
		defer redis.CloseRedis()
	}

	// 3. Initialize Repositories
	devopsRepo := devops_repo.NewDevOpsRepository(database.DB)

	// 4. Initialize Services
	devopsService := devops.NewDevOpsService(devopsRepo)

	// Monitor Service
	monitorService := monitor.NewMonitorService()
	if redis.Client != nil {
		monitorService.StartCollector()
		defer monitorService.StopCollector()
	}

	// 5. Initialize Handlers
	devOpsH := devopsHandler.NewDevOpsHandler(devopsService)
	monitorH := monitorHandler.NewMonitorHandler(monitorService)

	// 6. Setup Router
	r := gin.Default()
	r.Use(middleware.CORS()) // Ensure CORS is enabled for 8080/8081 cross-origin

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OpsGo is running"})
	})

	// API Routes
	v1 := r.Group("/api/v1/devops")
	{
		// Public SSE Route
		v1.GET("/events", devOpsH.StreamLogs)

		v1.POST("/config", devOpsH.ConfigRepo)
		v1.DELETE("/config/:id", devOpsH.DeleteConfig)
		v1.GET("/summary", devOpsH.GetSummary)
		v1.POST("/deploy", devOpsH.TriggerDeployment)
		v1.GET("/logs/:id", devOpsH.GetServiceLog)

		v1.GET("/monitor/stats", monitorH.GetStats)

		v1.POST("/webhooks/github", devOpsH.HandleGitHubWebhook)
		v1.POST("/webhooks/ci", devOpsH.HandleCICallback)
	}

	// 7. Start Server on 8081
	// We use 8081 specifically for OpsGo to avoid conflict with FlowGo (8080)
	port := "8081"
	fmt.Printf("OpsGo starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
