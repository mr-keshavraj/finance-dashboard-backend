package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all environment-based configuration for the app
type Config struct {
	Port      string
	MongoURI  string
	DBName    string
	JWTSecret string
}

// AppConfig will be used globally across the application
var AppConfig Config

func LoadConfig() {
	// Try loading variables from .env file (useful for local development)
	err := godotenv.Load()
	if err != nil {
		// If .env is not found, fallback to system environment variables
		log.Println(".env file not found, using system environment variables")
	}

	// Initialize the global config struct with values or defaults
	AppConfig = Config{
		Port:      getEnv("PORT", "8080"),
		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:    getEnv("DB_NAME", "finance_dashboard"),
		JWTSecret: getEnv("JWT_SECRET", "my_super_secret_key"),
	}
}

// getEnv fetches environment variable or returns default if not set
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	// Return default value if env variable is empty
	if value == "" {
		return defaultValue
	}

	return value
}