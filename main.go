package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
	"golang.org/x/net/html"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const CHAT_ID int64 = 75808241
const BOT_TOKEN string = "347808432:AAFJQQOUDKCHFBaSxAbVCykyIMa-D9dCcE4"
const SPB_WEATHER_URL = "https://www.gismeteo.ru/weather-sankt-peterburg-4079/"

func main() {
	//TODO получение token id чата из конфига

	var wg sync.WaitGroup
	wg.Add(2)

	bot, err := tgbotapi.NewBotAPI(BOT_TOKEN)
	tgbotapi.NewCallback("asdf", "asdf")
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	go sheduler(bot, wg)
	go answerPhone(bot, wg)
	wg.Wait()
}
func sheduler(bot *tgbotapi.BotAPI, wg sync.WaitGroup) {
	//TODO получение списка из хранилища
	c := cron.New()

	c.AddFunc("0 40 05 * * *", func() {
		msg := tgbotapi.NewMessage(CHAT_ID, parseGismeteoWeather())
		bot.Send(msg)
	})

	//var tasks = []string{"1", "2", "3"}
	//for _, task := range tasks {
	//	fmt.Println("Add new cron")
	//	fmt.Println(task)
	//	c.AddFunc("0 05 20 * * *", func() {
	//		fmt.Println("new sheduled message")
	//		fmt.Println(task)
	//		msg := tgbotapi.NewMessage(CHAT_ID, test())
	//		bot.Send(msg)
	//	})
	//	//fmt.Println(task)
	//	//time.Sleep(time.Second * 3)
	//	//fmt.Println("new sheduled message")
	//	//msg := tgbotapi.NewMessage(CHAT_ID, task)
	//	//bot.Send(msg)
	//}
	c.Start()
	wg.Done()
}

func answerPhone(bot *tgbotapi.BotAPI, wg sync.WaitGroup) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.Chat.ID, update.Message.Text)
		if update.Message.Text == "погода" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, parseGismeteoWeather())
			//btn := tgbotapi.NewKeyboardButton("test")
			//row := tgbotapi.NewKeyboardButtonRow(btn)

			bot.Send(msg)
		}
	}
	wg.Done()
}

func parseGismeteoWeather() string {
	client := &http.Client{
		Timeout: time.Duration(time.Duration(5) * time.Second),
	}

	req, _ := http.NewRequest("GET", SPB_WEATHER_URL, nil)
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

	temp := "not Found"
	doc.Find(".js_meas_container.temperature").Each(func(index int, item *goquery.Selection) {
		title := item.Text()
		temp, _ = item.Attr("data-value")
		// linkTag := item.Find("a")
		// link, _ := linkTag.Attr("href")
		fmt.Printf("Post: %s - %s\n", title, temp)
	})

	return temp
}
