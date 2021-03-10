// Package gomoex реализует часть запросов к MOEX ISS
// Официальный справочник запросов https://iss.moex.com/iss/reference/
// Официальный справочник разработчика https://fs.moex.com/files/6523
package gomoex

// Ключевые плейсхолдеры запросов — полный справочник https://iss.moex.com/iss/index.json
const (
	EngineStock    = "stock"    // Фондовый рынок и рынок депозитов
	EngineCurrency = "currency" // Валютный рынок
	EngineFutures  = "futures"  // Срочный рынок

	MarketIndex         = "index"         // Индексы фондового рынка
	MarketShares        = "shares"        // Рынок акций
	MarketBonds         = "bonds"         // Рынок облигаций
	MarketForeignShares = "foreignshares" // Иностранные ценные бумаги
	MarketSelt          = "selt"          // Биржевые сделки с ЦК
	MarketFutures       = "futures"       // Поставочные фьючерсы
	MarketFORTS         = "forts"         // ФОРТС
	MarketOptions       = "options"       // Опционы ФОРТС

	BoardTQBR = "TQBR" // Т+: Акции и ДР — безадресные сделки
	BoardTQTF = "TQTF" // Т+: ETF — безадресные сделки
	BoardFQBR = "FQBR" // Т+ Иностранные Акции и ДР — безадресные сделки
)
