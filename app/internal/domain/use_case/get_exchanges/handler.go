package get_exchanges

import (
	"Currency/internal/config"
	"Currency/internal/domain/model"
	"Currency/internal/domain/service"
	"Currency/internal/domain/use_case/get_exchanges/dto"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type GetExchangeRateHandler struct {
	service *service.ExchangeRateService
}

func NewHandler(srv *service.ExchangeRateService) *GetExchangeRateHandler {
	return &GetExchangeRateHandler{
		service: srv,
	}
}

func (h GetExchangeRateHandler) ExchangeRate(r *http.Request) model.CurrencyRate {
	currencyFrom, currencyTo, onDate := extractRateFilters(r.URL.Query())

	return h.service.GetExchangeRate(currencyFrom, currencyTo, onDate)
}

func (h GetExchangeRateHandler) Convert(r *http.Request) dto.ConvertResult {
	currencyFrom, currencyTo, value, onDateFilter := extractConvertFilters(r.URL.Query())
	exchangeRate := h.service.GetExchangeRate(currencyFrom, currencyTo, onDateFilter)

	return dto.NewConvertResult(
		exchangeRate.CurrencyFrom,
		exchangeRate.CurrencyTo,
		value,
		exchangeRate.ExchangeRate,
		exchangeRate.ExchangeRate*value,
		onDateFilter,
	)
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

	onDate, _ := time.Parse(config.ApiDateFormat, onDateParam[0])

	return fromParam[0], toParam[0], onDate
}

func extractConvertFilters(query url.Values) (string, string, float64, time.Time) {
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
	onDate, _ := time.Parse(config.ApiDateFormat, onDateParam[0])

	return fromParam[0], toParam[0], value, onDate
}
