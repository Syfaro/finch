package finchcommandcancel

import (
	"github.com/syfaro/finch"
	"gopkg.in/telegram-bot-api.v1"
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

	return cmd.Finch.SendMessage(msg)
}
