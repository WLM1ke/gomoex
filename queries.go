package gomoex

import (
	"strconv"
)

type ISSQuery struct {
	history   bool
	engine    string
	market    string
	board     string
	security  string
	endPoint  string
	table     string
	multipart bool
}

func (query ISSQuery) String(start int) (url string) {
	url = "https://iss.moex.com/iss"

	if query.history {
		url += "/history"
	}
	if query.engine != "" {
		url += "/engines/" + query.engine
	}
	if query.market != "" {
		url += "/markets/" + query.market
	}
	if query.board != "" {
		url += "/boards/" + query.board
	}
	if query.security != "" {
		url += "/securities/" + query.security
	}
	if query.endPoint != "" {
		url += "/" + query.endPoint
	}
	url += ".json?iss.json=extended&iss.meta=off&iss.only=history.cursor," + query.table + "&start=" + strconv.Itoa(start)

	return url
}

func (query ISSQuery) Multipart() bool {
	return query.multipart
}
