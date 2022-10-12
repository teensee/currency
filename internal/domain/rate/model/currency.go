package model

import (
	"gorm.io/datatypes"
	"time"
)

type CurrencyRate struct {
	//ID           uuid.UUID
	CurrencyFrom string         `gorm:"primaryKey:idx_name" json:"from"`
	CurrencyTo   string         `gorm:"primaryKey:idx_name" json:"to"`
	OnDate       datatypes.Date `gorm:"primaryKey:idx_name" json:"onDate"`
	ExchangeRate float64        `json:"rate"`
	CreatedAt    time.Time      `json:"-"` //ignore this field
}

type RatePair struct {
	Rate        CurrencyRate
	ReverseRate CurrencyRate
}

type RatePairCollection struct {
	RateList        []CurrencyRate
	ReverseRateList []CurrencyRate
}

func (rpc *RatePairCollection) AddRate(rate CurrencyRate) []CurrencyRate {
	rpc.RateList = append(rpc.RateList, rate)

	return rpc.RateList
}

func (rpc *RatePairCollection) AddReverseRate(rate CurrencyRate) []CurrencyRate {
	rpc.ReverseRateList = append(rpc.ReverseRateList, rate)

	return rpc.ReverseRateList
}

func (rpc *RatePairCollection) AddRatePair(rate RatePair) ([]CurrencyRate, []CurrencyRate) {
	rpc.AddRate(rate.Rate)
	rpc.AddReverseRate(rate.ReverseRate)

	return rpc.RateList, rpc.ReverseRateList
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
