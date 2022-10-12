package rate_repository

import (
	"Currency/internal/config"
	"Currency/internal/infrastructure/model"
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type RateRepository struct {
	db *gorm.DB
}

func NewRateRepository(db *gorm.DB) *RateRepository {
	return &RateRepository{
		db: db,
	}
}

// Save Сохранение курса model.CurrencyRate
func (r *RateRepository) Save(rate model.CurrencyRate) {
	r.db.Save(rate)
}

// SavePair Сохраннеие пары model.RatePair
// Rate - прямой курс валют полученный от цб
// ReverseRate - обратный, полученный путем деления 1 / Rate
func (r *RateRepository) SavePair(rate model.RatePair) {
	r.db.Save(rate.Rate)
	r.db.Save(rate.ReverseRate)
}

// SavePairCollection Сохранение коллекций
// Rate 	   []model.CurrencyRate и
// ReverseRate []model.CurrencyRate
func (r *RateRepository) SavePairCollection(rate model.RatePairCollection) {
	r.db.Save(rate.RateList)
	r.db.Save(rate.ReverseRateList)
}

func (r *RateRepository) TriangulateRates(onDate time.Time) {
	r.db.Exec("insert into currency_rates (currency_from, currency_to, on_date, exchange_rate, created_at) "+
		"select pair.currency_from, pair.currency_to, pair.on_date, pair.exchange_rate, now()"+
		"from ("+
		"select f.currency_from, t.currency_from as currency_to, f.created_at, f.on_date, (f.exchange_rate / t.exchange_rate) as exchange_rate from currency_rates f, currency_rates t "+
		"where f.on_date = t.on_date and f.currency_to = t.currency_to"+
		") pair "+
		"LEFT OUTER JOIN currency_rates cr ON (pair.on_date = cr.on_date AND pair.currency_from = cr.currency_from AND pair.currency_to = cr.currency_to) "+
		"where pair.on_date = @onDate "+
		"group by pair.currency_from, pair.currency_to, pair.on_date, pair.exchange_rate", sql.Named("onDate", onDate.Format(config.DbDateFormat)))
}

func (r *RateRepository) ExchangeRate(currencyFrom, currencyTo string, onDate time.Time) model.CurrencyRate {
	var rate model.CurrencyRate
	r.db.
		Where(
			"currency_from = ? and currency_to = ? and on_date = ?",
			currencyFrom,
			currencyTo,
			onDate.Format(config.DbDateFormat),
		).
		First(&rate)

	return rate
}

func (r *RateRepository) IsExistOnDate(date time.Time) bool {
	var count int64
	r.db.Model(&model.CurrencyRate{}).Where("on_date = ?", date.Format(config.DbDateFormat)).Count(&count)

	return count > 0
}
