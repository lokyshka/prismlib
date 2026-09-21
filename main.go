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

var regions = [...]string{
	// республики
	"республика адыгея",
	"республика алтай",
	"республика башкортостан",
	"республика бурятия",
	"республика дагестан",
	"донецкая народная республика",
	"республика ингушетия",
	"кабардино-балкарская республика",
	"республика калмыкия",
	"карачаево-черкесская республика",
	"республика карелия",
	"республика коми",
	"республика крым",
	"луганская народная республика",
	"республика марий эл",
	"республика мордовия",
	"республика саха",
	"республика северная осетия",
	"республика татарстан",
	"республика тыва",
	"удмуртская республика",
	"республика хакасия",
	"чеченская республика",
	"чувашская республика",

	// края
	"алтайский край",
	"забайкальский край",
	"камчатский край",
	"краснодарский край",
	"красноярский край",
	"пермский край",
	"приморский край",
	"ставропольский край",
	"хабаровский край",

	// области
	"амурская область",
	"архангельская область",
	"астраханская область",
	"белгородская область",
	"брянская область",
	"владимирская область",
	"волгоградская область",
	"вологодская область",
	"воронежская область",
	"ивановская область",
	"иркутская область",
	"калининградская область",
	"калужская область",
	"кемеровская область",
	"кировская область",
	"костромская область",
	"курганская область",
	"курская область",
	"ленинградская область",
	"липецкая область",
	"магаданская область",
	"московская область",
	"мурманская область",
	"нижегородская область",
	"новгородская область",
	"новосибирская область",
	"омская область",
	"оренбургская область",
	"орловская область",
	"пензенская область",
	"псковская область",
	"ростовская область",
	"рязанская область",
	"самарская область",
	"саратовская область",
	"сахалинская область",
	"свердловская область",
	"смоленская область",
	"тамбовская область",
	"тверская область",
	"томская область",
	"тульская область",
	"тюменская область",
	"ульяновская область",
	"челябинская область",
	"ярославская область",
	"запорожская область",
	"херсонская область",

	// города федерального значения
	"москва",
	"санкт-петербург",
	"севастополь",

	// автономная область
	"еврейская автономная область",

	// автономные округа
	"ненецкий автономный округ",
	"ханты-мансийский автономный округ",
	"чукотский автономный округ",
	"ямало-ненецкий автономный округ",
}

type user struct {
	sc     int8
	region int8
	city   string
	mu     sync.Mutex
}

var users = make(map[int64]*user)
var tgmutex sync.RWMutex
var tg *tele.Bot

const selrole int8 = 0
const idle int8 = 1

const clselreg int8 = 2
const clselcity int8 = 3
const clsellib int8 = 4
const clafter int8 = 5
const claskbook int8 = 6
const cledittimeget int8 = 7
const cledittimeret int8 = 8
const cllist int8 = 9

const libselreg int8 = -1
const libselmul int8 = -2
const libselname int8 = -3
const libseladress int8 = -4

func handletg(ctx tele.Context) {
	const welcomemsg string = "Привет! Я — бот, который поможет забронировать книгу в выбранной библиотеке! Выберите роль, чтобы продолжить:"

	userid := ctx.Sender().ID
	tgmutex.RLock()
	scene := users[userid].sc
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
			users[userid].sc = clselreg
			handletg(ctx)
			return nil
		})
		tg.Handle(&btncl, func(c tele.Context) error {
			users[userid].sc = libselreg
			handletg(ctx)
			return nil
		})

	case clselreg:
		ctx.Send("Введите название региона. Пример: \"Республика Татарстан\", \"Москва\".")
		users[userid].sc = clselcity

	case clselcity:
		var flag bool
		region := ctx.Message().Text

		for _, item := range regions {
			if item == region {
				flag = true
				break
			}
		}

		if !flag {
			ctx.Send("Неверное название региона! Попробуйте еще раз в другом формате.")
			return
		}

		// привязываем человека к региону

		ctx.Send("Отлично! Теперь, введите город.")
		users[userid].sc = clsellib

	case clsellib:
		// привязываем человека к городу

		ctx.Send("Выберите библиотеку. Нету в списке? Возможно, библиотека не зарегистрирована в боте или вы неверно указали город.")
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

	/*
		// запускаем бота max
		max, err := maxbot.NewApi(maxtoken)
		if err != nil {
			log.Fatal(err)
		}
	*/

	tg.Handle("/start", func(ctx tele.Context) error {
		userid := ctx.Sender().ID
		tgmutex.Lock()
		users[userid].sc = selrole
		tgmutex.Unlock()
		log.Println("new /start")

		handletg(ctx)
		return nil
	})
	/*
		max.Handle("/start", func(ctx maxbot.Context) error {
			return nil
		})
		go max.Start()
	*/

	tg.Start()
}
