package gomoex

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

// Ключевые плейсхолдеры запросов - полный справочник https://iss.moex.com/iss/index.json
const (
	EngineStock    = "stock"    // Фондовый рынок и рынок депозитов
	EngineCurrency = "currency" // Валютный рынок
	EngineFutures  = "futures"  // Срочный рынок

	MarketIndex = "index" // Индексы фондового рынка

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

const dates = "dates"

type ISSDate struct {
	time.Time
}

func (j *ISSDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	j.Time, err = time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	return nil
}

type BoardDates struct {
	From ISSDate `json:"from"`
	Till ISSDate `json:"till"`
}

func (iss ISSClient) GetBoardDates(ctx context.Context, engine string, market string, board string) (table []BoardDates, err error) {
	url := ISSQuery{
		history:  true,
		engine:   engine,
		market:   market,
		board:    board,
		endPoint: dates,
		table:    dates,
	}

	rows, errors := iss.getAll(ctx, url)

	for rawRow := range rows {
		table = append(table, BoardDates{})
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
