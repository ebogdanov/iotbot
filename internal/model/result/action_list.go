package result

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"su27bot/internal/data"
)

type ActionList struct {
	List     []data.Action
	Msg      string
	Previous string
}

func (a *ActionList) Render(chatID int64) tgbotapi.Chattable {
	text := a.Msg + ": \n\n"

	mask := "%s: %s - %s(%s), %v \n"

	for _, item := range a.List {
		validUtf8 := bytes.ToValidUTF8([]byte(item.Cmd), []byte{0xef, 0xbf, 0xbd})
		item.Cmd = string(validUtf8)

		if item.ID == 0 { // user is not registered
			text += fmt.Sprintf(mask, item.EventTime, "! "+item.User+" !", item.Cmd, item.HandlerName, item.Result)
		} else {
			text += fmt.Sprintf(mask, item.EventTime, "@"+item.User, item.Cmd, item.HandlerName, item.Result)
		}
	}

	msg := &tgbotapi.MessageConfig{
		Text:     text,
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}

	menu := tgbotapi.NewInlineKeyboardMarkup()
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		// В меню пользователей
		tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", a.Previous),
	))

	msg.ReplyMarkup = menu

	return msg
}
