package result

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	MenuAdmin = "MENU_ADMIN"
	MenuMain  = "MENU_MAIN"

	MenuUsers = "MENU_ADMIN_USERS"

	MenuGroups = "MENU_ADMIN_GROUPS"

	UserList   = "ADMIN_USER_LIST"
	UserDelete = "ADMIN_USER_DELETE"
	UserView   = "ADMIN_USER_VIEW"

	ActionsLast    = "ADMIN_LAST_ACTIONS"
	ActionsUnknown = "ADMIN_LAST_UNKNOWN"

	GroupsList        = "ADMIN_GROUP_LIST"
	GroupDelete       = "ADMIN_GROUP_DELETE"
	GroupMemberDelete = "ADMIN_GROUP_MEMBER_DELETE"
	GroupMemberAdd    = "ADMIN_GROUP_MEMBER_ADD"
	GroupView         = "ADMIN_GROUP_VIEW"
)

type AdminMenu struct {
	Success bool
	Msg     string
	Section string
}

func (a *AdminMenu) Render(chatID int64) tgbotapi.Chattable {
	msg := &tgbotapi.MessageConfig{
		Text:     a.Msg,
		BaseChat: tgbotapi.BaseChat{ChatID: chatID},
	}

	menu := tgbotapi.NewInlineKeyboardMarkup()

	switch a.Section {
	case MenuAdmin:
		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			// Submenu - users
			tgbotapi.NewInlineKeyboardButtonData("👥 Пользователи", MenuUsers),
		), tgbotapi.NewInlineKeyboardRow(
			// Submenu - groups
			tgbotapi.NewInlineKeyboardButtonData("👤 Группы", MenuGroups),
		), tgbotapi.NewInlineKeyboardRow(
			// Submenu - exit
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", MenuMain),
		))

	case MenuUsers:
		msg.Text += " (пользователи)"
		// Список пользователей
		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			// Список
			tgbotapi.NewInlineKeyboardButtonData("Список", UserList),
		), tgbotapi.NewInlineKeyboardRow(
			// Добавить пользователя
			tgbotapi.NewInlineKeyboardButtonData("Пригласить", InviteGenerate),
		), tgbotapi.NewInlineKeyboardRow(
			// Список последних действия
			tgbotapi.NewInlineKeyboardButtonData("Последние 20 действий", ActionsLast),
		), tgbotapi.NewInlineKeyboardRow(
			// Список последних действия
			tgbotapi.NewInlineKeyboardButtonData("Действия без авторизации", ActionsUnknown),
		), tgbotapi.NewInlineKeyboardRow(
			// В главное меню
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", MenuAdmin),
		))

	case MenuGroups:
		msg.Text += " (группы)"

		// Список групп
		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			// Список
			tgbotapi.NewInlineKeyboardButtonData("Список", GroupsList),
		), tgbotapi.NewInlineKeyboardRow(
			// В главное меню
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", MenuAdmin),
		))
	}

	msg.ReplyMarkup = &menu

	return msg
}
