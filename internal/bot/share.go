package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) Share(u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(
		u.Message.Chat.ID,
		fmt.Sprintf(
			"Чтобы связать другой чат с вашим, его администратор должен отправить туда следующую команду: <code>/%s %d</code>.",
			cmdLink, u.Message.Chat.ID,
		),
	)
	msg.ReplyToMessageID = u.Message.MessageID
	msg.ParseMode = "HTML"
	b.API.Send(msg)
}
