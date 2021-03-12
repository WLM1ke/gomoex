package gomoex

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/valyala/fastjson"
)

// ISSClient клиент для осуществления запросов к MOEX ISS.
type ISSClient struct {
	client  *http.Client
	parsers *fastjson.ParserPool
}

// NewISSClient создает клиент для осуществления запросов к MOEX ISS.
func NewISSClient(client *http.Client) *ISSClient {
	return &ISSClient{client, &fastjson.ParserPool{}}
}

func (iss ISSClient) rowGen(ctx context.Context, query *issQuery, rowsc chan interface{}, errc chan error) {
	defer close(rowsc)
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

		json, err := getPayload(parser, data)
		if err != nil {
			errc <- err

			return
		}

		nRows, err := yieldRows(json, query, rowsc)
		if err != nil {
			errc <- err

			return
		}

		if !query.multipart || nRows == 0 {
			return
		}

		if !haveNextBlock(json) {
			return
		}

		start += nRows
	}
}

func yieldRows(json *fastjson.Value, query *issQuery, rows chan interface{}) (int, error) {
	rawRows := json.GetArray(query.table)
	for n, rawRow := range rawRows {
		row, err := query.rowConverter(rawRow)
		if err != nil {
			return n, err
		}
		rows <- row
	}

	return len(rawRows), nil
}

// Полезные данные в первом элементе массива - в нулевом бесполезные данные о кодировке.
func getPayload(parser *fastjson.Parser, data []byte) (*fastjson.Value, error) {
	json, err := parser.ParseBytes(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Parse payload", err)
	}

	return json.Get("1"), nil
}

func haveNextBlock(json *fastjson.Value) bool {
	curData := json.Get("history.cursor", "0")
	if curData == nil {
		return true
	}

	if curData.GetInt("INDEX")+curData.GetInt("PAGESIZE") < curData.GetInt("TOTAL") {
		return true
	}

	return false
}

// ErrISSBadStatus - ответ ISS отличается от 200 OK.
var ErrISSBadStatus = errors.New("bad status code")

func (iss ISSClient) getJSON(ctx context.Context, query *issQuery, start int) (data []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, query.string(start), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "New ISS Request", err)
	}

	resp, err := iss.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Do ISS Request", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			data = nil
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s", ErrISSBadStatus, resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func (iss ISSClient) getRowsGen(ctx context.Context, query *issQuery) (rows chan interface{}, errc chan error) {
	rows = make(chan interface{})
	errc = make(chan error, 1)

	go iss.rowGen(ctx, query, rows, errc)

	return
}
