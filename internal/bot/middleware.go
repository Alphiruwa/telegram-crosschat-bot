package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *Bot) IsChatAdmin(userID int64, chatConfig tgbotapi.ChatConfig) bool {
	admins, _ := b.API.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{ChatConfig: chatConfig})
	for _, admin := range admins {
		if admin.User.ID == userID {
			return true
		}
	}
	return false
}

func (b *Bot) OnlyAdminCommand(f func(*tgbotapi.Update), u *tgbotapi.Update) {
	if !b.IsChatAdmin(u.Message.From.ID, u.Message.Chat.ChatConfig()) {
		msg := tgbotapi.NewMessage(
			u.Message.Chat.ID,
			"Эта команда доступна только администраторам чата.")
		msg.ReplyToMessageID = u.Message.MessageID
		b.API.Send(msg)
		return
	}
	f(u)
}

func (b *Bot) OnlyAdminButton(f func(*tgbotapi.Update), u *tgbotapi.Update) {
	if !b.IsChatAdmin(u.CallbackQuery.From.ID, u.CallbackQuery.Message.Chat.ChatConfig()) {
		b.API.Request(tgbotapi.NewCallbackWithAlert(
			u.CallbackQuery.ID,
			"Эта кнопка доступна только администраторам чата",
		))
		return
	}
	f(u)
}
