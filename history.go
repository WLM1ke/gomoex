package gomoex

import (
	"context"
	"math"
	"time"

	"github.com/tidwall/gjson"
)

const (
	_dateFrom = `from`
	_dateTill = `till`
)

// Date содержит информацию о диапазоне доступных торговых дат.
type Date struct {
	From time.Time
	Till time.Time
}

func dateConverter(row gjson.Result) (interface{}, error) {
	var (
		date Date
		err  error
	)

	date.From, err = time.Parse("2006-01-02", row.Get(_dateFrom).String())
	if err != nil {
		return nil, newParseErr(err)
	}

	date.Till, err = time.Parse("2006-01-02", row.Get(_dateTill).String())
	if err != nil {
		return nil, newParseErr(err)
	}

	return date, nil
}

// MarketDates получает таблицу с диапазоном дат с доступными данными для данного рынка.
//
// Описание запроса - https://iss.moex.com/iss/reference/83
func (iss *ISSClient) MarketDates(ctx context.Context, engine, market string) (table []Date, err error) {
	query := querySettings{
		history:      true,
		engine:       engine,
		market:       market,
		object:       "dates",
		table:        "dates",
		rowConverter: dateConverter,
	}

	for raw := range iss.getRowsGen(ctx, query.Make()) {
		switch row := raw.(type) {
		case Date:
			table = append(table, row)
		case error:
			return nil, row
		}
	}

	return table, nil
}

const (
	_quoteDate   = `TRADEDATE`
	_quoteOpen   = `OPEN`
	_quoteClose  = `CLOSE`
	_quoteHigh   = `HIGH`
	_quoteLow    = `LOW`
	_quoteValue  = `VALUE`
	_quoteVolume = `VOLUME`
)

// Quote представляет исторические дневные котировки в формате OCHL + объем торгов в деньгах и штуках.
type Quote struct {
	Date   time.Time
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Value  float64
	Volume int
}

func convertToNanFloat(value gjson.Result) float64 {
	if value.Type == gjson.Null {
		return math.NaN()
	}

	return value.Float()
}

func quoteConverter(row gjson.Result) (interface{}, error) {
	var (
		quote Quote
		err   error
	)

	quote.Date, err = time.Parse("2006-01-02", row.Get(_quoteDate).String())
	if err != nil {
		return nil, newParseErr(err)
	}

	quote.Open = convertToNanFloat(row.Get(_quoteOpen))
	quote.Close = convertToNanFloat(row.Get(_quoteClose))
	quote.High = convertToNanFloat(row.Get(_quoteHigh))
	quote.Low = convertToNanFloat(row.Get(_quoteLow))
	quote.Value = row.Get(_quoteValue).Float()
	quote.Volume = int(row.Get(_quoteVolume).Int())

	return quote, nil
}

// MarketHistory исторические котировки данного инструмента для всех торговых режимов для данного рынка.
//
// По сравнению со свечками обычно доступны за больший период, но имеются только дневные данные.
// Даты в формате YYYY-MM-DD или пустая строка для получения информации с начала или до конца доступного интервала дат.
//
// Описание запроса - https://iss.moex.com/iss/reference/63
func (iss *ISSClient) MarketHistory(
	ctx context.Context,
	engine,
	market,
	security,
	from,
	till string,
) (table []Quote, err error) {
	query := querySettings{
		history:      true,
		engine:       engine,
		market:       market,
		security:     security,
		table:        "history",
		from:         from,
		till:         till,
		multipart:    true,
		rowConverter: quoteConverter,
	}

	for raw := range iss.getRowsGen(ctx, query.Make()) {
		switch row := raw.(type) {
		case Quote:
			table = append(table, row)
		case error:
			return nil, row
		}
	}

	return table, nil
}
