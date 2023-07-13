package result

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UnknownCmd struct {
	Success bool
	UserID  string
}

func (u *UnknownCmd) Render(chatID int64) tgbotapi.Chattable {
	menu := &MainMenu{
		UserID: u.UserID,
	}

	msg := tgbotapi.MessageConfig{
		Text:     "Не понял команды, вот что я умею делать",
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}

	msg.ReplyMarkup = menu.Render(chatID)

	return &msg
}
