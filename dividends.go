package gomoex

import (
	"context"
	"time"

	"github.com/tidwall/gjson"
)

// Dividend содержит информацию дате закрытия реестра, дивиденде и валюте.
type Dividend struct {
	Ticker   string
	ISIN     string
	Date     time.Time
	Dividend float64
	Currency string
}

func dividendConverter(row gjson.Result) (interface{}, error) {
	var (
		err error
		div Dividend
	)

	div.Ticker = row.Get("secid").String()
	div.ISIN = row.Get("isin").String()

	div.Date, err = time.Parse("2006-01-02", row.Get("registryclosedate").String())
	if err != nil {
		return nil, newParseErr(err)
	}

	div.Dividend = row.Get("value").Float()
	div.Currency = row.Get("currencyid").String()

	return div, nil
}

// Dividends получает таблицу с дивидендами.
//
// Запрос не отражен в официальном справочнике. По многим инструментам дивиденды отсутствуют или отражены не полностью.
// Корректная информация содержится в основном только по наиболее ликвидным бумагам.
func (iss *ISSClient) Dividends(ctx context.Context, security string) (table []Dividend, err error) {
	query := querySettings{
		security:     security,
		object:       "dividends",
		table:        "dividends",
		rowConverter: dividendConverter,
	}

	for raw := range iss.getRowsGen(ctx, query.Make()) {
		switch row := raw.(type) {
		case Dividend:
			table = append(table, row)
		case error:
			return nil, row
		}
	}

	return table, nil
}
