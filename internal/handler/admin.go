package handler

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"su27bot/internal/acl"
	"su27bot/internal/model/result"
	"time"
)

type Admin struct {
	users *acl.Default
}

const (
	welcomeSuccess = "Привет, %s. Выбери действие:"
	incorrectCode  = "Код приглашения неверный."
)

const (
	adminPrefix       = "ADMIN_"
	startInvitePrefix = "/start I"
)

func NewAdmin(u *acl.Default) *Admin {
	return &Admin{users: u}
}

func (a *Admin) Supported(cmd string) bool {
	return strings.HasPrefix(cmd, result.MenuAdmin) ||
		strings.HasPrefix(cmd, adminPrefix) ||
		strings.HasPrefix(cmd, startInvitePrefix)
}

func (a *Admin) Allowed(cmd, userID string) (bool, error) {
	if !a.users.IsAdmin(userID) {
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

	if !a.users.CheckInvite(code) {
		return &result.Fail{Msg: incorrectCode}, errors.New("код приглашения не найден")
	}

	a.users.Add(userId, user)

	// Delete from invites
	a.users.RemoveInvite(code)

	return &result.MainMenu{Msg: welcomeSuccess}, nil
}

func (a *Admin) Execute(ctx context.Context, cmd, userId, user string) result.Message {
	if allow, _ := a.Allowed(cmd, userId); !allow {
		return &result.Joke{}
	}

	switch {
	// Меню администратора
	case strings.HasPrefix(cmd, result.MenuAdmin):
		return &result.AdminMenu{Success: true, Msg: "Выберите действие", Section: cmd}

	// Generate Invite
	case cmd == result.InviteGenerate:
		uid := "I" + randomString(10)
		a.users.AddInvite(uid)

		return &result.Invite{InviteId: uid}

		/*	case cmd == result.UserList:
				return &result.TextList{List: a.users.GetUserList()}

			case cmd == result.GroupsList:
				members := a.users.GetMembers()

				groups := make([]string, 0)
				members.Range(func(key, value any) bool {
					groups = append(groups, key.(string))

					return true
				})

				return &result.GroupList{Groups: groups}

			case strings.HasPrefix(cmd, result.GroupsList+":"):
				// If this is 1st phase - ask for UserName
				parts := strings.Split(cmd, ":")
				if len(parts) == 1 {
					return &model.UserRequestName{Msg: result.GroupsList}
				}
				requestedGroup := parts[1]

				members := a.users.GetGroupMembers(requestedGroup)

				return &result.TextList{List: members}

			case strings.HasPrefix(cmd, result.UserDelete):
				// If this is 1st phase - ask for UserName
				parts := strings.Split(cmd, ":")
				if len(parts) == 1 {
					return &model.UserRequestName{Msg: DeleteUserMsg}
				}
				requestedUser := parts[1]

				// If entered user is admin - do not delete it
				if a.users.IsAdmin(requestedUser) {
					// c.l.Error().Str("user", user).Str("action", cmd).Msgf("user is not allowed to delete admin")
					return &result.Joke{}
				}

				requestedUser = strings.Trim(requestedUser, " @.")
				a.users.Remove(requestedUser)

				return &result.Success{Base: result.Base{Success: true}}*/

		// Invite
	case strings.HasPrefix(cmd, startInvitePrefix):
		// Check if this /start is sent with command or not
		res, _ := a.RegisterByInvite(cmd, userId, user)

		return res
	}

	return nil
}

func (a *Admin) Menu(userID string) []result.MenuItem {
	res := []result.MenuItem{{ID: result.MenuAdmin, Title: "Admin section", Icon: `🚓`}}

	return res
}

func randomString(n int) string {
	chars := []rune("_123456789abcdefghijklmnpqrstuvwxyzABCDEFGHIJKLMNPRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
