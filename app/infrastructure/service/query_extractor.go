package service

import (
	"Currency/infrastructure/dto"
	"Currency/internal/config"
	"errors"
	"log"
	"net/url"
	"strconv"
	"time"
)

type QueryExtractor struct {
}

func NewQueryExtractor() *QueryExtractor {
	return &QueryExtractor{}
}

func (qe *QueryExtractor) ExtractDateParam(query url.Values) dto.Date {
	onDateFilter, present := query["onDate"]

	if !present && len(onDateFilter) == 0 {
		now := time.Now()
		y, m, d := now.Date()
		nowFtm := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		log.Printf("Time filter not presen, used now: %s", nowFtm)

		return dto.NewOnDate(nowFtm)
	} else {
		onDate, err := time.Parse(config.ApiDateFormat, onDateFilter[0])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Time filter passed, %s", onDate)

		return dto.NewOnDate(onDate)
	}
}

// ExtractRateParams
// Базовая валидация и проверка на наличие get-параметров 'from', 'to', 'onDate'
// Для запроса вида /exchange?from=USD&to=EUR&onDate=20.12.2000
func (qe *QueryExtractor) ExtractRateParams(query url.Values) (dto.ExchangeRequestParameters, error) {
	fromParam, fromPresent := query["from"]
	toParam, toPresent := query["to"]
	onDateParam, onDatePresent := query["onDate"]

	if !fromPresent && len(fromParam) == 0 {
		return dto.ExchangeRequestParameters{}, errors.New("'from' is required query parameter")
	}

	if !toPresent && len(toParam) == 0 {
		return dto.ExchangeRequestParameters{}, errors.New("'to' is required query parameter")
	}

	if !onDatePresent && len(onDateParam) == 0 {
		return dto.ExchangeRequestParameters{}, errors.New("'onDate' is required query parameter")
	}

	onDate, err := time.Parse(config.ApiDateFormat, onDateParam[0])

	if err != nil {
		return dto.ExchangeRequestParameters{}, errors.New("'onDate' is not correctly passed, please use a DD.MM.YYYY format")
	}

	return dto.NewExchangeRequestParameters(fromParam[0], toParam[0], onDate), nil
}

func (qe *QueryExtractor) ExtractConvertParams(query url.Values) (dto.ConvertRequestParameters, error) {
	exchangeParams, err := qe.ExtractRateParams(query)
	if err != nil {
		return dto.ConvertRequestParameters{}, err
	}

	valueParam, valuePresent := query["value"]

	if !valuePresent && len(valueParam) == 0 {
		return dto.ConvertRequestParameters{}, errors.New("'value' is required query parameter")
	}

	value, err := strconv.ParseFloat(valueParam[0], 64)

	if value < 0 {
		return dto.ConvertRequestParameters{}, errors.New("'value' must be greater than 0")
	}

	if err != nil {
		return dto.ConvertRequestParameters{}, err
	}

	return dto.NewConvertRequestParametersWithExchange(exchangeParams, value), nil
}
