package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database     DatabaseConfig
	Redis        RedisConfig
	JWT          JWTConfig
	Stripe       StripeConfig
	Server       ServerConfig
	OpenSSL      OpenSSLConfig
	RateLimit    RateLimitConfig
	CORS         CORSConfig
	FileUpload   FileUploadConfig
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

type StripeConfig struct {
	SecretKey     string
	WebhookSecret string
}

type ServerConfig struct {
	Port string
	Env  string
}

type OpenSSLConfig struct {
	BinaryPath string
}

type RateLimitConfig struct {
	Requests int
	Window   time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
}

type FileUploadConfig struct {
	MaxFileSize int64
	UploadDir   string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config := &Config{
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/openssl_ui?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379"),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "change-this-in-production"),
			ExpiresIn: parseDuration(getEnv("JWT_EXPIRES_IN", "24h")),
		},
		Stripe: StripeConfig{
			SecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
			WebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		OpenSSL: OpenSSLConfig{
			BinaryPath: getEnv("OPENSSL_BINARY_PATH", "openssl"),
		},
		RateLimit: RateLimitConfig{
			Requests: parseInt(getEnv("RATE_LIMIT_REQUESTS", "1000")),
			Window:   parseDuration(getEnv("RATE_LIMIT_WINDOW", "1h")),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		},
		FileUpload: FileUploadConfig{
			MaxFileSize: parseSize(getEnv("MAX_FILE_SIZE", "10MB")),
			UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
		},
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("Error parsing integer %s: %v", s, err)
		return 0
	}
	return i
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("Error parsing duration %s: %v", s, err)
		return time.Hour
	}
	return d
}

func parseSize(s string) int64 {
	// Simple size parser for MB/GB
	s = strings.TrimSpace(strings.ToUpper(s))

	if strings.HasSuffix(s, "MB") {
		size := parseInt(strings.TrimSuffix(s, "MB"))
		return int64(size) * 1024 * 1024
	}

	if strings.HasSuffix(s, "GB") {
		size := parseInt(strings.TrimSuffix(s, "GB"))
		return int64(size) * 1024 * 1024 * 1024
	}

	// Default to bytes
	return int64(parseInt(s))
}