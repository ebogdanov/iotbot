package result

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	InviteGenerate = "ADMIN_INVITE"
)

type Invite struct {
	Success     bool
	InviteId    string
	BotUserName string
}

func (i *Invite) Render() *tgbotapi.MessageConfig {
	inviteLink := fmt.Sprintf("https://t.me/%s?start=%s", i.BotUserName, i.InviteId)

	msg := &tgbotapi.MessageConfig{
		Text: fmt.Sprintf(
			"Отправьте уникальную одноразовую ссылку новому пользователю: %s", inviteLink),
	}

	return msg
}

type InviteError struct {
	Success bool
	Msg     string
	User    string
}

func (i *InviteError) Render() *tgbotapi.MessageConfig {
	return &tgbotapi.MessageConfig{
		Text: fmt.Sprintf(i.Msg, i.User),
	}
}
