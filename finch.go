// Package finch is a framework for Telegram Bots.
package finch

import (
	"encoding/json"
	"github.com/syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Config is a type used for storing configuration information.
type Config map[string]interface{}

// LoadConfig loads the saved config, if it exists.
//
// It looks for a FINCH_CONFIG environmental variable,
// before falling back to a file name config.json.
func LoadConfig() (*Config, error) {
	fileName := os.Getenv("FINCH_CONFIG")
	if fileName == "" {
		fileName = "config.json"
	}

	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return &Config{}, nil
	}

	var cfg Config
	json.Unmarshal(f, &cfg)

	return &cfg, nil
}

// Save saves the current Config struct.
//
// It uses the same file as LoadConfig.
func (c *Config) Save() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	fileName := os.Getenv("FINCH_CONFIG")
	if fileName == "" {
		fileName = "config.json"
	}

	return ioutil.WriteFile(fileName, b, 0600)
}

// Finch is a Telegram Bot, including API, Config, and Commands.
type Finch struct {
	API      *tgbotapi.BotAPI
	Config   Config
	Commands []*CommandState
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
	bot.Commands = commands

	c, _ := LoadConfig()
	bot.Config = *c

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

	_, err := f.API.SendMessage(message)
	return err
}
