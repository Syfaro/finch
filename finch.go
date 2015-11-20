// Package finch is a framework for Telegram Bots.
package finch

import (
	"github.com/syfaro/telegram-bot-api"
	"log"
	"net/http"
	"strings"
)

// Config is a type used for storing configuration information.
type Config map[string]interface{}

// Finch is a Telegram Bot, including API, Config, and Commands.
type Finch struct {
	API      *tgbotapi.BotAPI
	Config   Config
	Commands *[]Command
}

// NewFinch returns a new Finch instance, with Telegram API setup.
func NewFinch(token string) *Finch {
	return NewFinchWithClient(token, &http.Client{})
}

// NewFinchWithClient returns a new Finch instance,
// using a different net/http Client.
func NewFinchWithClient(token string, client *http.Client) *Finch {
	bot := &Finch{}

	api, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		panic(err)
	}

	bot.API = api
	bot.Commands = &commands

	bot.Config = make(Config)

	return bot
}

// Start initializes commands, and starts listening for messages.
func (f *Finch) Start() {
	f.commandInit()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 86400

	err := f.API.UpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range f.API.Updates {
		f.commandRouter(update)
	}
}

// StartWebhook initializes commands,
// then registers a webhook for the bot to listen on.
//
// This webhook URL is your bot token.
func (f *Finch) StartWebhook() {
	f.commandInit()

	f.API.ListenForWebhook("/" + f.API.Token)
}

// SendMessage sends a message with various changes, and does not return the Message.
//
// At some point, this may do more handling as needed.
func (f *Finch) SendMessage(message tgbotapi.MessageConfig) error {
	message.Text = strings.Replace(message.Text, "@@", "@"+f.API.Self.UserName, -1)

	_, err := f.API.SendMessage(message)
	return err
}
