package tg

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"strconv"
	"strings"
	"su27bot/internal/handler"
	result2 "su27bot/internal/model/result"
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
		c.logger.Debug().Interface("update.message", update.Message).Msg("Incoming request")

		if time.Now().Sub(update.Message.Time()) > acceptWindow {
			c.logger.Debug().Time("time", update.Message.Time()).Time("now", time.Now()).Msg("Message is too old, skipping")
			return
		}

		cmd = update.Message.Text
		userName = c.getFromUserMessage(update.Message.From)
		userId = strconv.FormatInt(update.Message.From.ID, 10)

		msg.ChatID = update.Message.From.ID

		// Ok, this message is reply to something
		if update.Message.ReplyToMessage != nil {
			if update.Message.ReplyToMessage.From.IsBot &&
				update.Message.ReplyToMessage.From.UserName == c.BotName() {
				if update.Message.ReplyToMessage.Text == handler.DeleteUserMsg {
					// CMD - delete userName, userId + input text
					cmd = fmt.Sprintf("%s:%s", result2.UserDelete, update.Message.Text)
				}
			}
		}
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

	res := h.Handler(context.Background(), cmd, userId, userName)

	// If this is Joke - add ChatID (Not Good - need to refactor somehow)
	if joke, ok := res.(*result2.Joke); ok {
		joke.ChatID = msg.ChatID
	}

	// If this is QR

	sendMsg := res.Render()
	sendMsg.ChatID = msg.ChatID
	sendMsg.DisableNotification = true
	sendMsg.ParseMode = tgbotapi.ModeMarkdownV2

	_, _ = c.bot.Send(sendMsg)
}

// func (c *Client) processHandlerResult(cmd string, user string, userId string, msg tgbotapi.MessageConfig, res result2.Message) {

// switch result.(type) {
// User requested Invite
/*	case *result2.Invite:
	// @todo Введите идентификатор пользователя в формате: кв, фамилия
	inviteLink := fmt.Sprintf("https://t.me/%s?start=%s", c.bot.Self.UserName, result.(*result2.Invite).InviteId)

	msg.Text = fmt.Sprintf("Отправьте уникальную ссылку новому пользователю: %s", inviteLink)
	_, _ = c.bot.Send(msg)*/

// We've some problem with permissions
/*case *result2.Joke:
// Some fun for not authorized users
joke := c.content.GetSticker(msg.ChatID)
joke.DisableNotification = true
joke.AllowSendingWithoutReply = true

_, _ = c.bot.Send(joke)*/

// NewIot User access granted
/*case *result2.Welcome:
msg.Text = fmt.Sprintf(result.(*result2.Welcome).Msg, user)
// Lookup if this text is within scenario actions
msg.ReplyMarkup = c.content.MainMenu(userId)
_, _ = c.bot.Send(msg)*/

/*case *result2.Success:
if len(result.(*result2.Success).Msg) == 0 {
	result.(*result2.Success).Msg = "Не понял команды, вот что я умею делать"
}
msg.Text = result.(*result2.Success).Msg
msg.ReplyMarkup = c.content.MainMenu(userId)
_, _ = c.bot.Send(msg)*/

// Incorrect invite code
/*case *result2.InviteError:
msg.Text = fmt.Sprintf(result.(*result2.InviteError).Msg, user)
_, _ = c.bot.Send(msg)*/

/*case *result2.AdminMenu:
msg.Text = "Выберите действие"
msg.ReplyMarkup = c.content.AdminMenu(cmd)
_, _ = c.bot.Send(msg)*/

/*case *result2.TextList:
	msg.ReplyMarkup = c.content.AddDeleteButtons(cmd)

	for _, item := range result.(*result2.TextList).List {
		msg.Text += item + "\n"
	}

	_, _ = c.bot.Send(msg)

case *model.UserRequestName:
	msg.Text = result.(*model.UserRequestName).Msg
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "",
		Selective:             true,
	}
	_, _ = c.bot.Send(msg)

case *result2.GroupList:
	msg.Text = "Список"
	msg.ReplyMarkup = c.content.GroupMenu(result.(*result2.GroupList).Groups, handler.GroupList)
	_, _ = c.bot.Send(msg)
*/
/*case *result2.Fail:
msg.Text = result.(*result2.Fail).Msg
_, _ = c.bot.Send(msg)*/

/*	default:
	msg.Text = "Я тут, и вот чем я могу быть полезен"
	// Lookup if this text is within scenario actions
	msg.ReplyMarkup = c.content.MainMenu(userId)
	_, _ = c.bot.Send(msg)
}*/
// }

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
