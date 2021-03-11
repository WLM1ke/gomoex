package gomoex

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarketCandleBorders(t *testing.T) {
	out := [][]string{
		{"2011-12-15 10:00:00", "2021-03-11 13:59:59", "1"},
		{"2003-07-01 00:00:00", "2021-03-10 00:00:00", "4"},
		{"2003-07-28 00:00:00", "2021-03-10 00:00:00", "7"},
		{"2011-12-08 10:00:00", "2021-03-11 13:56:03", "10"},
		{"2003-07-31 00:00:00", "2021-03-11 13:55:54", "24"},
		{"2003-07-01 00:00:00", "2021-03-10 00:00:00", "31"},
		{"2011-11-17 10:00:00", "2021-03-11 13:56:03", "60"},
	}

	cl := NewISSClient(http.DefaultClient)
	borders, err := cl.MarketCandleBorders(context.Background(), EngineStock, MarketShares, "SNGSP")
	assert.Nil(t, err)
	assert.Equal(t, len(borders), 7)

	for n, border := range borders {
		begin, _ := time.Parse("2006-01-02 15:04:05", out[n][0])
		assert.Equal(t, border.Begin, begin)

		end, _ := time.Parse("2006-01-02 15:04:05", out[n][1])
		assert.True(t, border.End.After(end))

		interval, _ := strconv.Atoi(out[n][2])
		assert.Equal(t, border.Interval, interval)
	}
}

func TestMarketCandlesFromBeginning(t *testing.T) {
	cl := NewISSClient(http.DefaultClient)
	candles, err := cl.MarketCandles(context.Background(), EngineStock, MarketShares, "RTKM", "", "2011-12-16", 1)
	assert.Nil(t, err)
	assert.Greater(t, len(candles), 1030)

	assert.Equal(t, candles[0].Open, 141.55)
	assert.Equal(t, candles[1].Close, 141.59)
	assert.Equal(t, candles[2].High, 142.4)
	assert.Equal(t, candles[3].Low, 140.81)
	assert.Equal(t, candles[4].Value, 2_586_296.9)
	assert.Equal(t, candles[5].Volume, 4140)
	assert.Equal(t, candles[6].Begin, time.Date(2011, 12, 15, 10, 0, 0, 0, time.UTC))
	assert.Equal(t, candles[len(candles)-1].End, time.Date(2011, 12, 16, 18, 44, 59, 0, time.UTC))
}

func TestMarketCandlesToEnd(t *testing.T) {
	cl := NewISSClient(http.DefaultClient)
	candles, err := cl.MarketCandles(context.Background(), EngineStock, MarketShares, "LSRG", "2020-08-20", "", 24)
	assert.Nil(t, err)
	assert.Greater(t, len(candles), 140)

	assert.Equal(t, candles[0].Open, 775.4)
	assert.Equal(t, candles[1].Close, 771.8)
	assert.Equal(t, candles[2].High, 779.8)
	assert.Equal(t, candles[3].Low, 770.2)
	assert.Equal(t, candles[4].Value, 59495740.6)
	assert.Equal(t, candles[6].Begin, time.Date(2020, 8, 28, 0, 0, 0, 0, time.UTC))
}

func TestMarketCandlesEmpty(t *testing.T) {
	cl := NewISSClient(http.DefaultClient)
	candles, err := cl.MarketCandles(context.Background(), EngineStock, MarketShares, "KSGR", "", "", 24)
	assert.Nil(t, err)
	assert.Equal(t, len(candles), 0)
}
