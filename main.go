package main

import (
	"context"
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tuya/tuya-connector-go/connector/httplib"
	"os"
	"su27bot/internal/acl"
	"su27bot/internal/config"
	"su27bot/internal/handler"
	"su27bot/internal/iot"
	"su27bot/internal/model/result"
	tg2 "su27bot/internal/tg"
)

func main() {
	cfgPath := flag.String("config", "conf/config.yaml", "config path")
	flag.Parse()

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	log.Output(logger)

	// Read YAML config file
	cfg, err := config.New(*cfgPath)
	if err != nil {
		log.Fatal().Err(err)
	}

	result.LoadStickers(cfg.Stickers)

	storage := acl.New(cfg)

	// Init handlers
	handlers := make([]handler.Handler, 0)
	if cfg.Tuya != nil {
		tuyaClient := iot.NewTuyaClient(httplib.URL_EU, httplib.MSG_EU, cfg.Tuya)

		if _, err := tuyaClient.Actions(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Failed load actions from Tuya Cloud")
		}

		handlers = append(handlers, handler.NewIot(tuyaClient, storage))
	}

	if cfg.Ewelink != nil {
		ewelinkClient := iot.NewEwelinkWebsocket(cfg.Ewelink.Region, cfg.Ewelink.Username, cfg.Ewelink.Password)

		if _, err := ewelinkClient.Actions(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Failed load actions from Ewelink Cloud")
		}

		handlers = append(handlers, handler.NewIot(ewelinkClient, storage))
	}

	if cfg.Qr.Enable {
		qrHandler := handler.NewQr(storage, cfg, handlers)
		handlers = append(handlers, handler.NewAdmin(storage), qrHandler)
	}
	// End of init handlers

	h := handler.NewEmitter(log.Logger, storage, handlers...)

	tgBot, err := tg2.NewClient(cfg.Telegram.TokenBot, log.Logger)
	if err != nil {
		log.Panic().Err(err)
	}

	log.Debug().Msg("Starting event handling...")
	tgBot.Start(h)
}
