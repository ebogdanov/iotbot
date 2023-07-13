package handler

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"su27bot/internal/config"
	"su27bot/internal/data"
	"su27bot/internal/model/result"
)

const (
	qrStartPrefix = "/start QR_"
	qrCmdPrefix   = "/qr"
	qrCmdAdd      = "/qr add" // N means how many times it can be used, 0 - for infinite
	qrCmdDelete   = "/qr delete"
	qrCmdList     = "/qr list"
	qrHelp        = "/qr help"
)

var (
	errQRNotValidFormat = errors.New("неверный формат команды")
)

var (
	regExp46Digits = regexp.MustCompilePOSIX("[0-9]{4,6}")
)

type Qr struct {
	s   *data.Storage
	cfg *config.Config
	cmd string
	h   []Handler
}

func NewQr(u *data.Storage, cfg *config.Config, h []Handler) *Qr {
	cmd := ""
	if cfg.Qr != nil {
		cmd = cfg.Qr.Cmd
	}
	return &Qr{
		s:   u,
		cfg: cfg,
		cmd: cmd,
		h:   h,
	}
}

func (q *Qr) Supported(cmd string) bool {
	return strings.HasPrefix(cmd, qrStartPrefix) ||
		strings.HasPrefix(cmd, qrCmdPrefix) ||
		regExp46Digits.MatchString(cmd)
}

func (q *Qr) Allowed(cmd, userID string) (bool, error) {
	if strings.HasPrefix(cmd, qrStartPrefix) {
		return true, nil
	}

	if strings.HasPrefix(cmd, qrCmdPrefix) && q.s.Users.Check(userID) {
		return true, nil
	}

	if q.cfg.Qr != nil && q.cfg.Qr.AllowCodes {
		if regExp46Digits.MatchString(cmd) && !q.s.Users.Check(userID) {
			return true, nil
		}
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
		// If user is registered - just execute action from config
		if q.s.Users.Check(userID) {
			// If this is format with QR:Action, use it
			cmd1 := strings.Replace(cmd, qrStartPrefix, "", 1)
			if cmd1 == "" {
				cmd1 = q.cmd
			}

			return q.executeAction(ctx, cmd1, userID, user)
		}

		// If not - enter code
		if q.cfg.Qr != nil && !q.cfg.Qr.Enable {
			return &result.QrHelp{}
		}

		return &result.QrEnterCode{}

	// Check code, if matches: execute requested method + send notification to user, who created it
	case regExp46Digits.MatchString(cmd):
		code := cmd
		if q.s.Codes.Check(code) {
			info, err := q.s.Codes.Info(code)
			if err != nil {
				return &result.QrError{Msg: "системная ошибка"}
			}

			owner := q.s.Users.Name(userID)
			res := q.executeAction(ctx, q.cmd, info.UserID, user)

			if _, ok := res.(*result.Success); ok {
				return &result.QrSuccess{Owner: owner, Title: info.Title}
			}

			return &result.QrError{Msg: "Неверный код"}
		}

		return &result.Sticker{}

	// A bit of help
	case cmd == qrCmdPrefix || cmd == qrHelp:
		msg := "Управление кодами для входа курьеров или гостей:\n\n"

		msg += fmt.Sprintf("%s _название_ – Добавление нового кода\n", qrCmdAdd)
		msg += fmt.Sprintf("%s _название_ _число_ – Добавление нового кода с ограничением использований\n", qrCmdAdd)
		msg += fmt.Sprintf("%s _название_ – Удаление существующего кода\n", qrCmdDelete)
		msg += fmt.Sprintf("%s – Вывести список существующих кодов\n\n", qrCmdList)

		msg += "Для того чтобы воспользоваться кодом входа, – нужно отсканировать QR над панелью вызова домофона и ввести " +
			"его после запроса\n"

		return &result.Success{Msg: msg}

	// Add new code
	case strings.HasPrefix(cmd, qrCmdAdd):
		res, err := q.addCode(cmd, userID)
		if err == nil {
			return &result.Success{Msg: fmt.Sprintf("Добавлен цифровой код для входа: *%s*", res)}
		}

		return &result.Fail{Msg: err.Error()}

	// List codes
	case strings.HasPrefix(cmd, qrCmdList):
		msg, err := q.listCodes(userID)
		if err == nil {
			return &result.Success{Msg: msg}
		}

		return &result.Fail{Msg: "Ошибка при получении списка кодов"}

	// Delete code
	case strings.HasPrefix(cmd, qrCmdDelete):
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

	var (
		res *data.Code
		err error
	)

	maxAttempts := 0
	if len(parts) >= 4 {
		if a, err := strconv.Atoi(parts[3]); err == nil {
			maxAttempts = a
		}
	}

	for i := 0; i < 3; i++ {
		code := randomCode(6)

		res, err = q.s.Codes.Info(code)
		if err != nil || res != nil {
			continue
		}

		item := &data.Code{
			UserID:      userID,
			Title:       parts[2],
			Code:        code,
			MaxAttempts: maxAttempts,
		}

		_, err = q.s.Codes.Add(item)
		if err != nil {
			return "", err
		}

		return code, err
	}

	return "", err
}

func (q *Qr) listCodes(userID string) (string, error) {
	list := ""

	codes, err := q.s.Codes.List(userID)
	if err != nil {
		return "", err
	}

	for _, item := range codes {
		if item.MaxAttempts > 0 {
			list += fmt.Sprintf("%s: *%s* - %d входов\n", item.Title, item.Code, item.MaxAttempts)
		} else {
			list += fmt.Sprintf("%s: *%s*\n", item.Title, item.Code)
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

	if q.s.Codes.Delete(titleOrCode, userID) {
		return fmt.Sprintf("Код (%s)", titleOrCode), nil
	}

	return "", nil
}

func (q *Qr) executeAction(ctx context.Context, cmd, userID, user string) result.Message {
	for i := 0; i < len(q.h); i++ {
		if q.h[i].Name() == q.Name() {
			// If we got loop for some strange issue - skip it
			continue
		}

		if !q.h[i].Supported(q.cmd) {
			continue
		}

		if ok, _ := q.h[i].Allowed(cmd, userID); !ok {
			continue
		}

		res := q.h[i].Execute(ctx, q.cmd, userID, user)

		return res
	}

	return &result.Fail{Msg: result.HelpText}
}

func (q *Qr) Name() string {
	return "qr"
}

func randomCode(n int) string {
	chars := []rune("0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
