package main

import (
	"bot/internal/parser"
	"bot/internal/request"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

const (
	helpText = `1. Команда /start – запускает чат-бота, выводит краткое описание его
возможностей и инструкции по использованию. Это помогает
пользователю быстро разобраться с интерфейсом и начать работу.
2. Команда /help – предоставляет справочную информацию о доступных
командах, объясняет, как вводить запросы и какие возможности
предоставляет бот. Также содержит примеры запроса и ответа.
3. Команда /quote – позволяет пользователю ввести ключевые слова или
тему, после чего бот предоставляет релевантные литературные цитаты.
Это особенно полезно для студентов, преподавателей и
исследователей, нуждающихся в точных формулировках и примерах из
классической литературы.
4. Команда /context – анализирует контекст цитаты и предоставляет
дополнительную информацию о произведении, его авторе и
историческом значении цитируемого фрагмента.
5. Дополнительные команды, такие как /random_quote, позволяют
пользователям изучать случайные цитаты, что способствует
расширению литературного кругозора.`
	context = `проанализируй контекст цитаты и предоставь
дополнительную информацию о произведении, его авторе и
историческом значении цитируемого фрагмента, уложи свой ответ в не более чем 50 слов:`
	query = `Предоставь релевантные литературные цитаты по этим ключевым словам, уложи свой ответ в не более чем 50 слов:`
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN") // Set your bot token as an environment variable
	api_token := strings.Trim(os.Getenv("API_KEY"), " ")

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
		return c.Send(helpText)
	})

	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Здравствуйте, чем могу вам помочь?")
	})

	bot.Handle("/random_quote", func(c telebot.Context) error {
		quote, err := parser.ParseQute()
		if err != nil {
			return c.Send("Couldn't find any quotes")
		}
		return c.Send(quote)
	})

	bot.Handle("/quote", func(c telebot.Context) error {
		messageText := c.Text()

		if len(messageText) > 0 {
			answer, err := request.SendRequest(query, messageText, api_token)
			if err != nil {
				log.Println(err)
				return c.Send("Произошла ошибка на стороне сервера.")
			}
			return c.Send(answer)
		}
		return c.Send("Пожалуйста, введите текст после команды /quert [текст].")
	})

	bot.Handle("/context", func(c telebot.Context) error {
		messageText := c.Text()

		if len(messageText) > 0 {
			answer, err := request.SendRequest(context, messageText, api_token)
			if err != nil {
				log.Println(err)
				return c.Send("Произошла ошибка на стороне сервера.")
			}
			return c.Send(answer)
		}
		return c.Send("Пожалуйста, введите текст после команды /context [текст].")
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		return c.Send(helpText)
	})

	log.Println("Bot is running...")
	go bot.Start()
	http.ListenAndServe(":8080", http.HandlerFunc(http.NotFound))
}
