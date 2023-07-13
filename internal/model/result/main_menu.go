package result

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MenuItem struct {
	ID    string
	Title string
	Icon  string
}

type MainMenu struct {
	Msg       string
	UserID    string
	GroupName string
	Actions   []MenuItem
}

func (m *MainMenu) Render(chatID int64) tgbotapi.Chattable {
	msg := &tgbotapi.MessageConfig{
		Text:     m.Msg,
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}

	menu := tgbotapi.NewInlineKeyboardMarkup()

	for i := 0; i < len(m.Actions); i++ {
		title := m.Actions[i].Title

		if m.Actions[i].Icon != "" {
			title = m.Actions[i].Icon + "   " + m.Actions[i].Title
		}

		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(title, m.Actions[i].ID),
		))
	}

	if len(menu.InlineKeyboard) > 0 {
		msg.ReplyMarkup = menu
	}

	return msg
}
