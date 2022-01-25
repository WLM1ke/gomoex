package gomoex

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoardSecurities(t *testing.T) {
	cl := NewISSClient(http.DefaultClient)
	sec, err := cl.BoardSecurities(context.Background(), EngineStock, MarketShares, BoardTQBR)
	assert.Nil(t, err)
	assert.Greater(t, len(sec), 250)
	assert.Equal(t, sec[0].Ticker, "ABRD")
	assert.Equal(t, sec[1].LotSize, 100)
	assert.Equal(t, sec[2].Board, BoardTQBR)
	assert.Equal(t, sec[3].Instrument, "EQIN")
	assert.Equal(t, sec[14].Type, "2")
	assert.Equal(t, sec[len(sec)-1].ISIN, "RU0009091300")
}
