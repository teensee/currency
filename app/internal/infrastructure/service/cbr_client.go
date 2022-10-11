package service

import (
	"Currency/internal/config"
	"Currency/internal/infrastructure/use_case/update_exchanges/dto"
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"time"
)

type CbrClient struct {
	BasePath string
}

func NewCbrClient() *CbrClient {
	return &CbrClient{}
}

// GetCbrRates Client request to CBR which return currency exchange rate list
func (c *CbrClient) GetCbrRates(onDate time.Time) dto.CbrRates {
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
