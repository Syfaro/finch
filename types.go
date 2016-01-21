package finch

import (
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Help contains information about a command,
// used for showing info in the help command.
type Help struct {
	Name        string
	Description string
	Example     string
	Botfather   [][]string
}

// String converts a Help struct into a pretty string.
//
// full makes each command item multiline with extra newlines.
func (h Help) String(full bool) string {
	b := &bytes.Buffer{}

	b.WriteString(h.Name)
	if full {
		b.WriteString("\n")
	} else {
		b.WriteString(" - ")
	}
	b.WriteString(h.Description)
	b.WriteString("\n")

	if full {
		b.WriteString("Example: ")
		b.WriteString(h.Example)
		b.WriteString("\n")
	}

	b.WriteString("\n")

	return b.String()
}

// BotfatherString formats a Help struct into something for Botfather.
func (h Help) BotfatherString() string {
	if len(h.Botfather) == 0 {
		return ""
	}

	b := bytes.Buffer{}

	for k, v := range h.Botfather {
		b.WriteString(v[0])
		b.WriteString(" - ")
		b.WriteString(v[1])
		if k+1 != len(h.Botfather) {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// Command contains the methods a must have.
type Command interface {
	Help() Help
	Init(*CommandState, *Finch) error
	ShouldExecute(tgbotapi.Message) bool
	Execute(tgbotapi.Message) error
	ExecuteKeyboard(tgbotapi.Message) error
}

// CommandBase is a default Command that handles various tasks for you,
// and allows for you to not have to write empty methods.
type CommandBase struct {
	*CommandState
	*Finch
}

// Help returns an empty Help struct.
func (CommandBase) Help() Help { return Help{} }

// Init sets MyState equal to the current CommandState.
//
// If you overwrite this method, you should still set MyState equal to CommandState!
func (cmd *CommandBase) Init(c *CommandState, f *Finch) error {
	cmd.CommandState = c
	cmd.Finch = f

	return nil
}

// ShouldExecute returns false, you should overwrite this method.
func (CommandBase) ShouldExecute(tgbotapi.Message) bool { return false }

// Execute returns nil to show no error, you should overwrite this method.
func (CommandBase) Execute(tgbotapi.Message) error { return nil }

// ExecuteKeyboard returns nil to show no error, you may overwrite this
// when you are expecting to get a reply that is not a command.
func (CommandBase) ExecuteKeyboard(tgbotapi.Message) error { return nil }

// Get fetches an item from the Config struct.
func (cmd CommandBase) Get(key string) interface{} {
	return cmd.Finch.Config[key]
}

// Set sets an item in the Config struct, then saves it.
func (cmd CommandBase) Set(key string, value interface{}) {
	cmd.Finch.Config[key] = value
	cmd.Finch.Config.Save()
}

type userWait map[int]bool

// CommandState is the current state of a command.
// It contains the command and if the command is waiting for a reply.
type CommandState struct {
	Command             Command
	waitingForReplyUser userWait
}

// NewCommandState creates a new CommandState with an initialized map.
func NewCommandState(cmd Command) *CommandState {
	return &CommandState{
		Command:             cmd,
		waitingForReplyUser: make(userWait),
	}
}

// IsWaiting checks if the current CommandState is waiting for input from
// this user.
func (state *CommandState) IsWaiting(user int) bool {
	if v, ok := state.waitingForReplyUser[user]; ok {
		return v
	}

	return false
}

// SetWaiting sets that the bot should expect user input from this user.
func (state *CommandState) SetWaiting(user int) {
	state.waitingForReplyUser[user] = true
}

// ReleaseWaiting sets that the bot should not expect any input from
// this user.
func (state *CommandState) ReleaseWaiting(user int) {
	state.waitingForReplyUser[user] = false
}

// InlineCommand is a single command executed for an Inline Query.
type InlineCommand interface {
	Execute(*Finch, tgbotapi.Update) error
}
