package get_exchanges

import (
	"Currency/internal/domain/use_case/update_exchanges"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type GetExchangeRateHandler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *GetExchangeRateHandler {
	return &GetExchangeRateHandler{
		db: db,
	}
}

func (h GetExchangeRateHandler) ExchangeRate(r *http.Request) update_exchanges.CurrencyRate {
	currencyFrom, currencyTo, onDate := extractRateFilters(r.URL.Query())

	var rate update_exchanges.CurrencyRate
	h.db.Where("currency_from = ? and currency_to = ? and on_date = ?", currencyFrom, currencyTo, onDate.Format("2006-01-02 00:00:00+00:00")).First(&rate)

	return rate
}

func extractRateFilters(query url.Values) (string, string, time.Time) {
	fromParam, fromPresent := query["from"]
	toParam, toPresent := query["to"]
	onDateParam, onDatePresent := query["onDate"]

	if !fromPresent {
		log.Fatal("'From' is required query parameter")
	}

	if !toPresent {
		log.Fatal("'To' is required query parameter")
	}

	if !onDatePresent {
		log.Fatal("'OnDate' is required query parameter")
	}

	onDate, _ := time.Parse("02.01.2006", onDateParam[0])

	return fromParam[0], toParam[0], onDate
}

func extractFilters(query url.Values) (string, string, float64, time.Time) {
	fromParam, fromPresent := query["from"]
	toParam, toPresent := query["to"]
	onDateParam, onDatePresent := query["onDate"]
	valueParam, valuePresent := query["value"]

	if !fromPresent {
		log.Fatal("'From' is required query parameter")
	}

	if !toPresent {
		log.Fatal("'To' is required query parameter")
	}

	if !valuePresent {
		log.Fatal("'Value' is required query parameter")
	}
	if !onDatePresent {
		log.Fatal("'OnDate' is required query parameter")
	}

	value, _ := strconv.ParseFloat(valueParam[0], 64)
	onDate, _ := time.Parse("02.01.2006", onDateParam[0])

	return fromParam[0], toParam[0], value, onDate
}
