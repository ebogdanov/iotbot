package result

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
)

const (
	defaultSticker = "CAACAgIAAxkBAAEGA5ZjP0t4akC6oPNFI5TwTZBzNK8wKAACrQMAAkcVaAlF6T5von_smSoE"
)

var (
	stickers = []string{defaultSticker}
)

func LoadStickers(items []string) {
	stickers = append(stickers, items...)
}

type Sticker struct {
	Success bool
	Text    string
}

func (j *Sticker) Render(chatID int64) tgbotapi.Chattable {
	randomIndex := rand.Intn(len(stickers))
	sticker := tgbotapi.NewSticker(chatID, tgbotapi.FileID(stickers[randomIndex]))

	sticker.DisableNotification = true
	sticker.AllowSendingWithoutReply = true

	sticker.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}

	return sticker
}
