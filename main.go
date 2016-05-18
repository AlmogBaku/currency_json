package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"io/ioutil"
)

const currencies_data_source = "http://www.currency-iso.org/dam/downloads/lists/list_one.xml"

type Currency struct {
	Code         string `xml:"Ccy",json:"iso_code"`
	Num          string `xml:"CcyNbr",json:"iso_num"`
	MinorUnit    string `xml:"CcyMnrUnts",json:"minor_unit,omitempty"`
	CountryName  string `xml:"CtryNm",json:"country_name,omitempty"`
	CurrencyName string `xml:"CcyNm",json:"currency_name"`
}

type ISO_4217 struct {
	Published  string     `xml:"Pblshd,attr",json:"published"`
	Currencies []Currency `xml:"CcyTbl>CcyNtry",json:"currencies"`
}

func main() {
	resp, err := http.Get(currencies_data_source)
	if err != nil {
		panic(err)
	}

	var list_one ISO_4217
	err = xml.NewDecoder(resp.Body).Decode(&list_one)
	if err != nil {
		panic(err)
	}

	err = obj2jsonfile("list_one.json", list_one)
	if err != nil {
		panic(err)
	}

	currenciesByCode := make(map[string]Currency)
	for _, c := range list_one.Currencies {
		if c.Code != "" {
			curr := *&c
			curr.CountryName = ""
			currenciesByCode[c.Code] = curr
		}
	}
	err = obj2jsonfile("currencies.json", currenciesByCode)
	if err != nil {
		panic(err)
	}
}

func obj2jsonfile(filename string, obj interface{}) error {
	jsonized, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, jsonized, 0644)
	if err != nil {
		return err
	}

	return nil
}
