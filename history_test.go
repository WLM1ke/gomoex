package gomoex

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestMarketDates(t *testing.T) {
	cl := NewISSClient(http.DefaultClient)
	dates, err := cl.MarketDates(context.Background(), EngineStock, MarketShares)
	assert.Nil(t, err)
	assert.Equal(t, len(dates), 1)
	assert.Equal(t, dates[0].From, time.Date(1997, 3, 24, 0, 0, 0, 0, time.UTC))
	assert.True(t, dates[0].Till.After(time.Date(2020, 3, 9, 0, 0, 0, 0, time.UTC)))
}

func TestMarketHistory(t *testing.T) {
	cl := NewISSClient(http.DefaultClient)
	candles, err := cl.MarketHistory(context.Background(), EngineStock, MarketShares, "AKRN", "2017-10-02", "2018-10-12")
	assert.Nil(t, err)
	assert.Equal(t, len(candles), 263)

	assert.Equal(t, candles[0].Date, time.Date(2017, 10, 2, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, candles[1].Open, 3460.0)
	assert.Equal(t, candles[2].Close, 3514.0)
	assert.Equal(t, candles[3].High, 3517.0)
	assert.Equal(t, candles[4].Low, 3510.0)
	assert.Equal(t, candles[5].Value, 3216573.0)
	assert.Equal(t, candles[6].Volume, 928)
	assert.Equal(t, candles[len(candles)-1].Date, time.Date(2018, 10, 12, 0, 0, 0, 0, time.UTC))
}
