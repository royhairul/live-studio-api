package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBPort       string
	ServerPort   string
	ShopeeAPIKey string
	JWTToken     string
	JWTExpiredAt string
}

// LoadConfig mengisi struct Config dari environment variable.
// File .env bersifat opsional: di dalam container, environment variable
// di-inject langsung oleh Docker/compose sehingga file .env tidak ada.
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading configuration from environment variables")
	}

	return &Config{
		DBHost:      os.Getenv("DB_HOST"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DBPort:      os.Getenv("DB_PORT"),
		ServerPort:  os.Getenv("SERVER_PORT"),
		ShopeeAPIKey: os.Getenv("SHOPEE_API_KEY"),
		JWTToken: os.Getenv("JWT_SECRET"),
		JWTExpiredAt: os.Getenv("JWT_EXPIRED_AT"),
	}
}
