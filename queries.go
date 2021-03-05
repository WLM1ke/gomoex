package gomoex

import (
	"strconv"
	"strings"
)

// ISSQuery содержит описание запроса к ISS и позволяет сформировать необходимый для его осуществления URL.
// Официальный справочник запросов https://iss.moex.com/iss/reference/
// Официальный справочник разработчика https://fs.moex.com/files/6523
type ISSQuery struct {
	// Нужен ли префикс history.
	history bool
	// Значение плейсхолдера engine. Для пустой строки не добавляется в запрос.
	engine string
	// Значение плейсхолдера market. Для пустой строки не добавляется в запрос.
	market string
	// Значение плейсхолдера board. Для пустой строки не добавляется в запрос.
	board string
	// Значение плейсхолдера security. Для пустой строки не добавляется в запрос.
	security string
	// Запрашиваемый объект. Для пустой строки не добавляется в запрос.
	object string
	// Запрашиваемая таблица внутри ответа.
	table string
	// Интервал свечек.
	interval int
	// Поисковый запрос о ценной бумаге.
	q string
	// Будет ли ответ разбит на несколько блоков, требующих последовательной загрузки со смещением стартовой позиции.
	multipart bool
}

// String формирует URL запроса на основании описания для заданной стартовой позиции.
// В базовый URL добавляются требование предоставить расширенный JSON без метаданных с таблицей курсора.
func (query ISSQuery) String(start int) (url string) {
	urlParts := []string{"https://iss.moex.com/iss"}

	if query.history {
		urlParts = append(urlParts, "/history")
	}
	if query.engine != "" {
		urlParts = append(urlParts, "/engines/", query.engine)
	}
	if query.market != "" {
		urlParts = append(urlParts, "/markets/", query.market)
	}
	if query.board != "" {
		urlParts = append(urlParts, "/boards/", query.board)
	}
	if query.security != "" {
		urlParts = append(urlParts, "/securities/", query.security)
	}
	if query.object != "" {
		urlParts = append(urlParts, "/", query.object)
	}
	urlParts = append(urlParts, ".json?iss.json=extended&iss.meta=off")
	urlParts = append(urlParts, "&iss.only=history.cursor,", query.table)
	urlParts = append(urlParts, "&interval=", strconv.Itoa(query.interval))
	urlParts = append(urlParts, "&q=", query.q)
	urlParts = append(urlParts, "&start=", strconv.Itoa(start))

	return strings.Join(urlParts, "")
}
