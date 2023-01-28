package main

import (
	"tele-notion-bot/config"
	"tele-notion-bot/logging"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var TeleNotionBotCommands []tgbotapi.BotCommand

func init() {
	// bot commands
	TeleNotionBotCommands = []tgbotapi.BotCommand{
		{
			Command:     "start",
			Description: "Connect with or Update your Notion Account connection",
		},
		{
			Command:     "search",
			Description: "Search for Pages and Databases on Notion",
		},
		{
			Command:     "help",
			Description: "At your service",
		},
		{
			Command:     "end",
			Description: "Stop chatting with moi~",
		},
	}
}

func main() {

	// general configs
	cfg := config.GetConfig()
	logger := logging.GetLogger()
	defer logger.Sync()
	sugar := logger.Sugar()

	// instantiate bot
	sugar.Infof("Starting %s", cfg.GetString("TELEGRAM.TEST_BOT_NAME"))
	bot, err := tgbotapi.NewBotAPI(cfg.GetString("TELEGRAM.TEST_BOT_TOKEN"))
	if err != nil {
		sugar.Fatalw(err.Error())
	}
	bot.Debug = false

	// bot configs
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	botCommands := tgbotapi.NewSetMyCommands(TeleNotionBotCommands...)
	resp, err := bot.Request(botCommands)
	if err != nil {
		sugar.Fatalw(err.Error())
	}
	sugar.Infof("Bot commands have been set with response: %b", resp.Ok)

	// process updates
	for update := range updates {
		if update.Message.IsCommand() {
			BotUpdateHandler(update, bot, cfg, sugar)
		} else {
			continue
		}
	}
}
