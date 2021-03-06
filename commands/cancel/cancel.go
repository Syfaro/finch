package commands

import (
	"github.com/Syfaro/finch"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	finch.RegisterCommand(&cancelCommand{})
}

type cancelCommand struct {
	finch.CommandBase
}

func (cancelCommand) ShouldExecute(message tgbotapi.Message) bool {
	return finch.SimpleCommand("cancel", message.Text)
}

func (cmd cancelCommand) Execute(message tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Your current command has been canceled.")
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      true,
	}

	for _, command := range cmd.Finch.Commands {
		command.ReleaseWaiting(message.From.ID)
	}

	return cmd.Finch.SendMessage(msg)
}

func (cmd cancelCommand) IsHighPriority(tgbotapi.Message) bool {
	return true
}

func (cancelCommand) Help() finch.Help {
	return finch.Help{
		Name: "Cancel",
		Botfather: [][]string{
			{"cancel", "Cancels the current command"},
		},
	}
}
