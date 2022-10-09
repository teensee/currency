package dto

import (
	"encoding/xml"
)

type CbrRates struct {
	XMLName xml.Name  `xml:"ValCurs"`
	Date    string    `xml:"Date,attr"`
	Name    string    `xml:"name,attr"`
	Valute  []CbrRate `xml:"Valute"`
}

type CbrRate struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}
