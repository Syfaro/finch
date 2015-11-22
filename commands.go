package finch

import (
	"github.com/syfaro/telegram-bot-api"
	"log"
	"regexp"
	"strings"
)

var commands []*CommandState

// RegisterCommand adds a command to the bot.
func RegisterCommand(cmd Command) {
	commands = append(commands, &CommandState{
		Command:         cmd,
		WaitingForReply: false,
	})
}

// SimpleCommand generates a command regex and matches it against a message.
//
// The trigger is the command without the slash,
// and the message is the text to check it against.
func SimpleCommand(trigger, message string) bool {
	return regexp.MustCompile("^/(" + trigger + ")(@\\w+)?( .+)?$").MatchString(message)
}

// SimpleArgCommand generates a command regex and matches it against a message,
// requiring a certain number of parameters.
//
// The trigger is the command without the slash, args is number of arguments,
// and the message is the text to check it against.
func SimpleArgCommand(trigger string, args int, message string) bool {
	matches := regexp.MustCompile("^/(" + trigger + ")(@\\w+)?( .+)?$").FindStringSubmatch(message)
	msgArgs := len(strings.Split(strings.Trim(matches[3], " "), " "))
	return args == msgArgs
}

func (f *Finch) commandRouter(update tgbotapi.Update) {
	for _, command := range f.Commands {
		if command.WaitingForReply {
			err := command.Command.ExecuteKeyboard(update)
			f.commandError(update, err)
		}

		if command.Command.ShouldExecute(update) {
			err := command.Command.Execute(update)
			f.commandError(update, err)
		}
	}
}

func (f *Finch) commandInit() {
	for _, command := range f.Commands {
		err := command.Command.Init(command, f)
		if err != nil {
			log.Printf("Error starting plugin %s: %s\n", command.Command.Help().Name, err.Error())
		} else {
			log.Printf("Started plugin %s!", command.Command.Help().Name)
		}
	}
}

func (f *Finch) commandError(update tgbotapi.Update, err error) {
	if err == nil {
		return
	}

	var msg tgbotapi.MessageConfig

	if f.API.Debug {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "An error occured processing a command!")
	}

	msg.ReplyToMessageID = update.Message.MessageID

	_, err = f.API.Send(msg)
	if err != nil {
		log.Printf("An error happened processing an error!\n%s\n", err.Error())
	}
}
