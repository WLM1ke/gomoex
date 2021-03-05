package gomoex

// Ключевые плейсхолдеры запросов — полный справочник https://iss.moex.com/iss/index.json
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
