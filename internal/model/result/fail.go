package result

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Fail struct {
	Success bool
	Msg     string
}

func (f *Fail) Render() *tgbotapi.MessageConfig {
	return &tgbotapi.MessageConfig{
		Text: f.Msg,
	}
}
