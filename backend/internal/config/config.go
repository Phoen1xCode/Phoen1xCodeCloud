package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort      string
	DatabaseURL     string
	JWTSecret       string
	StorageType     string // "local" or "s3"
	LocalStoragePath string
	S3Bucket        string
	S3Region        string
	S3Endpoint      string
	S3AccessKey     string
	S3SecretKey     string
}

func Load() *Config {
	cfg := &Config{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/phoen1xcloud?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		StorageType:      getEnv("STORAGE_TYPE", "local"),
		LocalStoragePath: getEnv("LOCAL_STORAGE_PATH", "./uploads"),
		S3Bucket:         getEnv("S3_BUCKET", ""),
		S3Region:         getEnv("S3_REGION", "auto"),
		S3Endpoint:       getEnv("S3_ENDPOINT", ""),
		S3AccessKey:      getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey:      getEnv("S3_SECRET_KEY", ""),
	}
	
	// Validate JWT secret - must be set and strong
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable must be set")
	}
	if cfg.JWTSecret == "your-secret-key" || len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be a strong secret (at least 32 characters) and not the default value")
	}
	
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
