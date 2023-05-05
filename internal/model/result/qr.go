package result

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type QrEnterCode struct {
	Count     int
	OnSuccess string
}

func (q *QrEnterCode) Render() *tgbotapi.MessageConfig {
	msg := &tgbotapi.MessageConfig{
		Text: "Привет, введи цифровой код",
	}

	return msg
}

type QrSuccess struct {
	Msg   string
	Owner string
	Title string
}

func (q *QrSuccess) Render() *tgbotapi.MessageConfig {
	msg := &tgbotapi.MessageConfig{
		Text: fmt.Sprintf("Код %s был использован", q.Title),
	}

	return msg
}

type QrError struct {
	Success bool
	Msg     string
	User    string
}

func (q *QrError) Render() *tgbotapi.MessageConfig {
	return &tgbotapi.MessageConfig{
		Text: fmt.Sprintf(q.Msg, q.User),
	}
}
