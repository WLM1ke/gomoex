package gomoex

import (
	"context"
	"github.com/valyala/fastjson"
)

// Security содержит информацию о ценной бумаге.
type Security struct {
	Ticker  string
	LotSize int
	ISIN    string
}

func securityConverter(row *fastjson.Value) (interface{}, error) {

	sec := Security{}
	var err error

	sec.Ticker = string(row.GetStringBytes("SECID"))
	sec.LotSize, err = row.Get("LOTSIZE").Int()
	if err != nil {
		return nil, err
	}
	sec.ISIN = string(row.GetStringBytes("ISIN"))

	return sec, nil
}

// BoardSecurities получает таблицу с торгуемыми бумагами в данном режиме торгов.
// Описание запроса - https://iss.moex.com/iss/reference/32
func (iss ISSClient) BoardSecurities(ctx context.Context, engine string, market string, board string) (table []Security, err error) {
	query := issQuery{
		engine:       engine,
		market:       market,
		board:        board,
		object:       "securities",
		table:        "securities",
		rowConverter: securityConverter,
	}

	rows, errors := iss.getRowsGen(ctx, query)

	for row := range rows {
		table = append(table, row.(Security))
	}

	if err = <-errors; err != nil {
		return nil, err
	}

	return table, nil
}
