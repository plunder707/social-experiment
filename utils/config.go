// utils/config.go
package utils

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

// Config holds all configuration variables
type Config struct {
	MongoURI        string
	JWTSecret       string
	ServerPort      string
	RateLimit       rate.Limit
	RateBurst       int
	CORSOrigins     []string
	SecurityHeaders bool
}

// LoadConfig loads environment variables and returns a Config struct
func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	config := Config{
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017/maliaki"),
		JWTSecret:       getEnv("JWT_SECRET", "your_jwt_secret"),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		RateLimit:       getEnvAsRateLimit("RATE_LIMIT", 10),
		RateBurst:       getEnvAsInt("RATE_BURST", 20),
		CORSOrigins:     splitEnv("CORS_ORIGINS", ","),
		SecurityHeaders: getEnvAsBool("SECURITY_HEADERS", true),
	}

	return config
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	if valueStr, exists := os.LookupEnv(name); exists {
		var value int
		_, err := strconv.Sscanf(valueStr, "%d", &value)
		if err != nil {
			log.Printf("Invalid integer for %s, using default %d", name, defaultVal)
			return defaultVal
		}
		return value
	}
	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	if valueStr, exists := os.LookupEnv(name); exists {
		value, err := strconv.ParseBool(valueStr)
		if err != nil {
			log.Printf("Invalid boolean for %s, using default %v", name, defaultVal)
			return defaultVal
		}
		return value
	}
	return defaultVal
}

func getEnvAsRateLimit(name string, defaultVal rate.Limit) rate.Limit {
	if valueStr, exists := os.LookupEnv(name); exists {
		var value float64
		_, err := strconv.Sscanf(valueStr, "%f", &value)
		if err != nil {
			log.Printf("Invalid rate limit for %s, using default %v", name, defaultVal)
			return defaultVal
		}
		return rate.Limit(value)
	}
	return defaultVal
}

func splitEnv(name string, sep string) []string {
	if value, exists := os.LookupEnv(name); exists {
		parts := strings.Split(value, sep)
		var trimmed []string
		for _, part := range parts {
			trimmed = append(trimmed, strings.TrimSpace(part))
		}
		return trimmed
	}
	return []string{}
}
