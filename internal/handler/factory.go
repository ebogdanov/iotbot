package handler

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/tuya/tuya-connector-go/connector/httplib"
	"su27bot/internal/acl"
	"su27bot/internal/config"
	"su27bot/internal/data"
	"su27bot/internal/iot"
)

func Factory(cfg *config.Config, rules *acl.Default, storage *data.Storage, log zerolog.Logger) []Handler {
	supported := rules.SupportActions()
	handlers := make([]Handler, 0, 2)

	if cfg.Tuya != nil {
		tuyaClient := iot.NewTuyaClient(httplib.URL_EU, httplib.MSG_EU, cfg.Tuya)

		if _, err := tuyaClient.Actions(context.Background(), supported); err != nil {
			log.Fatal().Err(err).Msg("Failed load actions from Tuya Cloud")
		}

		handlers = append(handlers, NewIot(tuyaClient, rules))
	}

	if cfg.Ewelink != nil {
		ewelinkClient := iot.NewEwelinkWebsocket(cfg.Ewelink.Region, cfg.Ewelink.Username, cfg.Ewelink.Password)

		if _, err := ewelinkClient.Actions(context.Background(), supported); err != nil {
			log.Fatal().Err(err).Msg("Failed load actions from Ewelink Cloud")
		}

		handlers = append(handlers, NewIot(ewelinkClient, rules))
	}

	qrHandler := NewQr(storage, cfg, handlers)
	handlers = append(handlers, qrHandler)

	handlers = append(handlers, NewAdmin(storage))

	return handlers
}
