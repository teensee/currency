package update_exchanges

import (
	"Currency/internal/config"
	"Currency/internal/domain/service"
	"Currency/internal/domain/use_case/update_exchanges/dto"
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type UpdateExchangeHandler struct {
	srv *service.ExchangeRateService
}

func NewHandler(srv *service.ExchangeRateService) *UpdateExchangeHandler {
	return &UpdateExchangeHandler{
		srv: srv,
	}
}

// ExchangeCbrRates Creates a CBR rates in database
func (h UpdateExchangeHandler) ExchangeCbrRates(req *http.Request) {
	onDate := extractFilters(req.URL.Query())
	cbrRates := getCbrRates(onDate)
	log.Printf("Cbr Rates on: %s successfully parsed from cbr", cbrRates.Date)

	for _, exchangeRate := range cbrRates.Valute {
		ratePair := processRate(exchangeRate, onDate)
		h.srv.GetRateRepository().SavePair(ratePair)

		log.Print(ratePair)
	}

	log.Printf("Pairs was inserted to db")

	h.srv.GetRateRepository().TriangulateRates(onDate)
}

func extractFilters(query url.Values) time.Time {
	filters, present := query["onDate"]

	if !present && len(filters) == 0 {
		onDate := time.Now()
		log.Printf("Time filter not presen, used now: %s", onDate)

		return onDate
	} else {
		onDate, err := time.Parse(config.ApiDateFormat, filters[0])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Time filter passed, %s", onDate)

		return onDate
	}
}

func processRate(cbrRate dto.CbrRate, onDate time.Time) RatePair {
	log.Printf("process %s/RUB", cbrRate.CharCode)

	exchangeRate, _ := strconv.ParseFloat(strings.Replace(cbrRate.Value, ",", ".", 1), 64)
	nominal, _ := strconv.ParseFloat(strings.Replace(cbrRate.Nominal, ",", ".", 1), 64)

	return NewRatePair(cbrRate.CharCode, "RUB", exchangeRate/nominal, onDate)
}

// Client request to CBR which return currency exchange rate list
func getCbrRates(onDate time.Time) dto.CbrRates {
	requestUrl := fmt.Sprintf("https://cbr.ru/scripts/XML_daily.asp?date_req=%s", onDate.Format(config.CbrDateFormat))
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
