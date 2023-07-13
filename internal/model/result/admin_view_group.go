package result

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdminViewGroup struct {
	Title    string
	GroupID  string
	Previous string

	Members map[string]string
	Users   map[string]string
}

func (a *AdminViewGroup) Render(chatID int64) tgbotapi.Chattable {
	text := fmt.Sprintf("Группа: %s \n\n", a.Title)

	menu := tgbotapi.NewInlineKeyboardMarkup()

	for id, item := range a.Users {
		if _, ok := a.Members[id]; ok {
			menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(item+" ❌", GroupMemberDelete+"_"+id+"_"+a.GroupID),
			))
		} else {
			menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(item+" ➕", GroupMemberAdd+"_"+id+"_"+a.GroupID),
			))
		}
	}

	msg := &tgbotapi.MessageConfig{
		Text:     text,
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}

	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		// Удалить группу целиком
		tgbotapi.NewInlineKeyboardButtonData("❌ Удалить группу", GroupDelete+"_"+a.GroupID),
	), tgbotapi.NewInlineKeyboardRow(
		// В меню пользователей
		tgbotapi.NewInlineKeyboardButtonData("◀️ Назад в меню", a.Previous),
	))

	msg.ReplyMarkup = menu

	return msg
}
