package gomoex

import (
	"context"
	"github.com/francoispqt/gojay"
)

// Security содержит информацию о ценной бумаге.
type Security struct {
	Ticker  string
	LotSize int
	ISIN    string
}

func (security *Security) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "SECID":
		return dec.String(&security.Ticker)
	case "LOTSIZE":
		return dec.Int(&security.LotSize)
	case "ISIN":
		return dec.String(&security.ISIN)
	}
	return nil
}
func (security *Security) NKeys() int {
	return 3
}

// BoardSecurities получает таблицу с торгуемыми бумагами в данном режиме торгов.
// Описание запроса - https://iss.moex.com/iss/reference/32
func (iss ISSClient) BoardSecurities(ctx context.Context, engine string, market string, board string) (table []Security, err error) {
	query := ISSQuery{
		engine: engine,
		market: market,
		board:  board,
		object: "securities",
		table:  "securities",
	}

	rows, errors := iss.getAll(ctx, query)

	for rawRow := range rows {
		table = append(table, Security{})
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

// SearchResult содержит результат поиска ценной бумаги.
type SearchResult struct {
	Ticker string
	ISIN   string
}

func (result *SearchResult) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "secid":
		return dec.String(&result.Ticker)
	case "isin":
		return dec.String(&result.ISIN)
	}
	return nil
}
func (result *SearchResult) NKeys() int {
	return 2
}

// FindSecurity осуществляет поиск ценной бумаги по тикеру, ISIN и прочей информации.
// Поисковые запросы длиной менее трёх букв игнорируются ISS. Если параметром передано два слова через пробел. То каждое
// должно быть длиной не менее трёх букв.
// Описание запроса - https://iss.moex.com/iss/reference/5
func (iss ISSClient) FindSecurity(ctx context.Context, q string) (table []SearchResult, err error) {
	query := ISSQuery{
		object:    "securities",
		table:     "securities",
		q:         q,
		multipart: true,
	}

	rows, errors := iss.getAll(ctx, query)

	for rawRow := range rows {
		table = append(table, SearchResult{})
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
