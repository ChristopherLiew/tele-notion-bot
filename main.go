// REFACTOR UPDATES cannot be nested since it conflict with the main instance
// Revert back to UPDATE TYPE

package main

import (
	"tele-notion-bot/config"
	"tele-notion-bot/logging"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

	// bot commands
	teleNotionBotCommands := []tgbotapi.BotCommand{
		{
			Command:     "start",
			Description: "Connect & Authenticate with Notion 🔗",
		},
		{
			Command:     "search",
			Description: "Search for Pages & Databases 🔎",
		},
		{
			Command:     "help",
			Description: "At your service 👋🏻",
		},
		{
			Command:     "end",
			Description: "Stop the bot 🛑",
		},
	}
	botCommands := tgbotapi.NewSetMyCommands(teleNotionBotCommands...)
	resp, err := bot.Request(botCommands)
	if err != nil {
		sugar.Fatalw(err.Error())
	}
	sugar.Infof("Bot commands have been set with response: %b", resp.Ok)

	// process updates
	for update := range updates {
		if update.InlineQuery != nil {
			teleSearchQueryHandler(update, bot, cfg, sugar)
		} else if update.CallbackQuery != nil {
			continue
		} else if update.Message.IsCommand() {
			teleCommandHandler(update, bot, cfg, sugar)
		} else {
			continue
		}
	}
}
