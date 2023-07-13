package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"su27bot/internal/acl"
	"su27bot/internal/config"
	"su27bot/internal/data"
	"su27bot/internal/handler"
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
		log.Fatal().Err(err).Msg("Unable to read config file")
	}
	result.LoadStickers(cfg.Stickers)

	storage, err := data.NewStorage(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to init DB connection")
	}

	rules := acl.New(cfg, storage)
	handlers := handler.Factory(cfg, rules, storage, log.Logger)

	h := handler.NewEmitter(log.Logger, storage, handlers...)

	tgBot, err := tg2.NewClient(cfg.Telegram.TokenBot, log.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create Telegram connection")
	}
	storage.BotName = tgBot.BotName()

	log.Debug().Msg("Starting event handling...")
	tgBot.Start(h)
}
