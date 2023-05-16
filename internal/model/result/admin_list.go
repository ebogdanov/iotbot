package result

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdminList struct {
	List         map[string]string
	Msg          string
	ActionView   string
	ActionDelete string
	Previous     string
	AllowDelete  bool
}

func (a *AdminList) Render(chatID int64) tgbotapi.Chattable {
	msg := &tgbotapi.MessageConfig{
		Text:     a.Msg,
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}

	menu := tgbotapi.NewInlineKeyboardMarkup()

	for key, item := range a.List {
		del := a.ActionDelete + "_" + key
		view := a.ActionView + "_" + key

		if a.AllowDelete {
			menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(item, view),
				tgbotapi.NewInlineKeyboardButtonData("❌", del),
			))
		} else {
			menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(item, view),
			))
		}
	}

	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		// В меню пользователей
		tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", a.Previous),
	))

	msg.ReplyMarkup = menu

	return msg
}
