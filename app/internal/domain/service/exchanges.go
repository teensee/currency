package service

import (
	"Currency/internal/domain/model"
	rateRepository "Currency/internal/domain/repository/rate"
	"Currency/internal/domain/use_case/update_exchanges/dto"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

type ExchangeRateService struct {
	repo *rateRepository.RateRepository
}

func NewExchangeRateService(db *gorm.DB) *ExchangeRateService {
	return &ExchangeRateService{
		repo: rateRepository.NewRateRepository(db),
	}
}

func (s ExchangeRateService) GetExchangeRate(currencyFrom, currencyTo string, onDate time.Time) model.CurrencyRate {
	return s.repo.ExchangeRate(currencyFrom, currencyTo, onDate)
}

func (s ExchangeRateService) CreateRatePair(rate dto.CbrRate, date time.Time) model.RatePair {
	log.Printf("process %s/RUB", rate.CharCode)

	exchangeRate, _ := strconv.ParseFloat(strings.Replace(rate.Value, ",", ".", 1), 64)
	nominal, _ := strconv.ParseFloat(strings.Replace(rate.Nominal, ",", ".", 1), 64)

	return model.NewRatePair(rate.CharCode, "RUB", exchangeRate/nominal, date)
}

func (s ExchangeRateService) SaveRatePairCollection(pairs model.RatePairCollection) {
	s.repo.SavePairCollection(pairs)
}

func (s ExchangeRateService) TriangulateRates(onDate time.Time) {
	s.repo.TriangulateRates(onDate)
}
