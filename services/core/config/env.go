package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GetEnv func to get env values
func GetEnv(key string) string {
	// load .env file
	err := godotenv.Load("/Users/joel/Desktop/Coding/zero_fintech/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(key)
}
