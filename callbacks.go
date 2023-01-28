package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"tele-notion-bot/notion"
)

// Search
func handleSearchCallback(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *viper.Viper, slogger *zap.SugaredLogger) {

	var err error
	notionIntToken := cfg.GetString("NOTION.INTEGRATION_SECRET")

	switch update.CallbackQuery.Data {
	case "get-pages":
		if err = getPagesCallback(update, bot, notionIntToken, slogger); err != nil {
			slogger.Errorw("callback error: [unable to get pages via notion api]")
		}
	case "get-databases":
		slogger.Info("Getting databases from user's Notion")
	case "get-specific-page":
		slogger.Info("Scurrying to look for user's specific page")
	default:
		err = fmt.Errorf("unknown callback: [%s]", update.CallbackQuery.Data)
		slogger.Errorw(err.Error())
	}

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Something went wrong with your Search request 🥲")
		bot.Send(msg)
	}
}

func getPagesCallback(update tgbotapi.Update, bot *tgbotapi.BotAPI, notionIntToken string, slogger *zap.SugaredLogger) error {

	response, err := notion.GetAllPages(notionIntToken, slogger, "", 10)
	if err != nil {
		slogger.Errorw(err.Error())
		return err
	}

	// Process response to get Title, URL + Icon
	var allResults []tgbotapi.InlineKeyboardButton

	// Use inline requests to handle this
	for _, page := range response.Results {
		button := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("test", page.URL),
		)
		allResults = append(
			allResults,
			button...,
		)
	}

	// Render
	resultsKeyboard := tgbotapi.NewInlineKeyboardMarkup(allResults)
	resultsMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Here are your 10 most recent Notion pages!")
	resultsMsg.ReplyMarkup = resultsKeyboard
	if _, err := bot.Send(resultsMsg); err != nil {
		slogger.Error(err.Error())
	}

	// Handle user selection
	return nil
}
