package finchcommandcancel

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/syfaro/finch"
)

func init() {
	finch.RegisterCommand(&cancelCommand{})
}

type cancelCommand struct {
	finch.CommandBase
}

func (cmd *cancelCommand) ShouldExecute(message tgbotapi.Message) bool {
	return finch.SimpleCommand("cancel", message.Text)
}

func (cmd *cancelCommand) Execute(message tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Your current command has been canceled.")
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardHide{
		HideKeyboard: true,
		Selective:    true,
	}

	cmd.ReleaseWaiting(message.From.ID)

	return cmd.Finch.SendMessage(msg)
}

func (cmd *cancelCommand) Help() finch.Help {
	return finch.Help{
		Name: "Cancel",
		Botfather: [][]string{
			[]string{"cancel", "Cancels the current command"},
		},
	}
}
