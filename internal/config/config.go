package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DbUrl        string
	JwtSecret    string
	RedisUrl     string
	GrpcAddr     string
	Env          string
	Origin       string
	ResendApiKey string
	AppUrl       string
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

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	origin := os.Getenv("CORS_ORIGIN")
	if origin == "" {
		origin = "localhost:8080"
	}

	resendApiKey := os.Getenv("RESEND_API_KEY")

	appUrl := os.Getenv("APP_URL")
	if appUrl == "" {
		appUrl = "http://localhost:3000"
	}

	return &Config{
		Port:         port,
		DbUrl:        dbUrl,
		JwtSecret:    jwtSecret,
		RedisUrl:     redisUlr,
		GrpcAddr:     grpcAddr,
		Env:          env,
		Origin:       origin,
		ResendApiKey: resendApiKey,
		AppUrl:       appUrl,
	}, nil

}
