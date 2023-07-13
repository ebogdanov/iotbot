package result

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"su27bot/internal/data"
)

type AdminViewUser struct {
	UserName string
	UserID   string
	Previous string

	Active   bool
	Actions  []data.Action
	MemberOf []string
	Groups   map[string]string
}

func (a *AdminViewUser) Render(chatID int64) tgbotapi.Chattable {
	const (
		mask = "%s: %s(%s), %v \n"
	)

	text := fmt.Sprintf("Пользователь: %s (%v) \n\n", a.UserName, a.Active)

	for _, item := range a.Actions {
		text += fmt.Sprintf(mask, item.EventTime, item.Cmd, item.HandlerName, item.Result)
	}
	text += "\n\nГруппы:\n"

	menu := tgbotapi.NewInlineKeyboardMarkup()

	for id, item := range a.Groups {
		if contains(a.MemberOf, item) {
			if item == "admin" {
				menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(item+" 🔎", GroupView+"_"+id),
				))
			} else {
				menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(item+" ❌", GroupMemberDelete+"_"+id+"_"+a.UserID),
				))
			}
		} else {
			menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(item+" ➕", GroupMemberAdd+"_"+id+"_"+a.UserID),
			))
		}
	}

	msg := &tgbotapi.MessageConfig{
		Text:     text,
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}

	menu.InlineKeyboard = append(menu.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			// Удалить пользователя
			tgbotapi.NewInlineKeyboardButtonData("❌ Удалить пользователя", UserDelete+"_"+a.UserID),
		),
		tgbotapi.NewInlineKeyboardRow(
			// В меню пользователей
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад в меню", a.Previous),
		))

	msg.ReplyMarkup = menu

	return msg
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
