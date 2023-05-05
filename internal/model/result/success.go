package result

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Success struct {
	User    string
	Success bool
	Msg     string
}

func (s *Success) Render() *tgbotapi.MessageConfig {
	menu := &MainMenu{UserID: s.User}

	msg := &tgbotapi.MessageConfig{
		Text: s.Msg,
	}

	msg.ReplyMarkup = menu.Render()

	return msg
}
