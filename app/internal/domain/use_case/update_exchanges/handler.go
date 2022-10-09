package update_exchanges

import (
	"Currency/internal/domain/use_case/update_exchanges/dto"
	"database/sql"
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type UpdateExchangeHandler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *UpdateExchangeHandler {
	return &UpdateExchangeHandler{
		db: db,
	}
}

// ExchangeCbrRates Creates a CBR rates in database
func (e UpdateExchangeHandler) ExchangeCbrRates(req *http.Request) {
	onDate := extractFilters(req.URL.Query())

	rates := getCbrRates(onDate)
	log.Printf("Cbr Rates on: %s successfully parsed from cbr", rates.Date)

	for _, rate := range rates.Valute {
		ratePair := processRate(rate, onDate)

		e.db.Save(ratePair.Rate)
		e.db.Save(ratePair.ReverseRate)

		log.Print(ratePair)
	}

	log.Printf("Pairs was inserted to db")

	//insert into currency_rates (currency_from, currency_to, created_at, on_date, exchange_rate)
	//select pair.currency_from, pair.currency_to, pair.created_at, pair.on_date, pair.exchange_rate
	//from (
	//         select f.currency_from, t.currency_from as currency_to, f.created_at, f.on_date, (f.exchange_rate / t.exchange_rate) as exchange_rate
	//         from currency_rates f, currency_rates t
	//         where f.on_date = t.on_date
	//           and f.currency_to = t.currency_to
	//     ) pair
	//         LEFT OUTER JOIN currency_rates cr
	//                         ON (
	//                                     pair.on_date = cr.on_date
	//                                 AND pair.currency_from = cr.currency_from
	//                                 AND pair.currency_to = cr.currency_to
	//                             )
	//group by pair.currency_from, pair.currency_to
	e.db.Exec("insert into currency_rates (currency_from, currency_to, created_at, on_date, exchange_rate) "+
		"select pair.currency_from, pair.currency_to, pair.created_at, pair.on_date, pair.exchange_rate "+
		"from ("+
		"select f.currency_from, t.currency_from as currency_to, f.created_at, f.on_date, (f.exchange_rate / t.exchange_rate) as exchange_rate from currency_rates f, currency_rates t "+
		"where f.on_date = t.on_date and f.currency_to = t.currency_to"+
		") pair "+
		"LEFT OUTER JOIN currency_rates cr ON (pair.on_date = cr.on_date AND pair.currency_from = cr.currency_from AND pair.currency_to = cr.currency_to) "+
		"where pair.on_date = @onDate "+
		"group by pair.currency_from, pair.currency_to", sql.Named("onDate", onDate.Format("2006-01-02 00:00:00+00:00")))
}

func extractFilters(query url.Values) time.Time {
	filters, present := query["onDate"]

	if !present && len(filters) == 0 {
		onDate := time.Now()
		log.Printf("Time filter not presen, used now: %s", onDate)

		return onDate
	} else {
		onDate, err := time.Parse("02.01.2006", filters[0])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Time filter passed, %s", onDate)

		return onDate
	}
}

func processRate(cbrRate dto.CbrRate, onDate time.Time) RatePair {
	log.Printf("process %s/RUB", cbrRate.CharCode)
	rate, _ := strconv.ParseFloat(strings.Replace(cbrRate.Value, ",", ".", 1), 64)
	nominal, _ := strconv.ParseFloat(strings.Replace(cbrRate.Nominal, ",", ".", 1), 64)

	return NewRatePair(cbrRate.CharCode, "RUB", rate/nominal, onDate)
}

// Client request to CBR which return currency exchange rate list
func getCbrRates(onDate time.Time) dto.CbrRates {
	requestUrl := "https://cbr.ru/scripts/XML_daily.asp?date_req=" + onDate.Format("02/01/2006")
	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	d := xml.NewDecoder(resp.Body)
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}

	var rates dto.CbrRates
	err = d.Decode(&rates)

	if err != nil {
		log.Fatal(err)
	}

	return rates
}
