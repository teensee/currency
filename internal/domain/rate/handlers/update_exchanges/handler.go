package update_exchanges

import (
	service2 "Currency/internal/domain/rate/service"
	"Currency/internal/infrastructure/service"
	"log"
	"net/http"
	"sync"
	"time"
)

const UpdateRateHandlerTag = "updExchangeH"

type UpdateExchangeHandler struct {
	srv            *service2.ExchangeRateService
	queryExtractor *service.QueryExtractor
}

func NewHandler(srv *service2.ExchangeRateService) *UpdateExchangeHandler {
	return &UpdateExchangeHandler{
		srv:            srv,
		queryExtractor: service.NewQueryExtractor(),
	}
}

// GetCbrExchangeRates Creates a CBR rates in database
func (h UpdateExchangeHandler) GetCbrExchangeRates(_ http.ResponseWriter, r *http.Request) {
	onDate := h.queryExtractor.ExtractDateParam(r.URL.Query())
	h.srv.GetNewExchangeRateOnDate(onDate.OnDate)
}

// SyncRatesOnStartup Синхронизация курсов валют при старте приложения
func (h UpdateExchangeHandler) SyncRatesOnStartup() {
	now := time.Now()
	y, m, d := now.Date()
	nowFtm := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	log.Printf("Syncing rates on %s", nowFtm.Format("02.01.2006"))

	h.srv.GetNewExchangeRateOnDate(nowFtm)
}

// RushRates Скачивание курсов валют из цб с 01/07/1992 по настоящее время
func (h UpdateExchangeHandler) RushRates(_ http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup

	startFromDate, _ := time.Parse("02/01/2006", "01/07/1992")
	now := time.Now()

	first := startFromDate.AddDate(0, 0, 1)
	sec := startFromDate.AddDate(0, 0, 2)
	third := startFromDate.AddDate(0, 0, 3)
	fourth := startFromDate.AddDate(0, 0, 4)
	fifth := startFromDate.AddDate(0, 0, 5)

	addDays := 0

	for now.Unix() > startFromDate.Unix() {
		wg.Add(5)

		first = first.AddDate(0, 0, addDays)
		sec = sec.AddDate(0, 0, addDays)
		third = third.AddDate(0, 0, addDays)
		fourth = fourth.AddDate(0, 0, addDays)
		fifth = fifth.AddDate(0, 0, addDays)
		log.Printf("Syncing rates on:\n%s,\n%s,\n%s,\n%s,\n%s\n", first.Format("02/01/2006"), sec.Format("02/01/2006"), third.Format("02/01/2006"), fourth.Format("02/01/2006"), fifth.Format("02/01/2006"))

		go h.srv.AsyncGetNewExchange(first, &wg)
		go h.srv.AsyncGetNewExchange(sec, &wg)
		go h.srv.AsyncGetNewExchange(third, &wg)
		go h.srv.AsyncGetNewExchange(fourth, &wg)
		go h.srv.AsyncGetNewExchange(fifth, &wg)

		wg.Wait()

		addDays = 5
		time.Sleep(time.Second * 2)
		startFromDate = startFromDate.AddDate(0, 0, 5)
	}
}
