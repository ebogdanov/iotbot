package result

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	InviteGenerate = "ADMIN_INVITE"
)

type Invite struct {
	InviteId    string
	BotUserName string
}

func (i *Invite) Render(chatID int64) tgbotapi.Chattable {
	inviteLink := fmt.Sprintf("https://t.me/%s?start=%s", i.BotUserName, i.InviteId)

	msg := &tgbotapi.MessageConfig{
		Text: fmt.Sprintf(
			"Отправьте уникальную одноразовую ссылку новому пользователю: %s", inviteLink),
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}

	return msg
}

type InviteError struct {
	Msg  string
	User string
}

func (i *InviteError) Render(chatID int64) tgbotapi.Chattable {
	return &tgbotapi.MessageConfig{
		Text:     fmt.Sprintf(i.Msg, i.User),
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}
}

type InviteSuccess struct {
	Msg  string
	User string
}

func (i *InviteSuccess) Render(chatID int64) tgbotapi.Chattable {
	return &tgbotapi.MessageConfig{
		Text:     fmt.Sprintf(i.Msg, i.User),
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}
}
