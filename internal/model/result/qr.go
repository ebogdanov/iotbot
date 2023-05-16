package result

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	HelpText = "Ошибка ⛔️\n\nЧтобы пользоваться входом по телефону вам нужно получить доступ к боту. Пожалуйста, сделайте запрос в общедомовом чате, с вами свяжутся"
)

type QrEnterCode struct {
	Count     int
	OnSuccess string
}

func (q *QrEnterCode) Render(chatID int64) tgbotapi.Chattable {
	msg := &tgbotapi.MessageConfig{
		Text:     "Введите код состоящий из 4-6 цифр",
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}

	return msg
}

type QrSuccess struct {
	Msg   string
	Owner string
	Title string
}

func (q *QrSuccess) Render(chatID int64) tgbotapi.Chattable {
	msg := &tgbotapi.MessageConfig{
		Text:     fmt.Sprintf("Код %s был использован", q.Title),
		BaseChat: tgbotapi.BaseChat{ChatID: chatID, DisableNotification: true},
	}

	return msg
}

type QrError struct {
	Success bool
	Msg     string
	User    string
}

func (q *QrError) Render(chatID int64) tgbotapi.Chattable {
	return &tgbotapi.MessageConfig{
		Text:     fmt.Sprintf(q.Msg, q.User),
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}
}

type QrHelp struct {
}

func (q *QrHelp) Render(chatID int64) tgbotapi.Chattable {
	return &tgbotapi.MessageConfig{
		Text:     HelpText,
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}
}
