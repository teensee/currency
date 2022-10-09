package rate_repository

import (
	"Currency/internal/config"
	"Currency/internal/domain/use_case/update_exchanges"
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

func (r RateRepository) Save(rate update_exchanges.CurrencyRate) {
	r.db.Save(rate)
}

func (r RateRepository) SavePair(rate update_exchanges.RatePair) {
	r.db.Save(rate.Rate)
	r.db.Save(rate.ReverseRate)
}

func (r RateRepository) TriangulateRates(onDate time.Time) {
	r.db.Exec("insert into currency_rates (currency_from, currency_to, created_at, on_date, exchange_rate) "+
		"select pair.currency_from, pair.currency_to, pair.created_at, pair.on_date, pair.exchange_rate "+
		"from ("+
		"select f.currency_from, t.currency_from as currency_to, f.created_at, f.on_date, (f.exchange_rate / t.exchange_rate) as exchange_rate from currency_rates f, currency_rates t "+
		"where f.on_date = t.on_date and f.currency_to = t.currency_to"+
		") pair "+
		"LEFT OUTER JOIN currency_rates cr ON (pair.on_date = cr.on_date AND pair.currency_from = cr.currency_from AND pair.currency_to = cr.currency_to) "+
		"where pair.on_date = @onDate "+
		"group by pair.currency_from, pair.currency_to", sql.Named("onDate", onDate.Format(config.DbDateFormat)))
}

func (r RateRepository) ExchangeRate(currencyFrom, currencyTo string, onDate time.Time) update_exchanges.CurrencyRate {
	var rate update_exchanges.CurrencyRate
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