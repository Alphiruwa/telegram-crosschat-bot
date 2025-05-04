package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Alphiruwa/telegram-crosschat-bot/internal/entity"
)

func (b *Bot) List(u *tgbotapi.Update) {
	links, err := b.LinkRepo.GetAllChatLinks(u.Message.Chat.ID)
	if err != nil {
		return
	}
	outRequests, err := b.RequestRepo.GetAllChatOutRequests(u.Message.Chat.ID)
	if err != nil {
		return
	}
	incRequests, err := b.RequestRepo.GetAllChatIncRequests(u.Message.Chat.ID)
	if err != nil {
		return
	}
	sb := strings.Builder{}
	sb.WriteString("<b>Связанные чаты</b>")
	if len(links) == 0 {
		sb.WriteString("\nнет")
	} else {
		for i, link := range links {
			linkedChatID := link.SrcChatID
			if link.SrcChatID == u.Message.Chat.ID {
				linkedChatID = link.TgtChatID
			}
			chat, err := b.API.GetChat(tgbotapi.ChatInfoConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: linkedChatID,
				},
			})
			if err != nil {
				sb.WriteString(fmt.Sprintf("\n%d. <b>?</b> (<code>%d</code>)", i+1, linkedChatID))
			} else {
				sb.WriteString(fmt.Sprintf("\n%d. <b>%s</b> (<code>%d</code>)", i+1, chat.Title, linkedChatID))
			}
		}
	}
	sb.WriteString("\n\n<b>Входящие запросы</b>")
	if len(incRequests) == 0 {
		sb.WriteString("\nнет")
	} else {
		for i, req := range incRequests {
			chat, err := b.API.GetChat(tgbotapi.ChatInfoConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: req.SrcChatID,
				},
			})
			if err != nil {
				sb.WriteString(fmt.Sprintf("\n%d. <b>?</b> (<code>%d</code>)", i+1, req.SrcChatID))
			} else {
				sb.WriteString(fmt.Sprintf("\n%d. <b>%s</b> (<code>%d</code>)", i+1, chat.Title, req.SrcChatID))
			}
		}
	}
	sb.WriteString("\n\n<b>Исходящие запросы</b>")
	if len(outRequests) == 0 {
		sb.WriteString("\nнет")
	} else {
		for i, req := range outRequests {
			sb.WriteString(fmt.Sprintf("\n%d. <code>%d</code>", i+1, req.TgtChatID))
		}
	}
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, sb.String())
	msg.ReplyToMessageID = u.Message.MessageID
	msg.ParseMode = "HTML"
	b.API.Send(msg)
}

func (b *Bot) AcceptRequest(srcChatID, tgtChatID int64, tgtMessageID int, fromUserName string) {
	err := b.RequestRepo.DeleteRequest(srcChatID, tgtChatID)
	if err != nil && err != entity.ErrRequestNotFound {
		return
	}
	delTgtMsg := tgbotapi.NewDeleteMessage(tgtChatID, tgtMessageID)
	b.API.Request(delTgtMsg)
	if err == entity.ErrRequestNotFound {
		return
	}
	if err := b.LinkRepo.CreateLink(srcChatID, tgtChatID); err != nil {
		return
	}
	tgtChat, err := b.API.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: tgtChatID,
		},
	})
	srcMsg := tgbotapi.NewMessage(
		srcChatID,
		"",
	)
	if err != nil {
		srcMsg.Text = fmt.Sprintf(
			"Запрос на связь к чату <code>%d</code> принят его администраторами. Сообщения оттуда будут дублироваться в ваш чат, а ваши сообщения — туда.",
			tgtChatID,
		)
	} else {
		srcMsg.Text = fmt.Sprintf(
			"Запрос на связь к чату <b>%s</b> (<code>%d</code>) принят его администраторами. Сообщения оттуда будут дублироваться в ваш чат, а ваши сообщения — туда.",
			tgtChat.Title, tgtChat.ID,
		)
	}
	srcMsg.ParseMode = "HTML"
	tgtMsg := tgbotapi.NewMessage(tgtChatID, "")
	if respSrcMsg, err := b.API.Send(srcMsg); err != nil {
		tgtMsg.Text = fmt.Sprintf(
			"Запрос на связь от чата <code>%d</code> принят @%s. Сообщения оттуда будут дублироваться в ваш чат, а ваши сообщения — туда.",
			srcChatID, fromUserName,
		)
	} else {
		tgtMsg.Text = fmt.Sprintf(
			"Запрос на связь от чата <b>%s</b> (<code>%d</code>) принят @%s. Сообщения оттуда будут дублироваться в ваш чат, а ваши сообщения — туда.",
			respSrcMsg.Chat.Title, respSrcMsg.Chat.ID, fromUserName,
		)
	}
	tgtMsg.ParseMode = "HTML"
	b.API.Send(tgtMsg)
}

