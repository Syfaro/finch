package finch

import (
	"github.com/syfaro/telegram-bot-api"
	"log"
	"net/http"
	"strings"
)

type Config map[string]interface{}

type Finch struct {
	API      *tgbotapi.BotAPI
	Config   Config
	Commands *[]Command
}

func NewFinch(token string) *Finch {
	return NewFinchWithClient(token, &http.Client{})
}

func NewFinchWithClient(token string, client *http.Client) *Finch {
	bot := &Finch{}

	api, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		panic(err)
	}

	bot.API = api
	bot.Commands = &Commands

	bot.Config = make(Config)

	return bot
}

func (f *Finch) Start() {
	for _, command := range *f.Commands {
		err := command.Init()
		if err != nil {
			log.Printf("Error starting plugin %s: %s\n", command.Help().Name, err.Error())
		} else {
			log.Printf("Started plugin %s!", command.Help().Name)
		}
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 86400

	err := f.API.UpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range f.API.Updates {
		f.CommandRouter(update)
	}
}

func (f *Finch) SendMessage(message tgbotapi.MessageConfig) error {
	message.Text = strings.Replace(message.Text, "@@", "@"+f.API.Self.UserName, -1)

	_, err := f.API.SendMessage(message)
	return err
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

	_, err = f.API.SendMessage(msg)
	if err != nil {
		log.Printf("An error happened processing an error!\n%s\n", err.Error())
	}
}
