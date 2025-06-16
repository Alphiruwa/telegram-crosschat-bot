package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken         string
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	PostgresDB       string
	PostgresSSLMode  string
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	return &Config{
		BotToken:         os.Getenv("BOT_TOKEN"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresDB:       os.Getenv("POSTGRES_NAME"),
		PostgresSSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
}
