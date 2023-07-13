package tg

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"strconv"
	"strings"
	"su27bot/internal/handler"
	"time"
)

const acceptWindow = 30 * time.Second

type Client struct {
	bot    *tgbotapi.BotAPI
	logger zerolog.Logger
}

func NewClient(token string, l zerolog.Logger) (*Client, error) {
	client := &Client{
		logger: l.With().Str("client", "telegram").Logger(),
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return client, err
	}

	bot.Debug = false
	client.bot = bot
	return client, err
}

// Start - processes cycle
func (c *Client) Start(h handler.Emitter) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		go c.processUpdate(h, update)
	}
}

func (c *Client) processUpdate(h handler.Emitter, update tgbotapi.Update) {
	cmd, userName, userId := "", "", ""

	msg := tgbotapi.NewMessage(0, "")

	if update.Message != nil {
		if time.Since(update.Message.Time()) > acceptWindow {
			c.logger.Debug().
				Time("time", update.Message.Time()).
				Time("now", time.Now()).
				Msg("Message is too old, skipping")

			return
		}

		c.logger.Debug().Interface("update.message", update.Message).Msg("Incoming request")

		cmd = update.Message.Text
		userName = c.getFromUserMessage(update.Message.From)
		userId = strconv.FormatInt(update.Message.From.ID, 10)

		msg.ChatID = update.Message.From.ID
	}

	if update.CallbackQuery != nil {
		c.logger.Debug().Interface("update.CallbackQuery", update.CallbackQuery).Msg("Incoming request")

		cmd = update.CallbackQuery.Data
		userName = c.getFromUserMessage(update.CallbackQuery.From)
		userId = strconv.FormatInt(update.CallbackQuery.From.ID, 10)

		msg.ChatID = update.CallbackQuery.Message.Chat.ID
		msg.Text = ""
	}

	typing := tgbotapi.NewChatAction(msg.ChatID, tgbotapi.ChatTyping)
	_, _ = c.bot.Send(typing)

	c.logger.Info().
		Str("userName", userName).
		Str("userId", userId).
		Str("cmd", cmd).
		Msg("User requested command")

	res := h.Handle(context.Background(), cmd, userId, userName)

	sendMsg := res.Render(msg.ChatID)

	_, err := c.bot.Send(sendMsg)

	if err != nil {
		c.logger.Info().
			Err(err).
			Str("userName", userName).
			Str("userId", userId).
			Str("cmd", cmd).
			Interface("result", res).
			Msg("User send reply error")
	}
}

func (c *Client) BotName() string {
	return c.bot.Self.UserName
}

func (c *Client) getFromUserMessage(from *tgbotapi.User) string {
	user := from.UserName

	if user == c.BotName() {
		user = ""
	}

	if user == "" {
		user = strings.Trim(from.FirstName+" "+from.LastName, " ")
	}

	if user == "" {
		user = strconv.FormatInt(from.ID, 10)
	}

	return user
}
