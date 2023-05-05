package handler

import (
	"context"
	"su27bot/internal/model/result"
)

const (
	DeleteUserMsg = "Введите имя пользователя для удаления"
)

type Handler interface {
	Supported(cmd string) bool
	Allowed(cmd, userID string) (bool, error)
	Menu(userID string) []result.MenuItem
	Execute(ctx context.Context, cmd, userID, user string) result.Message
}
