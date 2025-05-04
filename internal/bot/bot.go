package bot

import (
	"log"
	"strings"

	"github.com/Alphiruwa/telegram-crosschat-bot/internal/entity"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	cmdShare  = "share"
	cmdList   = "listlinks"
	cmdLink   = "linkchatid"
	cmdUnlink = "unlinkchatid"
)

type Bot struct {
	API         *tgbotapi.BotAPI
	LinkRepo    entity.LinkRepository
	RequestRepo entity.RequestRepository
}

func New(api *tgbotapi.BotAPI, linkRepo entity.LinkRepository, requestRepo entity.RequestRepository) *Bot {
	return &Bot{api, linkRepo, requestRepo}
}

func (b *Bot) SetCommands() error {
	_, err := b.API.Request(tgbotapi.NewSetMyCommandsWithScope(
		tgbotapi.BotCommandScope{Type: "all_chat_administrators"},
		tgbotapi.BotCommand{
			Command:     cmdShare,
			Description: "получить команду для связи с этим чатом",
		},
		tgbotapi.BotCommand{
			Command:     cmdList,
			Description: "показать связанные чаты",
		},
		tgbotapi.BotCommand{
			Command:     cmdLink,
			Description: "отправить запрос на связь с другим чатом",
		},
		tgbotapi.BotCommand{
			Command:     cmdUnlink,
			Description: "отменить связь с другим чатом",
		},
	))
	return err
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.API.GetUpdatesChan(u)
	for update := range updates {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println("Panic recovered:", err)
				}
			}()
			switch {
			case update.Message != nil:
				switch {
				case update.Message.Command() != "":
					switch update.Message.Command() {
					case cmdShare:
						b.OnlyAdminCommand(b.Share, &update)
					case cmdList:
						b.OnlyAdminCommand(b.List, &update)
					case cmdLink:
						b.OnlyAdminCommand(b.Link, &update)
					case cmdUnlink:
						b.OnlyAdminCommand(b.Unlink, &update)
					}
				default:
					b.ForwardMessage(&update)
				}
			case update.CallbackData() != "":
				switch {
				case strings.HasPrefix(update.CallbackData(), "a"):
					b.OnlyAdminButton(b.AcceptRequestBtn, &update)
				case strings.HasPrefix(update.CallbackData(), "d"):
					b.OnlyAdminButton(b.DeclineRequestBtn, &update)
				}
			}
		}()
	}
}
