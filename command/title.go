package command

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"

	hbot "github.com/whyrusleeping/hellabot"
	"golang.org/x/net/html"
)

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		if n.FirstChild != nil {
			rawTitle := n.FirstChild.Data
			actualTitle := strings.TrimSpace(rawTitle)
			return actualTitle, true
		} else {
			return "-- empty title --", true
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func GetHtmlTitle(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		panic("Fail to parse html")
	}

	return traverse(doc)
}

func DisplayHTMLTitle(bot *hbot.Bot, m *hbot.Message, url string) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		bot.Reply(m, "error with the cookiejar")
	}

	client := http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		bot.Reply(m, fmt.Sprintf("Http get error: %s", err))
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/81.0")
	resp, err := client.Do(req)
	if err != nil {
		bot.Reply(m, fmt.Sprintf("Http get error: %s", err))
	}
	defer resp.Body.Close()

	if title, ok := GetHtmlTitle(resp.Body); ok {
		bot.Reply(m, fmt.Sprintf("\x02%s", title))
	}
}
