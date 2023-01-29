package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"tele-notion-bot/notion"
)

// Search
func handleSearchQuery(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *viper.Viper, slogger *zap.SugaredLogger) (err error) {

	notionSecret := cfg.GetString("NOTION.INTEGRATION_SECRET")

	switch update.InlineQuery.Query {
	case "pages":
		if err = queryPages(update, bot, notionSecret, slogger); err != nil {
			slogger.Errorw("callback error: [unable to get pages via notion api]")
		}
	case "databases":
		slogger.Info("Getting databases from user's Notion")
	default:
		err = fmt.Errorf("unknown callback: [%s]", update.CallbackQuery.Data)
		slogger.Errorw(err.Error())
	}

	return
}

func queryPages(update tgbotapi.Update, bot *tgbotapi.BotAPI, notionIntToken string, slogger *zap.SugaredLogger) (err error) {

	response, err := notion.GetAllPages(notionIntToken, slogger, nil, 10)
	if err != nil {
		slogger.Errorw(err.Error())
		return err
	}
	snippets := notion.GetPageSnippets(response.Results, slogger)

	// construct inline query results
	results := make([]tgbotapi.InlineQueryResultArticle, 0)
	for i, snippet := range snippets {
		results = append(
			results,
			tgbotapi.NewInlineQueryResultArticle(
				fmt.Sprintf("%d", i),
				fmt.Sprintf("%s %s", snippet.Icon, snippet.Title),
				fmt.Sprintf("Page URL: %s\n", snippet.URL),
			),
		)
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       []interface{}{results},
	}

	if _, err := bot.Request(inlineConf); err != nil {
		slogger.Errorw(err.Error())
	}

	return
}