func (b *Bot) Link(u *tgbotapi.Update) {
	curMsg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
	curMsg.ReplyToMessageID = u.Message.MessageID
	tgtChatID, err := strconv.ParseInt(u.Message.CommandArguments(), 10, 64)
	if err != nil {
		curMsg.Text = "Добавьте к команде ID чата, с которым хотите связать ваш."
		b.API.Send(curMsg)
		return
	}
	if tgtChatID == u.Message.Chat.ID {
		curMsg.Text = "Добавьте к команде ID того чата, с которым хотите связать ваш."
		b.API.Send(curMsg)
		return
	}
	if exists, err := b.LinkRepo.IsLinkExists(u.Message.Chat.ID, tgtChatID); err != nil {
		return
	} else if exists {
		curMsg.Text = "Указанный чат уже связан с вашим."
		b.API.Send(curMsg)
		return
	}
	if exists, err := b.RequestRepo.IsRequestExists(u.Message.Chat.ID, tgtChatID); err != nil {
		return
	} else if exists {
		curMsg.Text = "Запрос на связь с указанным чатом уже отправлен."
		b.API.Send(curMsg)
		return
	}
	if req, err := b.RequestRepo.GetRequest(tgtChatID, u.Message.Chat.ID); err != nil && err != entity.ErrRequestNotFound {
		return
	} else if err == nil {
		b.AcceptRequest(req.SrcChatID, req.TgtChatID, int(req.TgtMessageID), u.Message.From.UserName)
		return
	}
	tgtMsg := tgbotapi.NewMessage(
		tgtChatID,
		fmt.Sprintf(
			"@%s предлагает связать ваш чат со своим: <b>%s</b> (<code>%d</code>). Если принять предложение, сообщения оттуда будут дублироваться в ваш чат, а ваши сообщения — туда.",
			u.Message.From.UserName, u.Message.Chat.Title, u.Message.Chat.ID,
		),
	)
	tgtMsg.ParseMode = "HTML"
	curChatIDStr := strconv.Itoa(int(u.Message.Chat.ID))
	acceptCb := fmt.Sprintf("a%s", curChatIDStr)
	declineCb := fmt.Sprintf("d%s", curChatIDStr)
	tgtMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:         "✅ Принять",
				CallbackData: &acceptCb,
			}, tgbotapi.InlineKeyboardButton{
				Text:         "❌ Отклонить",
				CallbackData: &declineCb,
			},
		),
	)
	respTgtMsg, err := b.API.Send(tgtMsg)
	if err != nil {
		curMsg.Text = fmt.Sprintf("Не удалось отправить запрос на связь с чатом <code>%d</code>. Скорее всего, у меня недостаточно прав в чате-получателе.", tgtChatID)
		curMsg.ParseMode = "HTML"
		b.API.Send(curMsg)
		b.RequestRepo.DeleteRequest(u.Message.Chat.ID, tgtChatID)
		return
	}
	if err = b.RequestRepo.CreateRequest(u.Message.Chat.ID, tgtChatID, int64(respTgtMsg.MessageID), u.Message.From.ID); err != nil {
		return
	}
	curMsg.Text = fmt.Sprintf("Запрос на связь с чатом <code>%d</code> отправлен его администраторам.", tgtChatID)
	curMsg.ParseMode = "HTML"
	b.API.Send(curMsg)
}

