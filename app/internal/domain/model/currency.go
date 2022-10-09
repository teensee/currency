package model

import (
	"gorm.io/datatypes"
	"time"
)

type CurrencyRate struct {
	//ID           uuid.UUID
	CurrencyFrom string         `gorm:"primaryKey:idx_name"`
	CurrencyTo   string         `gorm:"primaryKey:idx_name"`
	OnDate       datatypes.Date `gorm:"primaryKey:idx_name"`
	CreatedAt    time.Time
	ExchangeRate float64
}

type RatePair struct {
	Rate        CurrencyRate
	ReverseRate CurrencyRate
}

// NewRate return a new exchange rate
func NewRate(currencyFrom, currencyTo string, rate float64, onDate time.Time) CurrencyRate {
	return CurrencyRate{
		//ID:           uuid.New(),
		CurrencyFrom: currencyFrom,
		CurrencyTo:   currencyTo,
		CreatedAt:    time.Now(),
		OnDate:       datatypes.Date(onDate),
		ExchangeRate: rate,
	}
}

// NewRatePair return rate and reverse rate
func NewRatePair(currencyFrom, currencyTo string, rate float64, onDate time.Time) RatePair {
	return RatePair{
		Rate:        NewRate(currencyFrom, currencyTo, rate, onDate),
		ReverseRate: NewRate(currencyTo, currencyFrom, 1/rate, onDate),
	}
}
