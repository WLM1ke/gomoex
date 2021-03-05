package gomoex

import (
	"context"
	"fmt"
	"github.com/francoispqt/gojay"
	"time"
)

type Dates struct {
	From time.Time
	Till time.Time
}

func (date *Dates) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "from":
		return dec.Time(&date.From, "2006-01-02")
	case "till":
		return dec.Time(&date.Till, "2006-01-02")
	}
	return nil
}
func (date *Dates) NKeys() int {
	return 2
}

// MarketDates получает таблицу с диапазоном торговых дат для данного рынка.
// Описание запроса - https://iss.moex.com/iss/reference/83
func (iss ISSClient) MarketDates(ctx context.Context, engine string, market string) (table []Dates, err error) {
	query := ISSQuery{
		history: true,
		engine:  engine,
		market:  market,
		object:  "dates",
		table:   "dates",
	}

	rows, errors := iss.getAll(ctx, query)

	for rawRow := range rows {
		table = append(table, Dates{})
		err = gojay.Unmarshal(rawRow, &table[len(table)-1])
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	if err = <-errors; err != nil {

		return nil, err
	}
	return table, nil
}

// Quotes представляет исторические дневные котировки в формате OCHL + объем торгов в деньгах и штуках.
type Quotes struct {
	Date   time.Time
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Value  float64
	Volume int
}

func (quotes *Quotes) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "TRADEDATE":
		return dec.Time(&quotes.Date, "2006-01-02")
	case "OPEN":
		return dec.Float(&quotes.Open)
	case "CLOSE":
		return dec.Float(&quotes.Close)
	case "HIGH":
		return dec.Float(&quotes.High)
	case "LOW":
		return dec.Float(&quotes.Low)
	case "VALUE":
		return dec.Float(&quotes.Value)
	case "VOLUME":
		return dec.Int(&quotes.Volume)
	}
	return nil
}
func (quotes *Quotes) NKeys() int {
	return 7
}

// MarketHistory исторические котировки данного инструмента для всех торговых режимов для данного рынка.
//
// По сравнению со свечками обычно доступны за больший период, но имеются только дневные данные.
// Описание запроса - https://iss.moex.com/iss/reference/63
func (iss ISSClient) MarketHistory(ctx context.Context, engine string, market string, security string) (table []Quotes, err error) {
	query := ISSQuery{
		history:   true,
		engine:    engine,
		market:    market,
		security:  security,
		table:     "history",
		multipart: true,
	}

	rows, errors := iss.getAll(ctx, query)

	for rawRow := range rows {
		table = append(table, Quotes{})
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
