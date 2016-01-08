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

func (cmd *cancelCommand) ShouldExecute(update tgbotapi.Update) bool {
	return finch.SimpleCommand("cancel", update.Message.Text)
}

func (cmd *cancelCommand) Execute(update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your current command has been canceled.")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardHide{
		HideKeyboard: true,
		Selective:    true,
	}

	cmd.ReleaseWaiting(update.Message.From.ID)

	return cmd.Finch.SendMessage(msg)
}
