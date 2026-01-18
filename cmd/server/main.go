package main

import (
	"OpsGo/internal/application/service/devops"
	"OpsGo/internal/infrastructure/config"
	"OpsGo/internal/infrastructure/database"
	devops_repo "OpsGo/internal/infrastructure/repository/devops"
	devopsHandler "OpsGo/internal/interfaces/http/handler/devops"
	"OpsGo/internal/interfaces/http/middleware"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Warning: Could not load config, using defaults: %v", err)
	}

	// 2. Initialize Database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. Initialize Repositories
	devopsRepo := devops_repo.NewDevOpsRepository(db)

	// 4. Initialize Services
	devopsService := devops.NewDevOpsService(devopsRepo)

	// 5. Initialize Handlers
	devOpsH := devopsHandler.NewDevOpsHandler(devopsService)

	// 6. Setup Router
	r := gin.Default()
	r.Use(middleware.CORS()) // Ensure CORS is enabled for 8080/8081 cross-origin

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OpsGo is running"})
	})

	// API Routes
	v1 := r.Group("/api/v1")
	{
		// Public SSE Route
		v1.GET("/devops/events", devOpsH.StreamLogs)

		opsParams := v1.Group("/devops")
		// opsParams.Use(middleware.Auth()) // Optional: Add auth if needed
		{
			opsParams.POST("/config", devOpsH.ConfigRepo)
			opsParams.GET("/summary", devOpsH.GetSummary)
			opsParams.POST("/deploy", devOpsH.TriggerDeployment)
		}

		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("/github", devOpsH.HandleGitHubWebhook)
		}
	}

	// 7. Start Server on 8081
	port := "8081"
	fmt.Printf("OpsGo starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
