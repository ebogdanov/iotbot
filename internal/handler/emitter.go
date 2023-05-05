package handler

import (
	"context"
	"github.com/rs/zerolog"
	"su27bot/internal/acl"
	"su27bot/internal/model/result"
)

type Emitter interface {
	Handler(ctx context.Context, cmd, userId, user string) result.Message
}

type Default struct {
	h     []Handler
	l     zerolog.Logger
	users *acl.Default
}

func NewEmitter(logger zerolog.Logger, u *acl.Default, h ...Handler) *Default {
	return &Default{
		l:     logger.With().Str("component", "handler").Logger(),
		users: u,
		h:     h,
	}
}

func (c *Default) Handler(ctx context.Context, cmd, userID, user string) result.Message {
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

		if !c.h[i].Supported(cmd) {
			continue
		}

		ok, err := c.h[i].Allowed(cmd, userID)

		if err != nil {
			c.l.Debug().
				Str("cmd", cmd).
				Err(err).
				Msg("error check if cmd is allowed for handler")
		}

		if ok && res == nil {
			res = c.h[i].Execute(ctx, cmd, userID, user)

			// todo Here we can log actions
			continue
		}
	}

	// Default - show menu, or funny picture
	if res == nil {
		if len(menuItems) > 0 {
			res = &result.MainMenu{Msg: "Не понял команды, вот что я умею", Actions: menuItems}
		} else {
			res = &result.Joke{}
		}
	}

	return res
}
