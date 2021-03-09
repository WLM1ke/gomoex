package gomoex

import (
	"context"
	"errors"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
)

// ISSClient клиент для осуществления запросов к MOEX ISS
type ISSClient struct {
	client  *http.Client
	parsers *fastjson.ParserPool
}

// NewISSClient создает клиент для осуществления запросов к MOEX ISS
func NewISSClient(client *http.Client, parsers *fastjson.ParserPool) *ISSClient {
	return &ISSClient{client, parsers}
}

func (iss ISSClient) rowGen(ctx context.Context, query *issQuery, rows chan interface{}, errc chan error) {

	defer close(rows)
	defer close(errc)

	parser := iss.parsers.Get()
	defer iss.parsers.Put(parser)

	start := 0

	for {

		data, err := iss.getJSON(ctx, query, start)
		if err != nil {
			errc <- err
			return
		}

		json, err := parser.ParseBytes(data)
		if err != nil {
			errc <- err
			return
		}

		// Полезные данные в первом элементе массива
		json = json.Get("1")

		rawRows := json.GetArray(query.table)

		for _, rawRow := range rawRows {
			row, err := query.rowConverter(rawRow)
			if err != nil {
				errc <- err
				return
			}
			rows <- row
		}

		if !query.multipart || len(rawRows) == 0 {
			return
		}

		if !iss.loadNextBlock(json) {
			return
		}

		start += len(rawRows)
	}

}

func (iss ISSClient) loadNextBlock(json *fastjson.Value) bool {
	curData := json.Get("history.cursor", "0")
	if curData == nil {
		return true
	}

	if curData.GetInt("INDEX")+curData.GetInt("PAGESIZE") < curData.GetInt("TOTAL") {
		return true
	}
	return false
}

func (iss ISSClient) getJSON(ctx context.Context, query *issQuery, start int) (data []byte, err error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, query.string(start), nil)
	if err != nil {
		return nil, err
	}

	resp, err := iss.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			data = nil
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func (iss ISSClient) getRowsGen(ctx context.Context, query *issQuery) (rows chan interface{}, errc chan error) {
	rows = make(chan interface{})
	errc = make(chan error, 1)

	go iss.rowGen(ctx, query, rows, errc)

	return
}
