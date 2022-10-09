package service

import (
	rateRepository "Currency/internal/domain/repository/rate"
	"gorm.io/gorm"
)

type ExchangeRateService struct {
	repo *rateRepository.RateRepository
}

func NewExchangeRateService(db *gorm.DB) *ExchangeRateService {
	return &ExchangeRateService{
		repo: rateRepository.NewRateRepository(db),
	}
}

func (s ExchangeRateService) GetRateRepository() *rateRepository.RateRepository {
	return s.repo
}
