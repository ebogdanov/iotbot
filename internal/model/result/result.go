package result

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Message interface {
	Render() *tgbotapi.MessageConfig
}

type Base struct {
	Success bool
}

type Execute struct {
	Base
}

type TextList struct {
	Base
	List map[string]string
}

type GroupList struct {
	Base
	Groups []string
}
