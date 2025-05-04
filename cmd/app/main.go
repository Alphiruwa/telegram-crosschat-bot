package main

import (
	"context"
	"log"

	"github.com/Alphiruwa/telegram-crosschat-bot/internal/bot"
	"github.com/Alphiruwa/telegram-crosschat-bot/internal/config"
	"github.com/Alphiruwa/telegram-crosschat-bot/internal/storage/postgresql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.MustLoad()

	db, err := pgxpool.New(context.Background(), cfg.PostgresURI)
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
