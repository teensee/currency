package storage

import (
	"Currency/internal/domain/rate/model"
	"time"
)

type ExchangeRateRepository interface {
	Save(rate model.CurrencyRate)
	SavePair(rate model.RatePair)
	SavePairCollection(rate model.RatePairCollection)
	TriangulateRates(onDate time.Time)
	ExchangeRate(currencyFrom, currencyTo string, onDate time.Time) model.CurrencyRate
	IsExistOnDate(date time.Time) bool
}
