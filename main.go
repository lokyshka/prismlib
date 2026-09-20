package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/max-messenger/maxbot"
	tele "gopkg.in/telebot.v4"
)

/*
func tgtest(tg tele.Context) error {
	return tg.Send("ping-pong!")
}

func maxtest(ctx maxbot.Context) error {
	return ctx.Reply("ping-pong!")
}
*/

var tgscn map[int64]int8 = make(map[int64]int8)
var tgmutex sync.RWMutex
var tg *tele.Bot

const selrole int8 = 0
const idle int8 = 1

const clselreg int8 = 2
const clselmun int8 = 3
const clsellib int8 = 4
const claskbook int8 = 5
const cledittimeget int8 = 6
const cledittimeret int8 = 7
const cllist int8 = 8

const libselreg int8 = -1
const libselmul int8 = -2
const libselname int8 = -3
const libseladress int8 = -4

func handletg(ctx tele.Context) {
	const welcomemsg string = "Привет! Я — бот, который поможет забронировать книгу в выбранной библиотеке! Выберите роль, чтобы продолжить:"

	userid := ctx.Sender().ID
	tgmutex.RLock()
	scene := tgscn[userid]
	tgmutex.RUnlock()

	switch scene {
	case selrole:
		selector := &tele.ReplyMarkup{}
		btncl := selector.Data("клиент", "client")
		btnlib := selector.Data("библиотекарь", "library")

		selector.Inline(
			selector.Row(btncl, btnlib),
		)

		ctx.Send(welcomemsg, selector)
		tg.Handle(&btncl, func(c tele.Context) error {
			tgscn[userid] = clselreg
			return nil
		})
		tg.Handle(&btncl, func(c tele.Context) error {
			tgscn[userid] = libselreg
			return nil
		})

	case clselreg:

	}

}

func handlemax(ctx maxbot.Context) {

}

func main() {
	// загрузка токенов
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env файл не найден!")
	}

	tgtoken := os.Getenv("tgtoken")
	if tgtoken == "" {
		log.Fatal("токен telegram пуст!")
	}

	maxtoken := os.Getenv("maxtoken")
	if maxtoken == "" {
		log.Fatal("токен max пуст!")
	}

	// запускаем бота telegram
	pref := tele.Settings{
		Token:  tgtoken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	tg, err = tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// запускаем бота max
	/*
		max, err := maxbot.NewApi(maxtoken)
		if err != nil {
			log.Fatal(err)
		}
	*/

	tg.Handle("/start", func(ctx tele.Context) error {
		userid := ctx.Sender().ID
		tgmutex.Lock()
		tgscn[userid] = selrole
		tgmutex.Unlock()
		log.Println("new /start")

		handletg(ctx)
		return nil
	})
	/* max.Handle("/start", func(ctx maxbot.Context) error {
		return nil
	}) */

	tg.Start()
	// max.Start()
}
