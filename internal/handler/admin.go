package handler

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"su27bot/internal/data"
	"su27bot/internal/model/result"
	"time"
)

type Admin struct {
	s *data.Storage
}

const (
	welcomeSuccess = "Привет, %s.\n\nСпасибо за регистрацию! ❤️\nТеперь можно давать боту команды 👍"
	incorrectCode  = "Код приглашения неверный."
	systemFailure  = "Ошибка системы, попробуйте еще раз попозже"
)

const (
	adminPrefix       = "ADMIN_"
	startInvitePrefix = "/start I"
)

func NewAdmin(s *data.Storage) *Admin {
	return &Admin{s: s}
}

func (a *Admin) Supported(cmd string) bool {
	return strings.HasPrefix(cmd, result.MenuAdmin) ||
		strings.HasPrefix(cmd, adminPrefix) ||
		strings.HasPrefix(cmd, startInvitePrefix)
}

func (a *Admin) Allowed(cmd, userID string) (bool, error) {
	// Allow for registrations
	if strings.HasPrefix(cmd, startInvitePrefix) {
		return true, nil
	}

	if !a.s.Groups.IsAdmin(userID) {
		return false, fmt.Errorf("user %s is not allowed to execute admin command %s", userID, cmd)
	}

	return true, nil
}

func (a *Admin) RegisterByInvite(cmd, userId, user string) (result.Message, error) {
	parts := strings.Fields(cmd)
	if len(parts) != 2 {
		return &result.Fail{Msg: incorrectCode}, errors.New("неверный код приглашения")
	}
	code := parts[1]

	if !a.s.Invites.Check(code) {
		return &result.Fail{Msg: incorrectCode}, errors.New("код приглашения не найден")
	}

	if !a.s.Users.Check(userId) {
		res, err := a.s.Users.Add(userId, user)
		if err != nil || !res {
			return &result.Fail{Msg: systemFailure}, errors.New("ошибка добавления пользователя")
		}
	}

	// Delete from invites
	a.s.Invites.Delete(code)

	menuItems := []result.MenuItem{{ID: "Далее-" + time.Now().String(), Title: "Далее", Icon: `⏩`}}
	text := fmt.Sprintf(welcomeSuccess, user)

	return &result.MainMenu{Msg: text, Actions: menuItems, UserID: userId}, nil
}

