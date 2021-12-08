package gomoex

import (
	"context"
	"time"

	"github.com/tidwall/gjson"
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

const (
	_borderBegin    = `begin`
	_borderEnd      = `end`
	_borderInterval = `interval`
)

// CandleBorder содержит информацию о диапазоне доступных дат для свечек заданного интервала.
type CandleBorder struct {
	Begin    time.Time
	End      time.Time
	Interval int
}

func candleBorderConverter(row gjson.Result) (interface{}, error) {
	var (
		boarder CandleBorder
		err     error
	)

	boarder.Begin, err = time.Parse("2006-01-02 15:04:05", row.Get(_borderBegin).String())
	if err != nil {
		return nil, newParseErr(err)
	}

	boarder.End, err = time.Parse("2006-01-02 15:04:05", row.Get(_borderEnd).String())
	if err != nil {
		return nil, newParseErr(err)
	}

	boarder.Interval = int(row.Get(_borderInterval).Int())

	return boarder, nil
}

// MarketCandleBorders получает таблицу с периодами дат рассчитанных свечей для разных по размеру свечек.
//
// Описание запроса - https://iss.moex.com/iss/reference/156
func (iss *ISSClient) MarketCandleBorders(
	ctx context.Context,
	engine,
	market,
	security string,
) (table []CandleBorder, err error) {
	query := querySettings{
		engine:       engine,
		market:       market,
		security:     security,
		object:       "candleborders",
		table:        "borders",
		rowConverter: candleBorderConverter,
	}

	for raw := range iss.getRowsGen(ctx, query.Make()) {
		switch row := raw.(type) {
		case CandleBorder:
			table = append(table, row)
		case error:
			return nil, row
		}
	}

	return table, nil
}

const (
	_candleBegin  = `begin`
	_candleEnd    = `end`
	_candleOpen   = `open`
	_candleClose  = `close`
	_candleHigh   = `high`
	_candleLow    = `low`
	_candleValue  = `value`
	_candleVolume = `volume`
)

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

func candleConverter(row gjson.Result) (interface{}, error) {
	var (
		candle Candle
		err    error
	)

	candle.Begin, err = time.Parse("2006-01-02 15:04:05", row.Get(_candleBegin).String())
	if err != nil {
		return nil, newParseErr(err)
	}

	candle.End, err = time.Parse("2006-01-02 15:04:05", row.Get(_candleEnd).String())
	if err != nil {
		return nil, newParseErr(err)
	}

	candle.Open = row.Get(_candleOpen).Float()
	candle.Close = row.Get(_candleClose).Float()
	candle.High = row.Get(_candleHigh).Float()
	candle.Low = row.Get(_candleLow).Float()
	candle.Value = row.Get(_candleValue).Float()
	candle.Volume = int(row.Get(_candleVolume).Int())

	return candle, nil
}

// MarketCandles свечки данного инструмента и интервала свечки для основного режима данного рынка.
//
// По сравнению со свечками исторические котировки обычно доступны за больший период, но имеются только дневные данные.
// Даты в формате YYYY-MM-DD или пустая строка для получения информации с начала или до конца доступного интервала дат.
// Последняя свечка во время торгов может содержать неполную информацию.
//
// Описание запроса - https://iss.moex.com/iss/reference/155
func (iss *ISSClient) MarketCandles(
	ctx context.Context,
	engine, market, security, from, till string,
	interval int,
) ([]Candle, error) {
	table := make([]Candle, 0)
	query := querySettings{
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

	for raw := range iss.getRowsGen(ctx, query.Make()) {
		switch row := raw.(type) {
		case Candle:
			table = append(table, row)
		case error:
			return nil, row
		}
	}

	return table, nil
}
