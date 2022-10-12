package service

import (
	"Currency/internal/infrastructure/model"
	rateRepository "Currency/internal/infrastructure/repository/rate"
	exchangeDto "Currency/internal/infrastructure/use_case/get_exchanges/dto"
	"Currency/internal/infrastructure/use_case/update_exchanges/dto"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

type ExchangeRateService struct {
	repo   *rateRepository.RateRepository
	client *CbrClient
}

func NewExchangeRateService(db *gorm.DB) *ExchangeRateService {
	return &ExchangeRateService{
		repo: rateRepository.NewRateRepository(db),
	}
}

// GetExchangeRate Возвращает курс валют для пары currencyFrom/currencyTo на конкретную дату
func (s *ExchangeRateService) GetExchangeRate(currencyFrom, currencyTo string, onDate time.Time) model.CurrencyRate {
	return s.repo.ExchangeRate(currencyFrom, currencyTo, onDate)
}

// ConvertCurrency Конвертация валюты для пары currencyFrom/currencyTo на конкретную дату
func (s *ExchangeRateService) ConvertCurrency(from string, to string, onDate time.Time, value float64) exchangeDto.ConvertResult {
	exchangeRate := s.GetExchangeRate(from, to, onDate)

	return exchangeDto.NewConvertResult(
		exchangeRate.CurrencyFrom,
		exchangeRate.CurrencyTo,
		value,
		exchangeRate.ExchangeRate,
		exchangeRate.ExchangeRate*value,
		onDate,
	)
}

// IsExistOnDate Проверка существует ли на заданную дату курс валют
func (s *ExchangeRateService) IsExistOnDate(date time.Time) bool {
	return s.repo.IsExistOnDate(date)
}

func (s *ExchangeRateService) GetNewExchangeRateOnDate(onDate time.Time) {
	isExist := s.IsExistOnDate(onDate)

	if isExist == true {
		return
	}
	cbrRates := s.client.GetCbrRates(onDate)
	log.Printf("Cbr Rates on: %s date was successfully parsed", cbrRates.Date)

	var ratePairCollection model.RatePairCollection
	for _, exchangeRate := range cbrRates.Valute {
		pair := s.createRatePair(exchangeRate, onDate)
		ratePairCollection.AddRatePair(pair)
	}

	s.saveRatePairCollection(ratePairCollection)
	go func() { s.triangulateRates(onDate) }()

	log.Printf("Pairs was inserted to db")
}

func (s *ExchangeRateService) createRatePair(rate dto.CbrRate, date time.Time) model.RatePair {
	log.Printf("process %s/RUB", rate.CharCode)

	exchangeRate, _ := strconv.ParseFloat(strings.Replace(rate.Value, ",", ".", 1), 64)
	nominal, _ := strconv.ParseFloat(strings.Replace(rate.Nominal, ",", ".", 1), 64)

	return model.NewRatePair(rate.CharCode, "RUB", exchangeRate/nominal, date)
}

func (s *ExchangeRateService) saveRatePairCollection(pairs model.RatePairCollection) {
	s.repo.SavePairCollection(pairs)
}

func (s *ExchangeRateService) triangulateRates(onDate time.Time) {
	s.repo.TriangulateRates(onDate)
}
