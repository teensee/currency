package dto

import (
	"strings"
	"time"
)

// Date Дата на которую запрашивают курс валют
type Date struct {
	OnDate time.Time
}

// ExchangeRequestParameters Запрос обменного курса на дату
type ExchangeRequestParameters struct {
	From   string
	To     string
	OnDate Date
}

type ConvertRequestParameters struct {
	ExchangeRequestParameters
	Value float64
}

func NewExchangeRequestParameters(from, to string, onDate time.Time) ExchangeRequestParameters {
	return ExchangeRequestParameters{
		From:   strings.ToUpper(from),
		To:     strings.ToUpper(to),
		OnDate: NewOnDate(onDate),
	}
}

func NewConvertRequestParametersWithExchange(exchange ExchangeRequestParameters, value float64) ConvertRequestParameters {
	return ConvertRequestParameters{
		ExchangeRequestParameters: exchange,
		Value:                     value,
	}
}

func NewOnDate(onDate time.Time) Date {
	return Date{
		OnDate: onDate,
	}
}
