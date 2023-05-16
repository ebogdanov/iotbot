package handler

import (
	"context"
	"github.com/rs/zerolog"
	"su27bot/internal/data"
	"su27bot/internal/model/result"
)

type Emitter interface {
	Handle(context.Context, string, string, string) result.Message
}

type Handler interface {
	Supported(cmd string) bool
	Allowed(cmd, userID string) (bool, error)
	Menu(userID string) []result.MenuItem
	Execute(ctx context.Context, cmd, userID, user string) result.Message
	Name() string
}

type Default struct {
	h []Handler
	s *data.Storage
	l zerolog.Logger
}

func NewEmitter(logger zerolog.Logger, s *data.Storage, h ...Handler) *Default {
	return &Default{
		l: logger.With().Str("component", "handler").Logger(),
		h: h,
		s: s,
	}
}

func (c *Default) Handle(ctx context.Context, cmd, userID, user string) result.Message {
	c.l.Debug().
		Str("cmd", cmd).
		Str("user", user).
		Str("userID", userID).
		Msgf("new action requested")

	var (
		res result.Message
	)

	menuItems := make([]result.MenuItem, 0)
	// Find which one of handlers supports this command, and get result from it
	for i := 0; i < len(c.h); i++ {
		handlerItems := c.h[i].Menu(userID)

		if len(handlerItems) > 0 {
			menuItems = append(menuItems, handlerItems...)
		}

		if res != nil {
			continue
		}

		if !c.h[i].Supported(cmd) {
			continue
		}

		ok, err := c.h[i].Allowed(cmd, userID)

		if err != nil {
			c.l.Debug().
				Str("cmd", cmd).
				Str("handler", c.h[i].Name()).
				Err(err).
				Msg("error check if command is allowed for handler")
		}

		if ok {
			// Flood control goes here
			if !c.s.Groups.IsAdmin(userID) && c.s.Actions.CheckFlood(cmd, userID) {
				break
			}

			res = c.h[i].Execute(ctx, cmd, userID, user)
			c.s.Actions.Add(userID, cmd, c.h[i].Name(), true)
		}
	}

	// Default - show menu, or funny picture
	if res == nil {
		res = &result.Fail{Msg: result.HelpText}

		c.s.Actions.Add(userID, cmd, "none", false)
		cnt := c.s.Actions.Count("*", userID)

		if cnt > 1 {
			res = &result.Sticker{}
		}

		if len(menuItems) > 0 {
			res = &result.MainMenu{Msg: "Не понял команды, вот что я умею", Actions: menuItems}
		}
	} else {
		// Return menu in case of success
		if v, ok := res.(*result.Success); ok {
			res = &result.MainMenu{Msg: v.Msg, Actions: menuItems}
		}
	}

	return res
}
