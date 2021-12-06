package gomoex

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	_issURL     = `https://iss.moex.com/iss`
	_history    = `/history`
	_engines    = `/engines/`
	_markets    = `/markets/`
	_boards     = `/boards/`
	_securities = `/securities/`
	_object     = `/`
	_query      = `.json?iss.json=extended&iss.meta=off&iss.only=history.cursor,`
	_from       = `&from=`
	_till       = `&till=`
	_interval   = `&interval=`
	_start      = `&start=%d`
)

// converter преобразует строку таблицы из json.
type converter func(row gjson.Result) (interface{}, error)

// issQuery содержит данные необходимые для осуществления запроса.
type issQuery struct {
	fmt          string
	table        string
	multipart    bool
	rowConverter converter
}

func (query issQuery) URL(start int) string {
	return fmt.Sprintf(query.fmt, start)
}

// querySettings содержит настройки создания запроса к ISS.
//
// Позволяет сформировать необходимый для его осуществления URL и вспомогательные данные.
//
// Официальный справочник запросов https://iss.moex.com/iss/reference/
// Официальный справочник разработчика https://fs.moex.com/files/6523
type querySettings struct {
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
	// Дата, с которой выводить данные в формате ГГГГ-ММ-ДД.
	from string
	// Дата, до которой выводить данные в формате ГГГГ-ММ-ДД.
	till string
	// Интервал свечек.
	interval int
	// Будет ли ответ разбит на несколько блоков, требующих последовательной загрузки со смещением стартовой позиции.
	multipart bool
	// Конвертор данных — выбирает необходимые поля и преобразует данные.
	rowConverter converter
}

// Make формирует URL запроса на основании описания для заданной стартовой позиции.
// В базовый URL добавляются требование предоставить расширенный JSON без метаданных с таблицей курсора.
func (query *querySettings) Make() issQuery {
	var url strings.Builder

	url.WriteString(_issURL)

	if query.history {
		url.WriteString(_history)
	}

	if query.engine != "" {
		url.WriteString(_engines)
		url.WriteString(query.engine)
	}

	if query.market != "" {
		url.WriteString(_markets)
		url.WriteString(query.market)
	}

	if query.board != "" {
		url.WriteString(_boards)
		url.WriteString(query.board)
	}

	if query.security != "" {
		url.WriteString(_securities)
		url.WriteString(query.security)
	}

	if query.object != "" {
		url.WriteString(_object)
		url.WriteString(query.object)
	}

	url.WriteString(_query)
	url.WriteString(query.table)

	if query.from != "" {
		url.WriteString(_from)
		url.WriteString(query.from)
	}

	if query.till != "" {
		url.WriteString(_till)
		url.WriteString(query.till)
	}

	if query.interval != 0 {
		url.WriteString(_interval)
		url.WriteString(strconv.Itoa(query.interval))
	}

	url.WriteString(_start)

	return issQuery{
		fmt:          url.String(),
		table:        query.table,
		multipart:    query.multipart,
		rowConverter: query.rowConverter,
	}
}
