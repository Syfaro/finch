package finchcommandinfo

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/syfaro/finch"
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
		Botfather: [][]string{
			[]string{"info", "Information about the current user"},
		},
	}
}

func (cmd *infoCommand) ShouldExecute(update tgbotapi.Update) bool {
	return finch.SimpleCommand("info", update.Message.Text)
}

func (cmd *infoCommand) Execute(update tgbotapi.Update) error {
	text := fmt.Sprintf("Your ID: %d\nChat ID: %d", update.Message.From.ID, update.Message.Chat.ID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyToMessageID = update.Message.MessageID

	return cmd.Finch.SendMessage(msg)
}
