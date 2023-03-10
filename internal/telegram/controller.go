package telegram

import (
	"fmt"
	"tele-notion-bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// session variables (think of how to better cache / store these?)
var TeleUserName string

// teleCommandHandler handles all telegram commands from the user.
func TeleCommandHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *viper.Viper, slogger *zap.SugaredLogger) {

	switch update.Message.Command() {
	case "start":
		startCommand(update, bot, cfg, slogger)
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

// startCommand authenticates and connects the bot to the user's notion workspace.
func startCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *viper.Viper, slogger *zap.SugaredLogger) {

	// Start your connection
	TeleUserName = update.Message.From.UserName
	slogger.Infof("Begin authorisation to %s's notion workspace", TeleUserName)

	authURL := cfg.GetString("NOTION.AUTHORIZATION_URL")
	authKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text: "Authenticate Notion 🔑",
				URL:  &authURL,
			},
		),
	)

	// check if user with key exists else proceed with oauth2
	user := database.GetNotionUser(TeleUserName)
	if user.UserName == "" {
		slogger.Infof("%s not previously authenticated, starting OAuth2 auth flow", TeleUserName)
		authMsg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Welcome to NotionBot 😊\nLet's start by connecting your Workspace!",
		)
		authMsg.ReplyMarkup = authKeyboard
		if _, err := bot.Send(authMsg); err != nil {
			slogger.Error(err.Error())
		}
	} else {
		slogger.Infof("User previously authenticated as %s at timestamp: %s", user.UserName, user.Timestamp)
		completeAuthMsg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Welcome to NotionBot 😊\nYou have already authenticated!",
		)
		if _, err := bot.Send(completeAuthMsg); err != nil {
			slogger.Error(err.Error())
		}
	}
}

// searchCommand serves as the entrypoint to a user's notion search request.
func searchCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *viper.Viper, slogger *zap.SugaredLogger) {

	getPagesQuery := "pages"
	getDBQuery := "databases"

	user := database.GetNotionUser(TeleUserName)
	if user.UserName == "" {
		slogger.Infof("%s not previously authenticated, starting OAuth2 auth flow", TeleUserName)
		authMsg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Please hit /start and authenticate your Notion Workspace",
		)
		if _, err := bot.Send(authMsg); err != nil {
			slogger.Error(err.Error())
		}
		return
	}

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

// helpCommand provides basic help information to the user.
func helpCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	// Provide help and link to any tutorials
	helpMsg := `
	Hi there! I can ...

	/start 	- Connect with your Notion workspace 🔗
	/search - Search for Pages & Databases within it 🔎
	/end 	- End this conversation 🛑

	See the Menu for all possible commands
	`
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMsg)
	bot.Send(msg)
}

// endCommand terminates the conversation between the user and the telegram bot.
func endCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, slogger *zap.SugaredLogger) {

	userName := update.Message.From.UserName
	userId := update.Message.From.ID
	exitMsg := fmt.Sprintf(`Good bye <a href="tg://user?id=%d">@%s</a> 👋🏻`, userId, userName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, exitMsg)
	msg.ParseMode = "HTML"
	bot.Send(msg)

	slogger.Infof("Terminating session for tele user: %s", userName)
	bot.StopReceivingUpdates()
}
