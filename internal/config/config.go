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
	RedisUrl  string
	GrpcAddr  string
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

	redisUlr := os.Getenv("REDIS_URL")
	if redisUlr == "" {
		return nil, errors.New("REDIS_URL not found.")
	}

	grpcAddr := os.Getenv("DELIVERY_GRPC_ADDR")

	return &Config{
		Port:      port,
		DbUrl:     dbUrl,
		JwtSecret: jwtSecret,
		RedisUrl:  redisUlr,
		GrpcAddr:  grpcAddr,
	}, nil

}
