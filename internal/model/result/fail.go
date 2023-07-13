package result

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Fail struct {
	Success bool
	Msg     string
}

func (f *Fail) Render(chatID int64) tgbotapi.Chattable {
	msg := &tgbotapi.MessageConfig{
		Text:     f.Msg,
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}

	return msg
}
