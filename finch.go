// Package finch is a framework for Telegram Bots.
package finch

import (
	"net/http"
	"strings"

	"github.com/getsentry/raven-go"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var sentryEnabled = false

// Finch is a Telegram Bot, including API, Config, and Commands.
type Finch struct {
	API      *tgbotapi.BotAPI
	Config   Config
	Commands []*CommandState
	Inline   InlineCommand
}

// NewFinch returns a new Finch instance, with Telegram API setup.
func NewFinch(token string) *Finch {
	return NewFinchWithClient(token, &http.Client{})
}

// NewFinchWithClient returns a new Finch instance,
// using a different net/http Client.
func NewFinchWithClient(token string, client *http.Client) *Finch {
	bot := &Finch{}

	c, _ := LoadConfig()
	bot.Config = *c

	if token == "" {
		val := bot.Config.Get("token")
		if val == nil {
			panic("no token provided")
		}

		token = val.(string)
	}

	api, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		panic(err)
	}

	bot.API = api
	bot.Commands = commands
	bot.Inline = inline

	return bot
}

// Start initializes commands, and starts listening for messages.
func (f *Finch) Start() {
	if v := f.Config.Get("sentry_dsn"); v != nil {
		sentryEnabled = true
		raven.SetDSN(v.(string))
	}

	if v := f.Config.Get("sentry_env"); v != nil {
		raven.SetEnvironment(v.(string))
	}

	f.commandInit()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 86400

	updates := f.API.GetUpdatesChan(u)

	for update := range updates {
		go f.commandRouter(update)
	}
}

// StartWebhook initializes commands,
// then registers a webhook for the bot to listen on
func (f *Finch) StartWebhook(endpoint string) {
	f.commandInit()

	f.API.ListenForWebhook(endpoint)
}

// SendMessage sends a message with various changes, and does not return the Message.
//
// At some point, this may do more handling as needed.
func (f *Finch) SendMessage(message tgbotapi.MessageConfig) error {
	message.Text = strings.Replace(message.Text, "@@", "@"+f.API.Self.UserName, -1)

	_, err := f.API.Send(message)
	if err != nil && sentryEnabled {
		raven.CaptureError(err, nil)
	}
	return err
}

// QuickReply quickly sends a message as a reply.
func (f *Finch) QuickReply(message tgbotapi.Message, text string) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID

	return f.SendMessage(msg)
}
