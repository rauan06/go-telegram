package parser

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseQute() (string, error) {
	url := "https://ru.citaty.net/tsitaty/sluchainaia-tsitata/"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	quote := ""
	author := ""
	doc.Find("h1.blockquote-display").Each(func(i int, s *goquery.Selection) {
		quote = strings.TrimSpace(s.Text())
	})
	doc.Find("div.blockquote-origin a").Each(func(i int, s *goquery.Selection) {
		author = strings.TrimSpace(s.Text())
	})

	if quote != "" && author != "" {
		return fmt.Sprintf("Цитата: %s\nАвтор: %s\n", quote, author), nil
	} else {
		return "", fmt.Errorf("quote or author not found")
	}
}
