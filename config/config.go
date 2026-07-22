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

// LoadConfig membaca file .env dan mengisi struct Config.
// Di dalam container tidak ada file .env (env sudah di-inject compose lewat
// env_file), jadi file yang tidak ditemukan bukan error.
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
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
