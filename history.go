package gomoex

import (
	"context"
	"fmt"
	"github.com/francoispqt/gojay"
	"time"
)

// Ключевые плейсхолдеры запросов - полный справочник https://iss.moex.com/iss/index.json
const (
	EngineStock    = "stock"    // Фондовый рынок и рынок депозитов
	EngineCurrency = "currency" // Валютный рынок
	EngineFutures  = "futures"  // Срочный рынок

	MarketIndex         = "index"         // Индексы фондового рынка
	MarketShares        = "shares"        // Рынок акций
	MarketBonds         = "bonds"         // Рынок облигаций
	MarketForeignShares = "foreignshares" // Иностранные ц.б.

	MarketSelt    = "selt"    // Биржевые сделки с ЦК
	MarketFutures = "futures" //Поставочные фьючерсы

	MarketFORTS   = "forts"   // ФОРТС
	MarketOptions = "options" //Опционы ФОРТС

	BoardTQBR = "TQBR" // Т+: Акции и ДР - безадрес.
	BoardTQTF = "TQTF" // Т+: ETF - безадрес.
)

// Ключевые таблицы
const (
	dates   = "dates"
	history = "history"
)

type Dates struct {
	From time.Time
	Till time.Time
}

func (date *Dates) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "from":
		return dec.Time(&date.From, "2006-01-02")
	case "till":
		return dec.Time(&date.Till, "2006-01-02")
	}
	return nil
}
func (date *Dates) NKeys() int {
	return 2
}

// MarketDates получает таблицу с биапазоном торговых дат для данного рынка.
// Описание запроса - https://iss.moex.com/iss/reference/83
func (iss ISSClient) MarketDates(ctx context.Context, engine string, market string) (table []Dates, err error) {
	query := ISSQuery{
		history: true,
		engine:  engine,
		market:  market,
		object:  dates,
		table:   dates,
	}

	rows, errors := iss.getAll(ctx, query)

	for rawRow := range rows {
		table = append(table, Dates{})
		err = gojay.Unmarshal(rawRow, &table[len(table)-1])
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	if err = <-errors; err != nil {

		return nil, err
	}
	return table, nil
}

type Quotes struct {
	Date  time.Time
	Close float64
}

func (quotes *Quotes) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "TRADEDATE":
		return dec.Time(&quotes.Date, "2006-01-02")
	case "CLOSE":
		return dec.Float(&quotes.Close)
	}
	return nil
}
func (quotes *Quotes) NKeys() int {
	return 2
}

// MarketHistory исторические котировки для данного инструмента и всех торговоых режимов для данного рынка.
// По сравнению со свечками имеют доступны за больший период, но имеются только дневные данные.
// Описание запроса - https://iss.moex.com/iss/reference/63
func (iss ISSClient) MarketHistory(ctx context.Context, engine string, market string, security string) (table []Quotes, err error) {
	query := ISSQuery{
		history:   true,
		engine:    engine,
		market:    market,
		security:  security,
		table:     history,
		multipart: true,
	}

	rows, errors := iss.getAll(ctx, query)

	for rawRow := range rows {
		table = append(table, Quotes{})
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
