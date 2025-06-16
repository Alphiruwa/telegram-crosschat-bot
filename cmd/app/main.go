package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Alphiruwa/telegram-crosschat-bot/internal/bot"
	"github.com/Alphiruwa/telegram-crosschat-bot/internal/config"
	"github.com/Alphiruwa/telegram-crosschat-bot/internal/storage/postgresql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.MustLoad()
	
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB, cfg.PostgresSSLMode)
	if m, err := migrate.New("file://migrations/postgresql", dbURI); err != nil {
		log.Fatal(err)
	} else if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	db, err := pgxpool.New(context.Background(), dbURI)
	if err != nil {
		log.Fatal(err)
	}
	linkRepo := postgresql.NewLinkRepository(db)
	requestRepo := postgresql.NewRequestRepository(db)

	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Authorized as @%s", api.Self.UserName)
	}

	bot := bot.New(api, linkRepo, requestRepo)
	if err := bot.SetCommands(); err != nil {
		log.Fatal("Failed to set bot commands:", err)
	}
	bot.Run()
}
