package main

import (
	"bot/internal/parser"
	"bot/internal/request"
	"bot/logger"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

const (
	helpText = `/start – Запускает чат-бота, дает краткое описание возможностей и инструкции.
/help – Предоставляет справочную информацию о командах, запросах и примеры.
/quote <ключевые слова> – Позволяет вводить ключевые слова для получения литературных цитат.
/context <текст цитаты> – Анализирует цитату, предоставляя информацию о произведении и авторе.
/random_quote – Позволяет изучать случайные цитаты для расширения кругозора.`
	context = `проанализируй контекст цитаты и предоставь
дополнительную информацию о произведении, его авторе и
историческом значении цитируемого фрагмента, уложи свой ответ в не более чем 20 слов:`
	query = `Предоставь релевантные литературные цитаты по этим ключевым словам, уложи свой ответ в не более чем 20 слов:`
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN") // Set your bot token as an environment variable
	api_token := strings.Trim(os.Getenv("API_KEY"), " ")

	logger := logger.SetupPrettySlog(os.Stdout)
	slog.SetDefault(logger)

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

	// Log when the bot starts
	log.Println("Bot is starting...")

	bot.Handle("/help", func(c telebot.Context) error {
		log.Println("Received /help command")
		return c.Send(helpText)
	})

	bot.Handle("/start", func(c telebot.Context) error {
		log.Println("Received /start command")
		go request.SendRequest("", "", api_token)
		go parser.ParseQute()
		return c.Send("Бот готов принимать ваши запросы")
	})

	bot.Handle("/random_quote", func(c telebot.Context) error {
		log.Println("Received /random_quote command")
		quote, err := parser.ParseQute()
		if err != nil {
			log.Println("Error fetching random quote:", err)
			return c.Send("Couldn't find any quotes")
		}
		return c.Send(quote)
	})

	bot.Handle("/quote", func(c telebot.Context) error {
		log.Println("Received /quote command")
		messageText := c.Text()

		if len(messageText) > 0 {
			answer, err := request.SendRequest(query, messageText, api_token)
			if err != nil {
				log.Println("Error sending request for /quote:", err)
				return c.Send("Произошла ошибка на стороне сервера.")
			}
			return c.Send(answer)
		}
		return c.Send("Пожалуйста, введите текст после команды /quert [текст].")
	})

	bot.Handle("/context", func(c telebot.Context) error {
		log.Println("Received /context command")
		messageText := c.Text()

		if len(messageText) > 0 {
			answer, err := request.SendRequest(context, messageText, api_token)
			if err != nil {
				log.Println("Error sending request for /context:", err)
				return c.Send("Произошла ошибка на стороне сервера.")
			}
			return c.Send(answer)
		}
		return c.Send("Пожалуйста, введите текст после команды /context [текст].")
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		log.Println("Received text message")
		return c.Send(helpText)
	})

	log.Println("Bot is running...")
	go bot.Start()
	http.ListenAndServe(":8080", http.HandlerFunc(http.NotFound))
}
