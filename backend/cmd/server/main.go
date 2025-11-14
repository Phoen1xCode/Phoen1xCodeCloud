package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/phoen1xcode/phoen1xcodecloud/internal/config"
	"github.com/phoen1xcode/phoen1xcodecloud/internal/handlers"
	"github.com/phoen1xcode/phoen1xcodecloud/internal/middleware"
	"github.com/phoen1xcode/phoen1xcodecloud/internal/models"
	"github.com/phoen1xcode/phoen1xcodecloud/pkg/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Share{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	var stor storage.Storage
	if cfg.StorageType == "s3" {
		stor, err = storage.NewS3Storage(cfg.S3Endpoint, cfg.S3Region, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Bucket)
	} else {
		stor, err = storage.NewLocalStorage(cfg.LocalStoragePath)
	}
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}

	r := gin.Default()
	
	// Configure CORS with specific allowed origins
	allowedOrigins := []string{"http://localhost:3000", "http://localhost:5173"}
	if envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); envOrigins != "" {
		allowedOrigins = strings.Split(envOrigins, ",")
	}
	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}))
	
	// Add rate limiting: 100 requests per minute per IP
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	r.Use(rateLimiter.Middleware())

	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
	shareHandler := handlers.NewShareHandler(db, stor)
	adminHandler := handlers.NewAdminHandler(db)

	api := r.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		api.GET("/share/:code", shareHandler.GetShare)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			auth.POST("/upload", shareHandler.UploadFile)
			auth.POST("/text", shareHandler.CreateTextShare)
			auth.GET("/shares", shareHandler.ListUserShares)
			auth.DELETE("/share/:code", shareHandler.DeleteShare)

			admin := auth.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/stats", adminHandler.GetStats)
				admin.GET("/shares", adminHandler.ListAllShares)
				admin.GET("/users", adminHandler.ListUsers)
			}
		}
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
