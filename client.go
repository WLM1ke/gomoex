package gomoex

import (
	"bytes"
	"context"
	"net/http"
	"sync"

	"github.com/valyala/fastjson"
)

// ISSClient клиент для осуществления запросов к MOEX ISS.
type ISSClient struct {
	client  *http.Client
	parsers *fastjson.ParserPool
	buffers *sync.Pool
}

// NewISSClient создает клиент для осуществления запросов к MOEX ISS.
func NewISSClient(client *http.Client) *ISSClient {
	pool := sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	return &ISSClient{
		client:  client,
		parsers: &fastjson.ParserPool{},
		buffers: &pool,
	}
}

func (iss ISSClient) rowGen(ctx context.Context, query *issQuery, rowsc chan interface{}, errc chan error) {
	defer close(rowsc)
	defer close(errc)

	parser := iss.parsers.Get()
	defer iss.parsers.Put(parser)

	start := 0

	buffer := iss.buffers.Get().(*bytes.Buffer)
	defer iss.buffers.Put(buffer)

	for {
		buffer.Reset()

		err := iss.getJSON(ctx, buffer, query, start)
		if err != nil {
			errc <- err

			return
		}

		json, err := getPayload(parser, buffer.Bytes())
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
		return nil, wrapParseErr(err)
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

func (iss ISSClient) getJSON(ctx context.Context, buffer *bytes.Buffer, query *issQuery, start int) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, query.string(start), http.NoBody)
	if err != nil {
		return warpErrWithMsg("can't create request", err)
	}

	resp, err := iss.client.Do(req)
	if err != nil {
		return warpErrWithMsg("can't make request", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			err = warpErrWithMsg("can't close request", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return warpErrWithMsg("got request status", err)
	}

	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return warpErrWithMsg("can't read request", err)
	}

	return nil
}

func (iss ISSClient) getRowsGen(ctx context.Context, query *issQuery) (rows chan interface{}, errc chan error) {
	rows = make(chan interface{})
	errc = make(chan error, 1)

	go iss.rowGen(ctx, query, rows, errc)

	return
}
