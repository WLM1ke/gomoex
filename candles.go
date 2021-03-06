package gomoex

import (
	"context"
	"fmt"
	"time"

	"github.com/valyala/fastjson"
)

// Доступные интервалы свечек.
const (
	IntervalMin1   = 1
	IntervalMin10  = 10
	IntervalHour   = 60
	IntervalDay    = 24
	IntervalWeek   = 7
	IntervalMonth  = 31
	IntervalQuoter = 4
)

// CandleBorder содержит информацию о диапазоне доступных дат для свечек заданного интервала.
type CandleBorder struct {
	Begin    time.Time
	End      time.Time
	Interval int
}

func candleBorderConverter(row *fastjson.Value) (interface{}, error) {
	var (
		boarder = CandleBorder{}
		err     error
	)

	boarder.Begin, err = time.Parse("2006-01-02 15:04:05", string(row.GetStringBytes("begin")))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "CandleBorder.Begin", err)
	}

	boarder.End, err = time.Parse("2006-01-02 15:04:05", string(row.GetStringBytes("end")))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "CandleBorder.End", err)
	}

	boarder.Interval, err = row.Get("interval").Int()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "CandleBorder.Interval", err)
	}

	return boarder, nil
}

// MarketCandleBorders получает таблицу с периодами дат рассчитанных свечей для разных по размеру свечек.
//
// Описание запроса - https://iss.moex.com/iss/reference/156
func (iss ISSClient) MarketCandleBorders(ctx context.Context, engine, market, security string) (table []CandleBorder, err error) {
	query := issQuery{
		engine:       engine,
		market:       market,
		security:     security,
		object:       "candleborders",
		table:        "borders",
		rowConverter: candleBorderConverter,
	}

	rows, errors := iss.getRowsGen(ctx, &query)

	for row := range rows {
		table = append(table, row.(CandleBorder))
	}

	if err := <-errors; err != nil {
		return nil, err
	}

	return table, nil
}

// Candle представляет исторические котировки в формате OCHL + объем торгов в деньгах и штуках.
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

func candleConverter(row *fastjson.Value) (interface{}, error) {
	var (
		candle = Candle{}
		err    error
	)

	candle.Begin, err = time.Parse("2006-01-02 15:04:05", string(row.GetStringBytes("begin")))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Candle.Begin", err)
	}

	candle.End, err = time.Parse("2006-01-02 15:04:05", string(row.GetStringBytes("end")))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Candle.End", err)
	}

	candle.Open, err = row.Get("open").Float64()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Candle.Open", err)
	}

	candle.Close, err = row.Get("close").Float64()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Candle.Close", err)
	}

	candle.High, err = row.Get("high").Float64()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Candle.High", err)
	}

	candle.Low, err = row.Get("low").Float64()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Candle.Low", err)
	}

	candle.Value, err = row.Get("value").Float64()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Candle.Value", err)
	}

	candle.Volume, err = row.Get("volume").Int()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Candle.Volume", err)
	}

	return candle, nil
}

// MarketCandles свечки данного инструмента и интервала свечки для основного режима данного рынка.
//
// По сравнению со свечками исторические котировки обычно доступны за больший период, но имеются только дневные данные.
// Даты в формате YYYY-MM-DD или пустая строка для получения информации с начала или до конца доступного интервала дат.
// Последняя свечка во время торгов может содержать неполную информацию.
//
// Описание запроса - https://iss.moex.com/iss/reference/155
func (iss ISSClient) MarketCandles(ctx context.Context, engine, market, security, from, till string, interval int) ([]Candle, error) {
	table := make([]Candle, 0)
	query := issQuery{
		engine:       engine,
		market:       market,
		security:     security,
		object:       "candles",
		table:        "candles",
		from:         from,
		till:         till,
		interval:     interval,
		multipart:    true,
		rowConverter: candleConverter,
	}

	rows, errors := iss.getRowsGen(ctx, &query)

	for row := range rows {
		table = append(table, row.(Candle))
	}

	if err := <-errors; err != nil {
		return nil, err
	}

	return table, nil
}
