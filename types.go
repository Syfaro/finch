package finch

import (
	"bytes"
	"github.com/syfaro/telegram-bot-api"
)

type Help struct {
	Name        string
	Description string
	Example     string
}

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

type Command interface {
	Help() Help
	Init() error
	ShouldExecute(tgbotapi.Update) bool
	Execute(tgbotapi.Update, *Finch) error
}
