package handler

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"su27bot/internal/acl"
	"su27bot/internal/config"
	"su27bot/internal/model/result"
	"time"
)

const (
	qrStartPrefix   = "/start QR:"
	qrPrefix        = "/qr"
	qrCommandAdd    = "/qr add" // N means how many times it can be used, 0 - for infinite
	qrCommandDelete = "/qr delete"
	qrCommandList   = "/qr list"
	qrHelp          = "/qr help"
)

var (
	errQRNotValidFormat = errors.New("неверный формат команды")
)

var (
	regExp6Digits = regexp.MustCompilePOSIX("[0-9]{6}")
)

type Qr struct {
	users *acl.Default
	cfg   *config.Config
	cmd   string
	h     []Handler
}

func NewQr(u *acl.Default, cfg *config.Config, h []Handler) *Qr {
	return &Qr{
		users: u,
		cfg:   cfg,
		cmd:   cfg.Qr.Cmd,
		h:     h,
	}
}

func (q *Qr) Supported(cmd string) bool {
	return strings.HasPrefix(cmd, qrStartPrefix) ||
		strings.HasPrefix(cmd, qrPrefix) ||
		regExp6Digits.MatchString(cmd)
}

func (q *Qr) Allowed(cmd, userID string) (bool, error) {
	if strings.HasPrefix(cmd, qrStartPrefix) {
		return true, nil
	}

	if strings.HasPrefix(cmd, qrPrefix) && q.users.IsMember("all", userID) {
		return true, nil
	}

	// @todo: 6 Digits allowed for not members only
	if regExp6Digits.MatchString(cmd) /*&& !q.users.IsMember("all", userID)*/ {
		return true, nil
	}

	// Check somehow that method is allowed for user
	return false, nil
}

func (q *Qr) Menu(_ string) []result.MenuItem {
	return nil
}

func (q *Qr) Execute(ctx context.Context, cmd, userID, user string) result.Message {
	switch {

	case strings.HasPrefix(cmd, qrStartPrefix):
		return &result.QrEnterCode{}

	// Check code, if matches: execute requested method + send notification to user, who created it
	case regExp6Digits.MatchString(cmd):
		// @todo Flood control
		if owner, title, _ := q.check(cmd); owner != "" {
			for i := 0; i < len(q.h); i++ {
				if !q.h[i].Supported(q.cmd) {
					continue
				}

				ok, _ := q.h[i].Allowed(q.cmd, owner)
				if !ok {
					continue
				}

				// todo Here we can log actions

				_ = q.h[i].Execute(ctx, q.cmd, userID, user)

				if true { // Check if success here
					return &result.QrSuccess{Owner: owner, Title: title}
				} else {
					return &result.QrError{Msg: ""}
				}
			}
		}

		return &result.Joke{}

	// A bit of help
	case cmd == qrPrefix || cmd == qrHelp:
		msg := "Управление кодами для входа курьеров или гостей:\n\n"

		msg += fmt.Sprintf("%s _название_ – Добавление нового кода\n", qrCommandAdd)
		msg += fmt.Sprintf("%s _название_ _число_ – Добавление нового кода с ограничением использований\n", qrCommandAdd)
		msg += fmt.Sprintf("%s _название_ – Удаление существующего кода\n", qrCommandDelete)
		msg += fmt.Sprintf("%s – Вывести список существующих кодов\n\n", qrCommandList)

		msg += "Для того чтобы воспользоваться кодом входа, – нужно отсканировать QR над панелью вызова домофона и ввести " +
			"его после запроса\n"

		return &result.Success{Msg: msg}

	// Request code
	case strings.HasPrefix(cmd, qrCommandAdd):
		res, err := q.addCode(cmd, userID)
		if err == nil {
			return &result.Success{Msg: fmt.Sprintf("Добавлен цифровой код для входа: *%s*", res)}
		}

		return &result.Fail{Msg: err.Error()}

	// List codes
	case strings.HasPrefix(cmd, qrCommandList):
		msg, err := q.listCodes(userID)
		if err == nil {
			return &result.Success{Msg: msg}
		}

		return &result.Fail{Msg: "Ошибка при получении списка кодов"}

	// Delete code
	case strings.HasPrefix(cmd, qrCommandDelete):
		res, _ := q.deleteCode(cmd, userID)

		if res != "" {
			return &result.Success{Msg: fmt.Sprintf("Код %s удален", res)}
		}

		return &result.Fail{Msg: "Ошибка удаления кода"}
	}

	return nil
}

func (q *Qr) addCode(cmd, userID string) (string, error) {
	parts := strings.Split(cmd, " ")

	if len(parts) < 3 {
		return "", errQRNotValidFormat
	}

	code := randomCode(6)

	item := &config.Code{
		User:  userID,
		Title: parts[2],
		Code:  code,
		Times: 0,
	}

	// @todo Lookup same code for this user

	if len(parts) > 3 {
		if a, err := strconv.Atoi(parts[3]); err == nil {
			item.Times = a
		}
	}

	// *todo Mutex
	q.cfg.Qr.Codes = append(q.cfg.Qr.Codes, *item)
	err := q.cfg.SaveFile()

	return code, err
}

func (q *Qr) listCodes(userID string) (string, error) {
	list := ""

	for _, item := range q.cfg.Qr.Codes {
		if item.User == userID {
			if item.Times > 0 {
				list += fmt.Sprintf("%s: *%s* - %d входов\n", item.Title, item.Code, item.Times)
			} else {
				list += fmt.Sprintf("%s: *%s*\n", item.Title, item.Code)
			}
		}
	}

	if list == "" {
		return "У вас пока нет кодов. Можно попробовать создать\n", nil
	}

	list = "Список кодов:\n\n" + list

	return list, nil
}

func (q *Qr) deleteCode(cmd, userID string) (string, error) {
	parts := strings.Split(cmd, " ")

	if len(parts) < 3 {
		return "", errQRNotValidFormat
	}

	titleOrCode := parts[2]

	for i, item := range q.cfg.Qr.Codes {
		if item.User == userID || userID == "any" {
			if item.Title == titleOrCode || item.Code == titleOrCode {
				// @todo mutex
				q.cfg.Qr.Codes = remove(q.cfg.Qr.Codes, i)
				err := q.cfg.SaveFile()

				return fmt.Sprintf("%s (%s)", item.Title, item.Code), err
			}
		}
	}

	return "", nil
}

func (q *Qr) check(code string) (string, string, error) {
	if !regExp6Digits.MatchString(code) {
		return "", "", errQRNotValidFormat
	}

	for _, item := range q.cfg.Qr.Codes {
		if item.Code == code {
			return item.User, item.Title, nil
		}
	}

	return "", "", nil
}

func randomCode(n int) string {
	chars := []rune("0123456789")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}

func remove(s []config.Code, i int) []config.Code {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
