package notion

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func HandleBotCommands(update tgbotapi.Update, bot *tgbotapi.BotAPI, slogger *zap.SugaredLogger) {

	// slogger.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	// msg.ReplyToMessageID = update.Message.MessageID

	// bot.Send(msg)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	// Extract the command from the Message.
	switch update.Message.Command() {
	case "help":
		msg.Text = "I understand /start and /search nothing else really matters ;)."
	case "start":
		msg.Text = "Click on the button below to connect to your Notion Account!"
	case "search":
		msg.Text = "What Page or Database might you be looking for? "
	default:
		msg.Text = "讲啥你"
	}

	if _, err := bot.Send(msg); err != nil {
		slogger.Panic(err)
	}

}

// Handlers for each command (Move to individual files if too large)
