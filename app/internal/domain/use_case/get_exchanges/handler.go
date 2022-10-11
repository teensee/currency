package get_exchanges

import (
	"Currency/internal/config"
	"Currency/internal/domain/model"
	"Currency/internal/domain/service"
	"Currency/internal/domain/use_case/get_exchanges/dto"
	"errors"
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

func (h GetExchangeRateHandler) ExchangeRate(r *http.Request) (model.CurrencyRate, error) {
	reqParam, err := extractRateFilters(r.URL.Query())

	if err != nil {
		return model.CurrencyRate{}, err
	}

	return h.service.GetExchangeRate(reqParam.From, reqParam.To, reqParam.OnDate), nil
}

func (h GetExchangeRateHandler) Convert(r *http.Request) (dto.ConvertResult, error) {
	reqParam, err := extractConvertFilters(r.URL.Query())

	if err != nil {
		return dto.ConvertResult{}, err
	}
	exchangeRate := h.service.GetExchangeRate(reqParam.From, reqParam.To, reqParam.OnDate)

	return dto.NewConvertResult(
		exchangeRate.CurrencyFrom,
		exchangeRate.CurrencyTo,
		reqParam.Value,
		exchangeRate.ExchangeRate,
		exchangeRate.ExchangeRate*reqParam.Value,
		reqParam.OnDate,
	), nil
}

// Базовая валидация и проверка на наличие get-параметров 'from', 'to', 'onDate'
// Для запроса вида /exchange?from=USD&to=EUR&onDate=20.12.2000
func extractRateFilters(query url.Values) (dto.ExchangeRequestParameters, error) {
	fromParam, fromPresent := query["from"]
	toParam, toPresent := query["to"]
	onDateParam, onDatePresent := query["onDate"]

	if !fromPresent {
		return dto.ExchangeRequestParameters{}, errors.New("'from' is required query parameter")
	}

	if !toPresent {
		return dto.ExchangeRequestParameters{}, errors.New("'to' is required query parameter")
	}

	if !onDatePresent {
		return dto.ExchangeRequestParameters{}, errors.New("'onDate' is required query parameter")
	}

	onDate, err := time.Parse(config.ApiDateFormat, onDateParam[0])

	if err != nil {
		return dto.ExchangeRequestParameters{}, errors.New("'onDate' is not correctly passed, please use a DD.MM.YYYY format")
	}

	return dto.NewExchangeRequestParameters(fromParam[0], toParam[0], onDate), nil
}

func extractConvertFilters(query url.Values) (dto.ConvertRequestParameters, error) {
	exchangeParams, err := extractRateFilters(query)
	if err != nil {
		return dto.ConvertRequestParameters{}, err
	}

	valueParam, valuePresent := query["value"]

	if !valuePresent {
		return dto.ConvertRequestParameters{}, errors.New("'value' is required query parameter")
	}

	value, err := strconv.ParseFloat(valueParam[0], 64)

	if err != nil {
		return dto.ConvertRequestParameters{}, err
	}

	return dto.NewConvertRequestParametersWithExchange(exchangeParams, value), nil
}
