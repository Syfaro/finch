package finchcommandstats

import (
	"bytes"
	"github.com/syfaro/finch"
	"github.com/syfaro/telegram-bot-api"
	"strconv"
)

type UserMessageCount map[string]int

var userMessages UserMessageCount

func init() {
	finch.RegisterCommand(&infoCollector{})
	finch.RegisterCommand(&infoCommand{})
	userMessages = make(UserMessageCount)
}

type infoCommand struct {
}

func (cmd *infoCommand) Help() finch.Help {
	return finch.Help{
		Name:        "Stats",
		Description: "Displays some statistics",
		Example:     "/stats@@",
	}
}

func (cmd *infoCommand) Init() error {
	return nil
}

func (cmd *infoCommand) ShouldExecute(update tgbotapi.Update) bool {
	return finch.SimpleCommand("stats", update.Message.Text)
}

func (cmd *infoCommand) Execute(update tgbotapi.Update, f *finch.Finch) error {
	b := &bytes.Buffer{}

	b.WriteString("Users seen\n\n")

	for userName, count := range userMessages {
		b.WriteString(userName)
		b.WriteString(" - ")
		b.WriteString(strconv.Itoa(count))
		b.WriteString("\n")
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, b.String())
	msg.ReplyToMessageID = update.Message.MessageID

	return f.SendMessage(msg)
}

type infoCollector struct {
}

func (cmd *infoCollector) Help() finch.Help {
	return finch.Help{Name: "Stats Collector"}
}

func (cmd *infoCollector) Init() error {
	return nil
}

func (cmd *infoCollector) ShouldExecute(update tgbotapi.Update) bool {
	return true
}

func (cmd *infoCollector) Execute(update tgbotapi.Update, f *finch.Finch) error {
	if _, ok := userMessages[update.Message.From.String()]; !ok {
		userMessages[update.Message.From.String()] = 1
	} else {
		userMessages[update.Message.From.String()] += 1
	}

	return nil
}
