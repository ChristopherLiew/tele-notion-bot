package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"tele-notion-bot/internal/database"
	"tele-notion-bot/internal/notion"
)

// Search

// teleSearchQueryHandler handles all notion search related inline queries such as page and database search.
func TeleSearchQueryHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *viper.Viper, slogger *zap.SugaredLogger) (err error) {

	// Get notion auth token
	user := database.GetNotionUser(update.Message.From.UserName)
	if err != nil {
		slogger.Error(err.Error())
	}

	switch update.InlineQuery.Query {
	case "pages":
		err = queryPages(update, bot, user.Token, slogger, 50)
	case "databases":
		err = queryDBs(update, bot, user.Token, slogger, 50)
	default:
		err = fmt.Errorf("unknown callback: [%s]", update.CallbackQuery.Data)
	}

	return
}

// queryPages retrieves all notion pages from a notion workspace and returns results as telegram inline query result articles.
func queryPages(update tgbotapi.Update, bot *tgbotapi.BotAPI, notionIntToken string, slogger *zap.SugaredLogger, pageSize int) (err error) {

	response, err := notion.GetAllPages(notionIntToken, slogger, nil, pageSize)
	if err != nil {
		slogger.Errorw(err.Error())
		return err
	}
	snippets := notion.GetPageSnippets(response.Results, slogger)

	// construct inline query results
	results := make([]interface{}, len(snippets))
	for i, snippet := range snippets {
		results[i] = tgbotapi.NewInlineQueryResultArticleHTML(
			fmt.Sprintf("%d", i),
			fmt.Sprintf("%s %s", snippet.Icon, snippet.Title),
			fmt.Sprintf(`
				<b>Page</b>: %s %s
				<b>Link</b>: %s`,
				snippet.Icon, snippet.Title, snippet.URL),
		)
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       results,
	}

	slogger.Infof("Returning %d most recent page search results", pageSize)

	if _, err := bot.Request(inlineConf); err != nil {
		slogger.Errorw(err.Error())
	}

	return
}

// queryDBs retrieves all notion databases from a notion workspace and returns results as telegram inline query result articles.
func queryDBs(update tgbotapi.Update, bot *tgbotapi.BotAPI, notionIntToken string, slogger *zap.SugaredLogger, pageSize int) (err error) {

	response, err := notion.GetAllDBs(notionIntToken, slogger, nil, pageSize)
	if err != nil {
		slogger.Errorw(err.Error())
		return err
	}

	// construct inline query results
	results := make([]interface{}, len(response.Results))
	for i, db := range response.Results {

		icon := db.Icon.Emoji
		if icon == "" {
			icon = "💽"
		}
		title := db.Title[0].PlainText
		if title == "" {
			title = fmt.Sprintf(
				`Database has no title! (Id: %s)`,
				db.Id,
			)
		}
		url := db.URL

		results[i] = tgbotapi.NewInlineQueryResultArticleHTML(
			fmt.Sprintf("%d", i),
			fmt.Sprintf("%s %s", icon, title),
			fmt.Sprintf(`
				<b>Database</b>: %s %s
				<b>Link</b>: %s`, icon, title, url),
		)
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       results,
	}

	slogger.Infof("Returning %d most recent database search results", pageSize)

	if _, err := bot.Request(inlineConf); err != nil {
		slogger.Errorw(err.Error())
	}

	return
}
