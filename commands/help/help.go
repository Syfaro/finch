package finchcommandhelp

import (
	"bytes"
	"github.com/syfaro/finch"
	"github.com/syfaro/telegram-bot-api"
)

func init() {
	finch.RegisterCommand(&helpCommand{})
}

type helpCommand struct {
}

func (cmd *helpCommand) Help() finch.Help {
	return finch.Help{
		Name:        "Help",
		Description: "Displays loaded commands and their help text",
		Example:     "/help@@",
	}
}

func (cmd *helpCommand) Init() error {
	return nil
}

func (cmd *helpCommand) ShouldExecute(update tgbotapi.Update) bool {
	return finch.SimpleCommand("help", update.Message.Text)
}

func (cmd *helpCommand) Execute(update tgbotapi.Update, f *finch.Finch) error {
	b := &bytes.Buffer{}

	b.WriteString("Loaded commands:\n\n")

	for _, command := range *f.Commands {
		help := command.Help()

		if help.Description == "" {
			continue
		}

		b.WriteString(help.String(true))
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, b.String())
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = tgbotapi.ModeMarkdown
	return f.SendMessage(msg)
}
