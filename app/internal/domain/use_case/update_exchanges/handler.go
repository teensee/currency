package update_exchanges

import (
	"Currency/internal/config"
	"Currency/internal/domain/model"
	"Currency/internal/domain/service"
	"log"
	"net/http"
	"net/url"
	"time"
)

type UpdateExchangeHandler struct {
	srv    *service.ExchangeRateService
	client *service.CbrClient
}

func NewHandler(srv *service.ExchangeRateService) *UpdateExchangeHandler {
	return &UpdateExchangeHandler{
		srv:    srv,
		client: service.NewCbrClient(),
	}
}

// GetCbrExchangeRates Creates a CBR rates in database
func (h UpdateExchangeHandler) GetCbrExchangeRates(req *http.Request) {
	onDate := extractFilters(req.URL.Query())
	isExist := h.srv.IsExistOnDate(onDate)

	if isExist == true {
		return
	}

	cbrRates := h.client.GetCbrRates(onDate)
	log.Printf("Cbr Rates on: %s successfully parsed from cbr", cbrRates.Date)

	var ratePairCollection model.RatePairCollection
	for _, exchangeRate := range cbrRates.Valute {
		pair := h.srv.CreateRatePair(exchangeRate, onDate)
		ratePairCollection.AddRatePair(pair)
	}

	log.Printf("Pairs was inserted to db")

	h.srv.SaveRatePairCollection(ratePairCollection)
	go func() { h.srv.TriangulateRates(onDate) }()
}

func extractFilters(query url.Values) time.Time {
	filters, present := query["onDate"]

	if !present && len(filters) == 0 {
		onDate := time.Now()
		log.Printf("Time filter not presen, used now: %s", onDate)

		return onDate
	} else {
		onDate, err := time.Parse(config.ApiDateFormat, filters[0])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Time filter passed, %s", onDate)

		return onDate
	}
}