func (a *Admin) Execute(_ context.Context, cmd, userId, user string) result.Message {
	if allow, _ := a.Allowed(cmd, userId); !allow {
		return &result.Sticker{}
	}

	switch {
	// Admin menu
	case strings.HasPrefix(cmd, result.MenuAdmin):
		return &result.AdminMenu{Success: true, Msg: "Выберите действие", Section: cmd}

	// Generate Invite
	case cmd == result.InviteGenerate:
		uid := "I" + randomString(10)

		res, err := a.s.Invites.Add(uid)
		if err != nil || !res {
			return &result.Fail{Msg: systemFailure}
		}

		return &result.Invite{InviteId: uid, BotUserName: a.s.BotName}

	// Log of Actions
	case cmd == result.ActionsLast:
		list, err := a.s.Actions.List(20)
		if err != nil {
			return &result.Fail{Msg: err.Error()}
		}

		return &result.ActionList{
			Msg:      "Список действий",
			List:     list,
			Previous: result.MenuUsers,
		}

	case cmd == result.ActionsUnknown:
		list, err := a.s.Actions.ListUnknown(20)
		if err != nil {
			return &result.Fail{Msg: err.Error()}
		}

		return &result.ActionList{
			Msg:      "Список действий",
			List:     list,
			Previous: result.MenuUsers,
		}

	// Request users list
	case cmd == result.UserList:
		users := a.s.Users.Active()

		return &result.AdminList{
			List:         users,
			Msg:          "Список пользователей",
			ActionView:   result.UserView,
			ActionDelete: result.UserDelete,
			Previous:     result.MenuUsers,
		}

	// View user
	case strings.HasPrefix(cmd, result.UserView):
		parts := strings.Split(cmd, "_")

		if len(parts) > 3 {
			userID := parts[3]

			// Get user info
			_, userName, active, err := a.s.Users.Info(userID)
			if err == nil {
				// Get last actions for user
				actions, _ := a.s.Actions.ListUser(userID, 10)
				// Get user groups
				memberOf := a.s.Groups.MemberOf(userID)
				// Get other groups
				groups := a.s.Groups.List()

				return &result.AdminViewUser{
					UserName: userName,
					UserID:   userID,
					Previous: result.MenuUsers,
					Active:   active,
					Actions:  actions,
					MemberOf: memberOf,
					Groups:   groups,
				}
			}
		}

		return &result.Fail{Msg: "Не удалось загрузить информацию о пользователе"}

	// Delete user
	case strings.HasPrefix(cmd, result.UserDelete):
		parts := strings.Split(cmd, "_")

		if len(parts) > 3 {
			userID := parts[3]

			res := a.s.Users.Delete(userID)
			if res {
				a.s.Groups.DeleteMember(userID, "*")

				return &result.AdminMenu{Msg: "Пользователь был удален", Section: result.MenuUsers}
			}
		}

	// Groups list
	case cmd == result.GroupsList:
		members := a.s.Groups.List()

		return &result.AdminList{
			List:         members,
			Msg:          "Список групп",
			ActionView:   result.GroupView,
			ActionDelete: result.GroupDelete,
			Previous:     result.MenuGroups,
		}

	// Delete group
	case strings.HasPrefix(cmd, result.GroupDelete):
		parts := strings.Split(cmd, "_")

		if len(parts) > 3 && parts[3] != "1" {
			groupID := parts[3]

			res := a.s.Groups.Delete(groupID)
			if res {
				return &result.AdminMenu{Msg: "Группа была удалена", Section: result.MenuGroups}
			}
		}

		return &result.Fail{Msg: "Не удалось удалить группу"}

	// View group
	case strings.HasPrefix(cmd, result.GroupView):
		parts := strings.Split(cmd, "_")

		if len(parts) > 3 {
			groupID := parts[3]
			var members map[string]string

			title := a.s.Groups.Title(groupID)

			if title != "" {
				users := a.s.Users.Active()
				list := a.s.Groups.List()

				for id, item := range list {
					if id == groupID {
						members = a.s.Groups.Members(item)
						break
					}
				}

				return &result.AdminViewGroup{
					Title:    title,
					GroupID:  groupID,
					Previous: result.MenuGroups,
					Members:  members,
					Users:    users,
				}
			}
		}
		return &result.Fail{Msg: "Не удалось загрузить информацию о группе"}

	// Add user to group
	case strings.HasPrefix(cmd, result.GroupMemberAdd):
		parts := strings.Split(cmd, "_")

		if len(parts) > 3 {
			userID := parts[3]
			groupID := parts[4]

			res := a.s.Groups.AddMember(groupID, userID)
			if res {
				return &result.AdminMenu{Msg: "Пользователь был добавлен в группу", Section: result.MenuGroups}
			}
		}

		return &result.Fail{Msg: "Не удалось добавить пользователя в группу"}

	// Delete user from group
	case strings.HasPrefix(cmd, result.GroupMemberDelete):
		parts := strings.Split(cmd, "_")

		if len(parts) > 3 {
			userID := parts[3]
			groupID := parts[4]

			if groupID != "1" {
				res := a.s.Groups.DeleteMember(userID, groupID)
				if res {
					return &result.AdminMenu{Msg: "Пользователь был удален из группы", Section: result.MenuGroups}
				}
			}
		}

		return &result.Fail{Msg: "Не удалось удалить пользователя из группы"}

	// Invite
	case strings.HasPrefix(cmd, startInvitePrefix):
		// Check if this /start is sent with command or not
		res, _ := a.RegisterByInvite(cmd, userId, user)

		return res
	}

	return nil
}

func (a *Admin) Menu(userID string) []result.MenuItem {
	if a.s.Groups.IsAdmin(userID) {
		res := []result.MenuItem{{ID: result.MenuAdmin, Title: "Управление", Icon: `🚓`}}

		return res
	}

	return nil
}

func (a *Admin) Name() string {
	return "admin"
}

func randomString(n int) string {
	chars := []rune("_123456789abcdefghijklmnpqrstuvwxyzABCDEFGHIJKLMNPRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
