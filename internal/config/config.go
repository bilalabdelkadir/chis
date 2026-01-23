package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DbUrl     string
	JwtSecret string
}

func LoadEnv() (*Config, error) {
	godotenv.Load()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		return nil, errors.New("DB_URL not found.")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET not found.")
	}

	return &Config{
		Port:      port,
		DbUrl:     dbUrl,
		JwtSecret: jwtSecret,
	}, nil

}
