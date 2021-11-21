package gomoex

import (
	"context"
	"time"

	"github.com/valyala/fastjson"
)

// Dividend содержит информацию дате закрытия реестра, дивиденде и валюте.
type Dividend struct {
	Ticker   string
	ISIN     string
	Date     time.Time
	Dividend float64
	Currency string
}

func dividendConverter(row *fastjson.Value) (interface{}, error) {
	var (
		div = Dividend{}
		err error
	)

	div.Ticker = string(row.GetStringBytes("secid"))
	div.ISIN = string(row.GetStringBytes("isin"))
	div.Date, err = time.Parse("2006-01-02", string(row.GetStringBytes("registryclosedate")))

	if err != nil {
		return nil, wrapParseErr(err)
	}

	div.Dividend, err = row.Get("value").Float64()

	if err != nil {
		return nil, wrapParseErr(err)
	}

	div.Currency = string(row.GetStringBytes("currencyid"))

	return div, nil
}

// Dividends получает таблицу с дивидендами.
//
// Запрос не отражен в официальном справочнике. По многим инструментам дивиденды отсутствуют или отражены не полностью.
// Корректная информация содержится в основном только по наиболее ликвидным бумагам.
func (iss ISSClient) Dividends(ctx context.Context, security string) (table []Dividend, err error) {
	query := issQuery{
		security:     security,
		object:       "dividends",
		table:        "dividends",
		rowConverter: dividendConverter,
	}

	rows, errors := iss.getRowsGen(ctx, &query)

	for row := range rows {
		table = append(table, row.(Dividend))
	}

	if err := <-errors; err != nil {
		return nil, err
	}

	return table, nil
}
