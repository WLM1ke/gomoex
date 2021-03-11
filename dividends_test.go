package gomoex

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestDividends(t *testing.T) {
	cl := NewISSClient(http.DefaultClient)
	div, err := cl.Dividends(context.Background(), "TATN")
	assert.Nil(t, err)
	assert.Greater(t, len(div), 7)

	assert.Equal(t, div[0].Ticker, "TATN")
	assert.Equal(t, div[1].ISIN, "RU0009033591")
	assert.Equal(t, div[2].Date, time.Date(2019, 1, 9, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, div[3].Currency, "RUB")
	assert.Equal(t, div[7].Dividend, 9.94)
}
