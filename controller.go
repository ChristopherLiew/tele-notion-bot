// Controller serves as the entrypoint into the telegram bot and
// comprises all of its commands

package main

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func botUpdateHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *viper.Viper, slogger *zap.SugaredLogger) {

	switch update.Message.Command() {
	case "start":
		startCommandHandler(update, bot, slogger)
	case "search":
		searchCommandHandler(update, bot, cfg, slogger)
	case "help":
		helpCommandHandler(update, bot)
	case "end":
		endCommandHandler(update, bot, slogger)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command 😮!")
		if _, err := bot.Send(msg); err != nil {
			slogger.Errorw(err.Error())
		}
	}
}

func startCommandHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, slogger *zap.SugaredLogger) {

	// Start your connection
	slogger.Infof("Begin authorisation to [%s]'s Notion workspace", update.Message.From.UserName)

	startMsg := "Welcome to the TeleNotionBot! Click on the button below to get authenticate and connect with your Notion Workspace :)"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, startMsg)
	bot.Send(msg)
}

func searchCommandHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *viper.Viper, slogger *zap.SugaredLogger) {

	getPagesQuery := "pages"
	getDBQuery := "databases"

	searchCommandKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:                         "Get Pages 📝",
				SwitchInlineQueryCurrentChat: &getPagesQuery,
			},
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:                         "Get Databases 💾",
				SwitchInlineQueryCurrentChat: &getDBQuery,
			},
		),
	)

	searchEntryMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "What would you like to search for on Notion?")
	searchEntryMsg.ReplyMarkup = searchCommandKeyboard
	if _, err := bot.Send(searchEntryMsg); err != nil {
		slogger.Error(err.Error())
	}

	// Get and validate response (Check if callback else exit)
	updateConfig := tgbotapi.NewUpdate(update.UpdateID + 1)
	updateConfig.Timeout = 30
	latestUpdates, err := bot.GetUpdates(updateConfig)
	if err != nil {
		slogger.Error(err.Error())
	}
	latestResp, hasResp := Last(latestUpdates)

	// Handle search related callbacks
	if hasResp {
		handleSearchQuery(latestResp, bot, cfg, slogger)
	} else {
		slogger.Fatalw("Unable to obtain latest response from user!")
	}
}

func helpCommandHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	// Provide help and link to any tutorials
	helpMsg := `
	Hi there! I only understand the following:

	/start 	- Connect with Notion 🔌
	/update - Update your Notion connection 🔐
	/search - Search for Pages & Databases 🔎
	/end 	- Stop the bot 🛑

	See the Menu for all possible commands
	`
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMsg)
	bot.Send(msg)
}

func endCommandHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, slogger *zap.SugaredLogger) {

	userName := update.Message.From.UserName
	exitMsg := fmt.Sprintf("Good bye! Till next time %s", userName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, exitMsg)
	bot.Send(msg)

	slogger.Infof("Terminating bot for User: %s", userName)
	os.Exit(0)
}
