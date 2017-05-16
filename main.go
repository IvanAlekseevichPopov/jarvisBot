package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {
	// parseGismeteo()
	// parseTest()
	bot, err := tgbotapi.NewBotAPI("347808432:AAFJQQOUDKCHFBaSxAbVCykyIMa-D9dCcE4")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, parseTest())
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

// func parseGismeteo() {
// 	doc, err := goquery.NewDocument("http://ivan-popov.tk")
// 	// doc, err := goquery.
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// use CSS selector found with the browser inspector
// 	// for each, use index and item
// 	doc.Find(".logo").Each(func(index int, item *goquery.Selection) {
// 		title := item.Text()
// 		linkTag := item.Find("a")
// 		link, _ := linkTag.Attr("href")
// 		fmt.Printf("Post #%d: %s - %s\n", index, title, link)

// 	})
// }

func parseTest() string {
	timeOut := 10
	client := &http.Client{
		Timeout: time.Duration(time.Duration(timeOut) * time.Second),
		// Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	}

	req, _ := http.NewRequest("GET", "https://www.gismeteo.ru/weather-sankt-peterburg-4079/", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/20100101 Firefox/51.0")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		//TODO
		return "err"
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%T\n", body)

	bodyString := string(body)
	fmt.Printf("%T\n", bodyString)
	// fmt.Printf("%s\n", bodyString)

	reader := strings.NewReader(bodyString)
	fmt.Printf("%T\n", reader)

	node, err := html.Parse(reader)
	fmt.Printf("%T\n", node)

	doc := goquery.NewDocumentFromNode(node)
	// doc, err := goquery.
	if err != nil {
		log.Fatal(err)
	}

	// use CSS selector found with the browser inspector
	// for each, use index and item
	temp := "not FOund"
	doc.Find(".js_meas_container.temperature").Each(func(index int, item *goquery.Selection) {
		title := item.Text()
		temp, _ = item.Attr("data-value")
		// linkTag := item.Find("a")
		// link, _ := linkTag.Attr("href")
		fmt.Printf("Post: %s - %s\n", title, temp)
		// return temp
	})

	fmt.Printf("%T\n", temp)
	fmt.Println(temp)
	// return string(temp)
	return temp
	// return "not Found", "adsf"
}
