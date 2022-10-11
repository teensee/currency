package update_exchanges

import (
	service2 "Currency/infrastructure/service"
	"log"
	"net/http"
	"time"
)

type UpdateExchangeHandler struct {
	srv            *service2.ExchangeRateService
	queryExtractor *service2.QueryExtractor
}

func NewHandler(srv *service2.ExchangeRateService) *UpdateExchangeHandler {
	return &UpdateExchangeHandler{
		srv:            srv,
		queryExtractor: service2.NewQueryExtractor(),
	}
}

// GetCbrExchangeRates Creates a CBR rates in database
func (h UpdateExchangeHandler) GetCbrExchangeRates(_ http.ResponseWriter, r *http.Request) {
	onDate := h.queryExtractor.ExtractDateParam(r.URL.Query())
	h.srv.GetNewExchangeRateOnDate(onDate.OnDate)
}

func (h UpdateExchangeHandler) SyncRatesOnStartup() {
	now := time.Now()
	y, m, d := now.Date()
	nowFtm := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	log.Printf("Syncing rates on %s", nowFtm.Format("02.01.2006"))

	h.srv.GetNewExchangeRateOnDate(nowFtm)
}
