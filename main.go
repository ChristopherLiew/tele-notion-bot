package main

// 1) Authenticate Pages + login on Notion (Public Integration)
// 2) Select Options by Page type DB, Pages, Users, Blocks, Comments
// 3) Perform operations on them
// 4) Render them in Telegram

import (
	"log"
	"tele-notion-bot/config"
	"tele-notion-bot/logging"
	"tele-notion-bot/notion"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Set up configs and logging
	cfg := config.GetConfig()

	logger := logging.GetLogger()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Instantiate Bot
	sugar.Infof("Starting %s", cfg.GetString("TELEGRAM.BOT_NAME"))
	bot, err := tgbotapi.NewBotAPI(cfg.GetString("TELEGRAM.BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {

		if update.Message != nil {
			notion.HandleNotionRequest(update, bot, sugar)
		}
	}
}
