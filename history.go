package gomoex

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/valyala/fastjson"
)

// Date содержит информацию о диапазоне доступных торговых дат.
type Date struct {
	From time.Time
	Till time.Time
}

func dateConverter(row *fastjson.Value) (interface{}, error) {
	var (
		date = Date{}
		err  error
	)

	date.From, err = time.Parse("2006-01-02", string(row.GetStringBytes("from")))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Date.From", err)
	}

	date.Till, err = time.Parse("2006-01-02", string(row.GetStringBytes("till")))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Date.Till", err)
	}

	return date, nil
}

// MarketDates получает таблицу с диапазоном дат с доступными данными для данного рынка.
//
// Описание запроса - https://iss.moex.com/iss/reference/83
func (iss ISSClient) MarketDates(ctx context.Context, engine, market string) (table []Date, err error) {
	query := issQuery{
		history:      true,
		engine:       engine,
		market:       market,
		object:       "dates",
		table:        "dates",
		rowConverter: dateConverter,
	}

	rows, errors := iss.getRowsGen(ctx, &query)

	for row := range rows {
		table = append(table, row.(Date))
	}

	if err := <-errors; err != nil {
		return nil, err
	}

	return table, nil
}

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

func convertToNanFloat(value *fastjson.Value) (float64, error) {
	if value.Type() == fastjson.TypeNull {
		return math.NaN(), nil
	}

	return value.Float64()
}

func quoteConverter(row *fastjson.Value) (interface{}, error) {
	var (
		quote = Quote{}
		err   error
	)

	quote.Date, err = time.Parse("2006-01-02", string(row.GetStringBytes("TRADEDATE")))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Quote.Date", err)
	}

	quote.Open, err = convertToNanFloat(row.Get("OPEN"))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Quote.Open", err)
	}

	quote.Close, err = convertToNanFloat(row.Get("CLOSE"))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Quote.Close", err)
	}

	quote.High, err = convertToNanFloat(row.Get("HIGH"))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Quote.High", err)
	}

	quote.Low, err = convertToNanFloat(row.Get("LOW"))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Quote.Low", err)
	}

	quote.Value, err = row.Get("VALUE").Float64()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Quote.Value", err)
	}

	quote.Volume, err = row.Get("VOLUME").Int()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Quote.Volume", err)
	}

	return quote, nil
}

// MarketHistory исторические котировки данного инструмента для всех торговых режимов для данного рынка.
//
// По сравнению со свечками обычно доступны за больший период, но имеются только дневные данные.
// Описание запроса - https://iss.moex.com/iss/reference/63
func (iss ISSClient) MarketHistory(ctx context.Context, engine, market, security, from, till string) (table []Quote, err error) {
	query := issQuery{
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

	rows, errors := iss.getRowsGen(ctx, &query)

	for row := range rows {
		table = append(table, row.(Quote))
	}

	if err := <-errors; err != nil {
		return nil, err
	}

	return table, nil
}
