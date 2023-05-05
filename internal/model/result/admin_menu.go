package result

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	MenuAdmin = "MENU_ADMIN"
	MenuMain  = "MENU_MAIN"

	MenuUsers   = "MENU_ADMIN_USERS"
	MenuGroups  = "MENU_ADMIN_GROUPS"
	MenuInvites = "MENU_ADMIN_INVITES"

	UserList   = "ADMIN_USER_LIST"
	UserDelete = "ADMIN_USER_DELETE"

	GroupsList       = "ADMIN_GROUP_LIST"
	GroupPermissions = "ADMIN_GROUP_PERMISSIONS"
)

type AdminMenu struct {
	Success bool
	Msg     string
	Section string
}

func (a *AdminMenu) Render() *tgbotapi.MessageConfig {
	msg := &tgbotapi.MessageConfig{
		Text: a.Msg,
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
			// Submenu - invites
			tgbotapi.NewInlineKeyboardButtonData("✉️ Приглашения", MenuInvites),
		), tgbotapi.NewInlineKeyboardRow(
			// Submenu - exit
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", MenuMain),
		))

	case MenuUsers:
		// Список пользователей
		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			// Добавить пользователя
			tgbotapi.NewInlineKeyboardButtonData("Пригласить", InviteGenerate),
		), tgbotapi.NewInlineKeyboardRow(
			// Список
			tgbotapi.NewInlineKeyboardButtonData("Список", UserList), // @todo
		), tgbotapi.NewInlineKeyboardRow(
			// Список последних действия
			tgbotapi.NewInlineKeyboardButtonData("Последние 10 записей", "LAST_ACTIONS"), // @todo
		), tgbotapi.NewInlineKeyboardRow(
			// Удалить пользователя
			tgbotapi.NewInlineKeyboardButtonData("Удалить", UserDelete), // @todo
		), tgbotapi.NewInlineKeyboardRow(
			// В главное меню
			tgbotapi.NewInlineKeyboardButtonData("Назад", MenuAdmin),
		))

	case MenuGroups:
		// Список групп
		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			// Список
			tgbotapi.NewInlineKeyboardButtonData("Список", GroupsList),
		), tgbotapi.NewInlineKeyboardRow(
			// Список последних действий
			tgbotapi.NewInlineKeyboardButtonData("Разрешения", GroupPermissions), // @todo
		), tgbotapi.NewInlineKeyboardRow(
			// В главное меню
			tgbotapi.NewInlineKeyboardButtonData("Назад", MenuAdmin),
		))

	case MenuInvites:
		// Список приглашений
		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			// Добавить пользователя
			tgbotapi.NewInlineKeyboardButtonData("Пригласить", InviteGenerate),
		), tgbotapi.NewInlineKeyboardRow(
			// Список
			tgbotapi.NewInlineKeyboardButtonData("Список", "INVITE_LIST"), // @todo
		), tgbotapi.NewInlineKeyboardRow(
			// Список последних действий
			tgbotapi.NewInlineKeyboardButtonData("Удалить все", "CLEAR_INVITES"), // @todo
		), tgbotapi.NewInlineKeyboardRow(
			// В главное меню
			tgbotapi.NewInlineKeyboardButtonData("Назад", MenuAdmin),
		))
	}

	msg.ReplyMarkup = menu

	return msg
}
