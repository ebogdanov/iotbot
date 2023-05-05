package result

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UnknownCmd struct {
	Success bool
	UserID  string
}

func (u *UnknownCmd) Render() *tgbotapi.MessageConfig {
	menu := &MainMenu{
		UserID: u.UserID,
	}

	msg := tgbotapi.MessageConfig{
		Text: "Не понял команды, вот что я умею делать",
	}

	msg.ReplyMarkup = menu.Render()

	return &msg
}
