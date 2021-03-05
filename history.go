package gomoex

import (
	"context"
	"encoding/json"
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
	From string `json:"from"`
	Till string `json:"till"`
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
		err = json.Unmarshal(rawRow, &table[len(table)-1])
		if err != nil {
			return nil, err
		}
	}
	if err = <-errors; err != nil {
		return nil, err
	}
	return table, nil
}

type Quotes struct {
	Date  string  `json:"TRADEDATE"`
	Close float64 `json:"CLOSE"`
}

// MarketHistory исторические котировки для данного инструмента и всех торговоых режимов для данного рынка.
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
		err = json.Unmarshal(rawRow, &table[len(table)-1])
		if err != nil {
			return nil, err
		}
	}
	if err = <-errors; err != nil {
		return nil, err
	}
	return table, nil
}
