package result

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Success struct {
	User    string
	Success bool
	Msg     string
}

func (s *Success) Render(chatID int64) tgbotapi.Chattable {
	menu := &MainMenu{UserID: s.User}

	msg := &tgbotapi.MessageConfig{
		Text:     s.Msg,
		BaseChat: tgbotapi.BaseChat{ChatID: chatID, DisableNotification: true},
	}

	msg.ReplyMarkup = menu.Render(chatID)

	return msg
}
