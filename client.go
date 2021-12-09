package gomoex

import (
	"bytes"
	"context"
	"errors"
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

var errLastBlock = errors.New("last block")

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

func (iss *ISSClient) rowsGen(ctx context.Context, query issQuery) chan interface{} {
	rows := make(chan interface{})

	go func(query issQuery) {
		defer close(rows)

		for query.err == nil {
			query = iss.parseBlock(ctx, query, rows)
		}

		if !errors.Is(query.err, errLastBlock) {
			rows <- query.err
		}
	}(query)

	return rows
}

func (iss *ISSClient) parseBlock(ctx context.Context, query issQuery, out chan<- interface{}) issQuery {
	buffer := iss.buffers.Get().(*bytes.Buffer) //nolint:forcetypeassert
	defer func() {
		buffer.Reset()
		iss.buffers.Put(buffer)
	}()

	json, err := iss.getJSON(ctx, query.URL(), buffer)
	if err != nil {
		query.err = err

		return query
	}

	return sendTable(json, query, out)
}

func (iss *ISSClient) getJSON(ctx context.Context, url string, buffer *bytes.Buffer) (json []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, newWarpedErr("can't create request", err)
	}

	resp, err := iss.client.Do(req)
	if err != nil {
		return nil, newWarpedErr("can't make request", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			json = nil
			err = newWarpedErr("can't close request", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, newWarpedErr("got request status", err)
	}

	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return nil, newWarpedErr("can't read request", err)
	}

	return buffer.Bytes(), nil
}

func sendTable(json []byte, query issQuery, out chan<- interface{}) issQuery {
	results := gjson.GetManyBytes(json, _payloadPath+query.table, _cursorIndex, _cursorPage, _cursorTotal)

	table := results[0]
	index := results[1]
	page := results[2]
	total := results[3]

	if !table.IsArray() {
		query.err = newErrWithMsg(fmt.Sprintf("can't find tableName %s", query.table))

		return query
	}

	cursorExists := index.Exists() && page.Exists() && total.Exists()

	count, err := sendRows(table, query.rowConverter, out)
	query.start += count
	query.err = err

	switch {
	case count == 0:
		query.err = errLastBlock
	case !query.multipart:
		query.err = errLastBlock
	case cursorExists:
		if index.Int()+page.Int() >= total.Int() {
			query.err = errLastBlock
		}
	}

	return query
}

func sendRows(table gjson.Result, converter converter, out chan<- interface{}) (count int, err error) {
	table.ForEach(
		func(_, rawRow gjson.Result) bool {
			if !rawRow.IsObject() {
				count = 0
				err = newErrWithMsg("can't parse row")

				return false
			}

			var row interface{}

			row, err = converter(rawRow)
			if err != nil {
				count = 0

				return false
			}

			out <- row
			count++

			return true
		},
	)

	return count, nil
}
