package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	// load .env
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	return os.Getenv(key)
}
