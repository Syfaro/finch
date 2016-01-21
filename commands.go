package finch

import (
	"gopkg.in/telegram-bot-api.v2"
	"log"
	"regexp"
	"strings"
)

var commands []*CommandState
var inline InlineCommand

// RegisterCommand adds a command to the bot.
func RegisterCommand(cmd Command) {
	commands = append(commands, NewCommandState(cmd))
}

// SetInline sets the Inline Query handler.
func SetInline(handler InlineCommand) {
	inline = handler
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
	if len(matches) < 4 {
		return false
	}
	msgArgs := len(strings.Split(strings.Trim(matches[3], " "), " "))
	return args == msgArgs
}

func (f *Finch) commandRouter(update tgbotapi.Update) {
	if update.InlineQuery.ID != "" {
		if f.Inline == nil {
			log.Println("Got inline query, but no handler is set!")
			return
		}

		if err := f.Inline.Execute(f, update); err != nil {
			log.Printf("Error processing inline query:\n%+v\n", err)
		}
		return
	}

	for _, command := range f.Commands {
		if command.IsWaiting(update.Message.From.ID) {
			err := command.Command.ExecuteKeyboard(update.Message)
			f.commandError(update.Message, err)
			return
		}

		if command.Command.ShouldExecute(update.Message) {
			err := command.Command.Execute(update.Message)
			f.commandError(update.Message, err)
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

func (f *Finch) commandError(message tgbotapi.Message, err error) {
	if err == nil {
		return
	}

	var msg tgbotapi.MessageConfig

	if f.API.Debug {
		msg = tgbotapi.NewMessage(message.Chat.ID, err.Error())
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, "An error occured processing a command!")
	}

	msg.ReplyToMessageID = message.MessageID

	_, err = f.API.Send(msg)
	if err != nil {
		log.Printf("An error happened processing an error!\n%s\n", err.Error())
	}
}
