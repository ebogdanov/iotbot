package result

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
)

const (
	onlyRegisteredUsers = "К сожалению, бот не сможет обработать ваш запрос."
	defaultSticker      = "CAACAgIAAxkBAAEGA5ZjP0t4akC6oPNFI5TwTZBzNK8wKAACrQMAAkcVaAlF6T5von_smSoE"
)

var stickers []string

func LoadStickers(items []string) {
	stickers = items
}

type Joke struct {
	Success bool
	ChatID  int64
	Text    string
}

func (j *Joke) Render() *tgbotapi.MessageConfig {
	msg := &tgbotapi.MessageConfig{Text: onlyRegisteredUsers}

	var (
		sticker tgbotapi.StickerConfig
	)

	if len(stickers) == 0 {
		sticker = tgbotapi.NewSticker(j.ChatID, tgbotapi.FileID(defaultSticker))
	} else {
		randomIndex := rand.Intn(len(stickers))

		sticker = tgbotapi.NewSticker(j.ChatID, tgbotapi.FileID(stickers[randomIndex]))
	}

	sticker.DisableNotification = true
	sticker.AllowSendingWithoutReply = true

	msg.ReplyMarkup = sticker

	return msg
}
