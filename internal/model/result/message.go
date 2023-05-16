package result

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Message interface {
	Render(int64) tgbotapi.Chattable
}
