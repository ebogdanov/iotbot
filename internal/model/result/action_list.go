package result

import (
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
		if item.UserID != "" {
			text += fmt.Sprintf(mask, item.EventTime, "@"+item.User, item.Cmd, item.HandlerName, item.Result)
		} else {
			text += fmt.Sprintf(mask, item.EventTime, item.UserID, item.Cmd, item.HandlerName, item.Result)
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
