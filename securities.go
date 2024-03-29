package gomoex

import (
	"context"

	"github.com/tidwall/gjson"
)

// Security содержит информацию о ценной бумаге.
type Security struct {
	Ticker     string
	LotSize    int
	ISIN       string
	Board      string
	Type       string
	Instrument string
}

const (
	_securitySECID      = `SECID`
	_securityLotSize    = `LOTSIZE`
	_securityISIN       = `ISIN`
	_securityBoard      = `BOARDID`
	_securityType       = `SECTYPE`
	_securityInstrument = `INSTRID`
)

func securityConverter(row gjson.Result) (interface{}, error) {
	var sec Security

	sec.Ticker = row.Get(_securitySECID).String()
	sec.LotSize = int(row.Get(_securityLotSize).Int())
	sec.ISIN = row.Get(_securityISIN).String()
	sec.Board = row.Get(_securityBoard).String()
	sec.Type = row.Get(_securityType).String()
	sec.Instrument = row.Get(_securityInstrument).String()

	return sec, nil
}

// BoardSecurities получает таблицу с торгуемыми бумагами в данном режиме торгов.
//
// Описание запроса - https://iss.moex.com/iss/reference/32
func (iss *ISSClient) BoardSecurities(ctx context.Context, engine, market, board string) (table []Security, err error) {
	query := querySettings{
		engine:       engine,
		market:       market,
		board:        board,
		object:       "securities",
		table:        "securities",
		rowConverter: securityConverter,
	}

	for raw := range iss.rowsGen(ctx, query.Make()) {
		switch row := raw.(type) {
		case Security:
			table = append(table, row)
		case error:
			return nil, row
		}
	}

	return table, nil
}
