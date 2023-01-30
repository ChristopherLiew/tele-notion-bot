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

func teleCommandHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *viper.Viper, slogger *zap.SugaredLogger) {

	switch update.Message.Command() {
	case "start":
		startCommand(update, bot, slogger)
	case "search":
		searchCommand(update, bot, cfg, slogger)
	case "help":
		helpCommand(update, bot)
	case "end":
		endCommand(update, bot, slogger)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command 😮!")
		if _, err := bot.Send(msg); err != nil {
			slogger.Errorw(err.Error())
		}
	}
}

func startCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, slogger *zap.SugaredLogger) {

	// Start your connection
	slogger.Infof("Begin authorisation to [%s]'s Notion workspace", update.Message.From.UserName)

	startMsg := "Welcome to the TeleNotionBot! Click on the button below to get authenticate and connect with your Notion Workspace :)"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, startMsg)
	bot.Send(msg)
}

func searchCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *viper.Viper, slogger *zap.SugaredLogger) {

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
}

func helpCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	// Provide help and link to any tutorials
	helpMsg := `
	Hi there! I can ...

	/start 	- Connect with your Notion workspace 🔗
	/search - Search for Pages & Databases on it 🔎
	/end 	- End this conversation 🛑

	See the Menu for all possible commands
	`
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMsg)
	bot.Send(msg)
}

func endCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, slogger *zap.SugaredLogger) {

	userName := update.Message.From.UserName
	userId := update.Message.From.ID
	exitMsg := fmt.Sprintf(
		`Good bye 👋🏻 Till next time <a href="tg://user?id=%d">@%s</a>!`,
		userId,
		userName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, exitMsg)
	msg.ParseMode = "HTML"
	bot.Send(msg)

	slogger.Infof("Terminating bot for User: %s", userName)
	os.Exit(0)
}
