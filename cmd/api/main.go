package main

import (
	"bot/internal/parser"
	"log"
	"os"
	"time"

	"gopkg.in/telebot.v3"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN") // Set your bot token as an environment variable
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/help", func(c telebot.Context) error {
		return c.Send("Использование: введите команду \"/получить/\"")
	})

	bot.Handle("/получить", func(c telebot.Context) error {
		quote, err := parser.ParseQute()
		if err != nil {
			return c.Send("Couldn't find any quotes")
		}
		return c.Send(quote)
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		return c.Send("Использование: введите команду \"/получить\"")
	})

	log.Println("Bot is running...")
	bot.Start()
}
