package dto

import (
	"Currency/internal/config"
	"fmt"
	"math"
	"time"
)

type ConvertResult struct {
	Pair          string  `json:"pair"`
	Value         float64 `json:"value"`
	Rate          float64 `json:"rate"`
	RawResult     float64 `json:"rawResult"`
	RoundedResult float64 `json:"roundedResult"`
	OnDate        string  `json:"onDate"`
}

func NewConvertResult(from, to string, value, rate, convertedResult float64, onDate time.Time) ConvertResult {
	return ConvertResult{
		Pair:          fmt.Sprintf("%s/%s", from, to),
		Value:         value,
		Rate:          rate,
		RawResult:     convertedResult,
		RoundedResult: math.Round(convertedResult*100) / 100,
		OnDate:        onDate.Format(config.ApiDateFormat),
	}
}
