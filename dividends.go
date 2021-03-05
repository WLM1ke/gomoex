package gomoex

import (
	"context"
	"github.com/francoispqt/gojay"
	"time"
)

// Dividend содержит информацию дате закрытия реестра, дивиденде и валюте.
type Dividend struct {
	Ticker   string
	ISIN     string
	Date     time.Time
	Dividend float64
	Currency string
}

func (dividend *Dividend) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "secid":
		return dec.String(&dividend.Ticker)
	case "isin":
		return dec.String(&dividend.ISIN)
	case "registryclosedate":
		return dec.Time(&dividend.Date, "2006-01-02")
	case "value":
		return dec.Float(&dividend.Dividend)
	case "currencyid":
		return dec.String(&dividend.Currency)
	}
	return nil
}
func (dividend *Dividend) NKeys() int {
	return 5
}

// SecurityDividends получает таблицу с дивидендами.
// Запрос не отражен в официальном справочнике. По многим инструментам дивиденды отсутствуют или отражены не полностью.
// Корректная информация содержится в основном только по наиболее ликвидным позициям.
func (iss ISSClient) SecurityDividends(ctx context.Context, security string) (table []Dividend, err error) {
	query := ISSQuery{
		security: security,
		object:   "dividends",
		table:    "dividends",
	}

	rows, errors := iss.getAll(ctx, query)

	for rawRow := range rows {
		table = append(table, Dividend{})
		err = gojay.Unmarshal(rawRow, &table[len(table)-1])
		if err != nil {
			return nil, err
		}
	}
	if err = <-errors; err != nil {

		return nil, err
	}
	return table, nil
}
