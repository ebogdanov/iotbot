package handler

import (
	"context"
	"strings"
	"su27bot/internal/acl"
	"su27bot/internal/iot"
	"su27bot/internal/model/result"
)

const (
	nfcPrefix = "/start NFC:"
)

type Iot struct {
	provider iot.Provider
	a        *acl.Default
}

func NewIot(p iot.Provider, a *acl.Default) *Iot {
	return &Iot{provider: p, a: a}
}

func (i *Iot) Supported(cmd string) bool {
	cmd = strings.Replace(cmd, nfcPrefix, "", 1) // Support of message from NFC tag

	items, err := i.provider.Actions(context.Background())
	if err == nil {
		for _, c := range items {
			if c.MenuId() == cmd {
				return true
			}

			if c.Name == cmd {
				return true
			}
		}
	}

	return false
}

func (i *Iot) Allowed(cmd, userID string) (bool, error) {
	cmd = strings.Replace(cmd, nfcPrefix, "", 1) // Support of message from NFC tag

	if i.a.IsAllowed(userID, cmd) {
		return true, nil
	}

	items, err := i.provider.Actions(context.Background())
	if err != nil {
		return false, nil
	}

	for _, c := range items {
		if c.MenuId() == cmd {
			return true, nil
		}

		if c.Name == cmd {
			return true, nil
		}
	}

	return false, nil
}

func (i *Iot) Execute(ctx context.Context, cmd, userId, user string) result.Message {
	// @todo: Count limit for this user per hour (20)

	// Check that this is allowed
	parts := strings.Split(cmd, "_")
	if len(parts) != 3 {
		// Try to find this id (maybe user sent us text command)
		items, err := i.provider.Actions(ctx)

		if err == nil {
			for _, c := range items {
				if c.Name == cmd {
					cmd = c.MenuId()
					parts = strings.Split(cmd, "_")

					break
				}
			}
		}

		if len(parts) != 3 {
			return &result.Fail{Msg: "Invalid request"}
		}
	}

	acts, err := i.provider.Actions(ctx)
	if err != nil {
		return &result.Fail{Msg: "Invalid request"}
	}

	// Lookup status of device
	for _, item := range acts {
		if item.ID != parts[2] {
			continue
		}

		res, err := i.provider.DeviceInfo(ctx, item.DeviceID)

		if err != nil {
			return &result.Fail{Msg: "Не получилось :( (ошибка связи с устройством)"}
		}

		if !res.Result.Online {
			// lastDateTime := time.Since(time.Unix(res.Result.ActiveTime, 0))

			return &result.Fail{Msg: "Не получилось из-за ошибки связи с устройством :("}
		}
	}

	res, err := i.provider.StartScenario(ctx, parts[1], parts[2])
	if err != nil {
		return &result.Fail{Msg: "Не получилось :( (ошибка запуска сценария)"}
	}

	return &result.Success{Success: res, Msg: "Успешно выполнено, жду приказов"}
}

func (i *Iot) Menu(userID string) []result.MenuItem {
	items := make([]result.MenuItem, 0)
	actions, err := i.provider.Actions(context.Background())

	if err != nil {
		return items
	}

	for c := 0; c < len(actions); c++ {
		if !i.a.IsAllowed(userID, actions[c].Name) {
			continue
		}

		items = append(items, result.MenuItem{
			ID:    actions[c].MenuId(),
			Title: actions[c].Name,
		})
	}

	return items
}
