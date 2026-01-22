package config

import (
	"errors"
	"os"
)

type Config struct {
	Port  string
	DbUrl string
}

func LoadEnv() (*Config, error) {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		return nil, errors.New("db url not found.")
	}

	return &Config{
		Port:  port,
		DbUrl: dbUrl,
	}, nil

}
