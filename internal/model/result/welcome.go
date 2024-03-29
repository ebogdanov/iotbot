package result

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Welcome struct {
	Success bool
	Msg     string
	UserID  int64
	User    string
}

func (w *Welcome) Render(chatID int64) tgbotapi.Chattable {
	menu := &MainMenu{UserID: w.User}
	msg := &tgbotapi.MessageConfig{
		Text:     fmt.Sprintf(w.Msg, w.User),
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}
	msg.ReplyMarkup = menu.Render(chatID)

	return msg
}
