package gomoex

import (
	"context"
	"github.com/francoispqt/gojay"
	"time"
)

// CandleBorders содержит информацию о диапазоне доступных дат для свечек заданного интервала.
type CandleBorders struct {
	Begin    time.Time
	End      time.Time
	Interval int
}

func (boarders *CandleBorders) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "begin":
		return dec.Time(&boarders.Begin, "2006-01-02 15:04:05")
	case "end":
		return dec.Time(&boarders.End, "2006-01-02 15:04:05")
	case "interval":
		return dec.Int(&boarders.Interval)
	}
	return nil
}
func (boarders *CandleBorders) NKeys() int {
	return 3
}

// MarketCandleBorders получает таблицу с диапазонами доступных свечек для данного рынка и тикера.
// Описание запроса - https://iss.moex.com/iss/reference/156
func (iss ISSClient) MarketCandleBorders(ctx context.Context, engine string, market string, security string) (table []CandleBorders, err error) {
	query := ISSQuery{
		engine:   engine,
		market:   market,
		security: security,
		object:   "candleborders",
		table:    "borders",
	}

	rows, errors := iss.getAll(ctx, query)

	for rawRow := range rows {
		table = append(table, CandleBorders{})
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

// Candle представляет исторические дневные котировки в формате OCHL + объем торгов в деньгах и штуках.
type Candle struct {
	Begin  time.Time
	End    time.Time
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Value  float64
	Volume int
}

func (candle *Candle) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "begin":
		return dec.Time(&candle.Begin, "2006-01-02 15:04:05")
	case "end":
		return dec.Time(&candle.End, "2006-01-02 15:04:05")
	case "open":
		return dec.Float(&candle.Open)
	case "close":
		return dec.Float(&candle.Close)
	case "high":
		return dec.Float(&candle.High)
	case "low":
		return dec.Float(&candle.Low)
	case "value":
		return dec.Float(&candle.Value)
	case "volume":
		return dec.Int(&candle.Volume)
	}
	return nil
}
func (candle *Candle) NKeys() int {
	return 8
}

// Доступные интервалы свечек
const (
	IntervalMin1   = 1
	IntervalMin10  = 10
	IntervalHour   = 60
	IntervalDay    = 24
	IntervalWeek   = 7
	IntervalMonth  = 31
	IntervalQuoter = 4
)

// MarketCandles свечки данного инструмента и интервала свечки для основного режима данного рынка.
//
// По сравнению со свечками исторические котировки обычно доступны за больший период, но имеются только дневные данные.
// Описание запроса - https://iss.moex.com/iss/reference/155
func (iss ISSClient) MarketCandles(ctx context.Context, engine string, market string, security string, interval int) (table []Candle, err error) {
	query := ISSQuery{
		engine:    engine,
		market:    market,
		security:  security,
		object:    "candles",
		table:     "candles",
		interval:  interval,
		multipart: true,
	}

	rows, errors := iss.getAll(ctx, query)

	for rawRow := range rows {
		table = append(table, Candle{})
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
