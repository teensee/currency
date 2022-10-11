package dto

import (
	"strings"
	"time"
)

type ExchangeRequestParameters struct {
	From   string
	To     string
	OnDate time.Time
}

type ConvertRequestParameters struct {
	ExchangeRequestParameters
	Value float64
}

func NewExchangeRequestParameters(from, to string, onDate time.Time) ExchangeRequestParameters {
	return ExchangeRequestParameters{
		From:   strings.ToUpper(from),
		To:     strings.ToUpper(to),
		OnDate: onDate,
	}
}

func NewConvertRequestParametersWithExchange(exchange ExchangeRequestParameters, value float64) ConvertRequestParameters {
	return ConvertRequestParameters{
		ExchangeRequestParameters: exchange,
		Value:                     value,
	}
}
