package main

import (
	"bot/internal/models"

	"github.com/gocolly/colly"
)

func main() {
	var quotes []models.Quote

	c := colly.NewCollector(
		colly.AllowedDomains("https://ru.citaty.net"),
	)

}
