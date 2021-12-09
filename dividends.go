package gomoex

import (
	"context"
	"time"

	"github.com/tidwall/gjson"
)

const (
	_dividendSECID     = `secid`
	_dividendISIN      = `isin`
	_dividendCloseDate = `registryclosedate`
	_dividendValue     = `value`
	_dividendCurrency  = `currencyid`
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

	div.Ticker = row.Get(_dividendSECID).String()
	div.ISIN = row.Get(_dividendISIN).String()

	div.Date, err = time.Parse("2006-01-02", row.Get(_dividendCloseDate).String())
	if err != nil {
		return nil, newParseErr(err)
	}

	div.Dividend = row.Get(_dividendValue).Float()
	div.Currency = row.Get(_dividendCurrency).String()

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

	for raw := range iss.rowsGen(ctx, query.Make()) {
		switch row := raw.(type) {
		case Dividend:
			table = append(table, row)
		case error:
			return nil, row
		}
	}

	return table, nil
}
