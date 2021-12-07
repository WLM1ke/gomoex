package gomoex

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/tidwall/gjson"
)

const (
	_payloadPath = `1.`
	_cursorPath  = `history\.cursor.`
	_cursorIndex = _payloadPath + _cursorPath + `INDEX`
	_cursorPage  = _payloadPath + _cursorPath + `PAGESIZE`
	_cursorTotal = _payloadPath + _cursorPath + `TOTAL`
)

// ISSClient клиент для осуществления запросов к MOEX ISS.
type ISSClient struct {
	client  *http.Client
	buffers sync.Pool
}

// NewISSClient создает клиент для осуществления запросов к MOEX ISS.
func NewISSClient(client *http.Client) *ISSClient {
	iss := ISSClient{client: client}
	iss.buffers.New = func() interface{} { return new(bytes.Buffer) }

	return &iss
}

func (iss *ISSClient) getRowsGen(ctx context.Context, query issQuery) chan interface{} {
	rows := make(chan interface{})

	go iss.rowGen(ctx, query, rows)

	return rows
}

func (iss *ISSClient) rowGen(ctx context.Context, query issQuery, out chan<- interface{}) {
	defer close(out)

	buffer := iss.buffers.Get().(*bytes.Buffer) //nolint:forcetypeassert
	defer iss.buffers.Put(buffer)

	for start := 0; ; {
		buffer.Reset()

		url := query.URL(start)

		err := iss.bufferJSON(ctx, url, buffer)
		if err != nil {
			out <- err

			return
		}

		table, haveNext := extractTable(buffer, query.table)
		if !query.multipart {
			haveNext = false
		}

		nRows := sendRows(table, query, out)

		if !haveNext || nRows == 0 {
			return
		}

		start += nRows
	}
}

func (iss *ISSClient) bufferJSON(ctx context.Context, url string, buffer *bytes.Buffer) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return newWarpedErr("can't create request", err)
	}

	resp, err := iss.client.Do(req)
	if err != nil {
		return newWarpedErr("can't make request", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			err = newWarpedErr("can't close request", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return newWarpedErr("got request status", err)
	}

	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return newWarpedErr("can't read request", err)
	}

	return nil
}

func extractTable(buffer *bytes.Buffer, tableName string) (gjson.Result, bool) {
	results := gjson.GetManyBytes(buffer.Bytes(), _payloadPath+tableName, _cursorIndex, _cursorPage, _cursorTotal)

	table := results[0]
	index := results[1]
	page := results[2]
	total := results[3]

	if !index.Exists() || !page.Exists() || !total.Exists() {
		return table, true
	}

	if index.Int()+page.Int() < total.Int() {
		return table, true
	}

	return table, false
}

func sendRows(table gjson.Result, query issQuery, out chan<- interface{}) int {
	if !table.IsArray() {
		out <- newErrWithMsg(fmt.Sprintf("can't find tableName %s", query.table))

		return 0
	}

	count := 0

	table.ForEach(
		func(_, rawRow gjson.Result) bool {
			if !rawRow.IsObject() {
				count = 0

				return false
			}

			row, err := query.rowConverter(rawRow)
			if err != nil {
				out <- err
				count = 0

				return false
			}

			out <- row
			count++

			return true
		},
	)

	return count
}
