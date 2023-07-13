package handler

import (
	"context"
	"strings"
	"su27bot/internal/acl"
	"su27bot/internal/iot"
	"su27bot/internal/model/result"
)

const (
	IotCmdPrefix = "DO_"
)

type Iot struct {
	provider iot.Provider
	a        *acl.Default
}

func NewIot(p iot.Provider, a *acl.Default) *Iot {
	return &Iot{provider: p, a: a}
}

func (i *Iot) Supported(cmd string) bool {
	items, err := i.provider.Actions(context.Background(), []string{})
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
	// Check that this is Direct text
	switch strings.HasPrefix(cmd, IotCmdPrefix) {
	case false:
		if i.a.IsAllowed(userID, cmd) {
			return true, nil
		}
	case true:
		// Ok, this is the button press event
		items, err := i.provider.Actions(context.Background(), []string{})
		if err != nil {
			return false, nil
		}

		for _, c := range items {
			if c.MenuId() == cmd {
				return i.a.IsAllowed(userID, c.Name), nil
			}

			if c.Name == cmd {
				return i.a.IsAllowed(userID, c.Name), nil
			}
		}
	}

	return false, nil
}

func (i *Iot) Execute(ctx context.Context, cmd, userId, user string) result.Message {
	parts := strings.Split(cmd, "_")
	if len(parts) < 3 {
		// Try to find this id (maybe user sent us text command)
		items, err := i.provider.Actions(ctx, []string{})

		if err == nil {
			for _, c := range items {
				if c.Name == cmd {
					cmd = c.MenuId()
					parts = strings.Split(cmd, "_")

					break
				}
			}
		}
	}

	if len(parts) < 3 {
		return &result.Fail{Msg: "Invalid request"}
	}

	acts, err := i.provider.Actions(ctx, []string{})
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

func (a *Iot) Name() string {
	return "iot"
}

func (i *Iot) Menu(userID string) []result.MenuItem {
	items := make([]result.MenuItem, 0)
	actions, err := i.provider.Actions(context.Background(), []string{})

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
