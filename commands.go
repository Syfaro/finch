package finch

import (
	"github.com/syfaro/telegram-bot-api"
	"regexp"
	"strings"
)

var Commands []Command

func RegisterCommand(cmd Command) {
	Commands = append(Commands, cmd)
}

func SimpleCommand(trigger string, message string) bool {
	return regexp.MustCompile("^/(" + trigger + ")(@\\w+)?( .+)?$").MatchString(message)
}

func SimpleArgCommand(trigger string, args int, message string) bool {
	matches := regexp.MustCompile("^/(" + trigger + ")(@\\w+)?( .+)?$").FindStringSubmatch(message)
	msgArgs := len(strings.Split(strings.Trim(matches[3], " "), " "))
	return args == msgArgs
}

func (f *Finch) CommandRouter(update tgbotapi.Update) {
	for _, command := range *f.Commands {
		if command.ShouldExecute(update) {
			err := command.Execute(update, f)
			f.commandError(update, err)
		}
	}
}