func (b *Bot) AcceptRequestBtn(u *tgbotapi.Update) {
	srcChatID, _ := strconv.ParseInt(strings.TrimPrefix(u.CallbackData(), "a"), 10, 64)
	b.AcceptRequest(srcChatID, u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID, u.CallbackQuery.From.UserName)
}

func (b *Bot) DeclineRequest(srcChatID, tgtChatID int64, tgtMessageID int, fromUserName string) {
	err := b.RequestRepo.DeleteRequest(srcChatID, tgtChatID)
	if err != nil && err != entity.ErrRequestNotFound {
		return
	}
	delTgtMsg := tgbotapi.NewDeleteMessage(tgtChatID, tgtMessageID)
	b.API.Request(delTgtMsg)
	if err == entity.ErrRequestNotFound {
		return
	}
	srcMsg := tgbotapi.NewMessage(
		srcChatID,
		fmt.Sprintf("Запрос на связь к чату <code>%d</code> был отклонен его администраторами.", tgtChatID),
	)
	srcMsg.ParseMode = "HTML"
	b.API.Send(srcMsg)
	tgtMsg := tgbotapi.NewMessage(
		tgtChatID,
		fmt.Sprintf("Запрос на связь от чата <code>%d</code> был отменен @%s.", srcChatID, fromUserName),
	)
	tgtMsg.ParseMode = "HTML"
	b.API.Send(tgtMsg)
}

func (b *Bot) Unlink(u *tgbotapi.Update) {
	curMsg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
	curMsg.ReplyToMessageID = u.Message.MessageID
	tgtChatID, err := strconv.ParseInt(u.Message.CommandArguments(), 10, 64)
	if err != nil {
		curMsg.Text = "Добавьте к команде ID чата, связь с которым хотите удалить."
		b.API.Send(curMsg)
		return
	}
	if err := b.LinkRepo.DeleteLink(u.Message.Chat.ID, tgtChatID); err == nil {
		tgtMsg := tgbotapi.NewMessage(
			tgtChatID,
			fmt.Sprintf(
				"Связь с чатом <b>%s</b> (<code>%d</code>) была удалена его администраторами.",
				u.Message.Chat.Title, u.Message.Chat.ID,
			),
		)
		tgtMsg.ParseMode = "HTML"
		b.API.Send(tgtMsg)
		curMsg.Text = fmt.Sprintf("Связь с чатом <code>%d</code> была удалена @%s.", tgtChatID, u.Message.From.UserName)
		curMsg.ParseMode = "HTML"
		b.API.Send(curMsg)
		return
	}
	if incReq, err := b.RequestRepo.GetRequest(tgtChatID, u.Message.Chat.ID); err == nil {
		b.DeclineRequest(incReq.SrcChatID, incReq.TgtChatID, int(incReq.TgtMessageID), u.Message.From.UserName)
		return
	}
	if outReq, err := b.RequestRepo.GetRequest(u.Message.Chat.ID, tgtChatID); err == nil {
		b.RequestRepo.DeleteRequest(outReq.SrcChatID, outReq.TgtChatID)
		delTgtMsg := tgbotapi.NewDeleteMessage(outReq.TgtChatID, int(outReq.TgtMessageID))
		b.API.Request(delTgtMsg)
		tgtMsg := tgbotapi.NewMessage(
			outReq.TgtChatID,
			fmt.Sprintf("Запрос на связь от чата <b>%s</b> (<code>%d</code>) был отменен его администраторами.", u.Message.Chat.Title, outReq.SrcChatID),
		)
		tgtMsg.ParseMode = "HTML"
		b.API.Send(tgtMsg)
		curMsg.Text = fmt.Sprintf("Запрос на связь к чату <code>%d</code> был отменен @%s.", outReq.TgtChatID, u.Message.From.UserName)
		curMsg.ParseMode = "HTML"
		b.API.Send(curMsg)
		return
	}
	curMsg.Text = "Указанный чат не связан с вашим, запросов на связь не найдено."
	b.API.Send(curMsg)
}

func (b *Bot) DeclineRequestBtn(u *tgbotapi.Update) {
	srcChatID, _ := strconv.ParseInt(strings.TrimPrefix(u.CallbackData(), "d"), 10, 64)
	b.DeclineRequest(srcChatID, u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID, u.CallbackQuery.From.UserName)
}
