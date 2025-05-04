package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) ForwardMessage(u *tgbotapi.Update) {
	links, err := b.LinkRepo.GetAllChatLinks(u.Message.Chat.ID)
	if err != nil {
		return
	}
	if len(links) == 0 {
		return
	}
	sendMsg := tgbotapi.NewMessage(
		0,
		fmt.Sprintf("🔗 <b>%s</b> (<code>%d</code>)", u.Message.Chat.Title, u.Message.Chat.ID),
	)
	sendMsg.ParseMode = "HTML"
	forwardMsg := tgbotapi.NewForward(0, u.Message.Chat.ID, u.Message.MessageID)
	for _, link := range links {
		toChatID := link.TgtChatID
		if link.TgtChatID == u.Message.Chat.ID {
			toChatID = link.SrcChatID
		}
		sendMsg.ChatID = toChatID
		b.API.Send(sendMsg)
		forwardMsg.ChatID = toChatID
		b.API.Request(forwardMsg)
	}
}
