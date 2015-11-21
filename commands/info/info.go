package finchcommandinfo

import (
	"fmt"
	"github.com/syfaro/finch"
	"github.com/syfaro/telegram-bot-api"
)

func init() {
	finch.RegisterCommand(&infoCommand{})
}

type infoCommand struct {
	finch.CommandBase
}

func (cmd *infoCommand) Help() finch.Help {
	return finch.Help{
		Name:        "Info",
		Description: "Displays information about the currently requesting user",
		Example:     "/info@@",
	}
}

func (cmd *infoCommand) ShouldExecute(update tgbotapi.Update) bool {
	return finch.SimpleCommand("info", update.Message.Text)
}

func (cmd *infoCommand) Execute(update tgbotapi.Update, f *finch.Finch) error {
	text := fmt.Sprintf("Your ID: %d\nChat ID: %d", update.Message.From.ID, update.Message.Chat.ID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyToMessageID = update.Message.MessageID

	return f.SendMessage(msg)
}
