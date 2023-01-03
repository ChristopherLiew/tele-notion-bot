package main

// 1) Authenticate Pages + login on Notion (Public Integration)
// 2) Select Commands
// - Search and Display Pages & DBs
// - Query selected DBs etc
// - Update selected Pages + DBs etc
// - Add new Pages and DBs
// 3) Perform operations on them
// 4) Render them in Telegram

import (
	"tele-notion-bot/config"
	"tele-notion-bot/logging"
	"tele-notion-bot/notion"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var TeleNotionBotCommands []tgbotapi.BotCommand

func init() {
	// Bot commands
	TeleNotionBotCommands = []tgbotapi.BotCommand{
		{
			Command:     "start",
			Description: "Connect with or Update your Notion Account",
		},
		{
			Command:     "search",
			Description: "Search for Pages or Databases on Notion",
		},
		{
			Command:     "help",
			Description: "All the help you could possibly need!",
		},
	}
}

func main() {
	// Set up configs and logging
	cfg := config.GetConfig()

	logger := logging.GetLogger()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Instantiate Bot
	sugar.Infof("Starting %s", cfg.GetString("TELEGRAM.TEST_BOT_NAME"))
	bot, err := tgbotapi.NewBotAPI(cfg.GetString("TELEGRAM.TEST_BOT_TOKEN"))
	if err != nil {
		sugar.Fatalw(err.Error())
	}
	bot.Debug = true

	// Bot configs
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	// Bot Commands
	botCommands := tgbotapi.NewSetMyCommands(TeleNotionBotCommands...)
	resp, err := bot.Request(botCommands)
	if err != nil {
		sugar.Fatalw(err.Error())
	}
	sugar.Infof("Bot commands have been set with response: %b", resp.Ok)

	// Process updates
	for update := range updates {

		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		notion.HandleBotCommands(update, bot, sugar)

	}
}
