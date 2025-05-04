package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken    string
	PostgresURI string
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	return &Config{
		BotToken:    os.Getenv("BOT_TOKEN"),
		PostgresURI: os.Getenv("POSTGRES_URI"),
	}
}
