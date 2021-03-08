package gomoex

import (
	"context"
	"fmt"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
)

type ISSClient struct {
	Client *http.Client
	pool   fastjson.ParserPool
}

func (iss *ISSClient) get(ctx context.Context, query issQuery, rows chan interface{}, errors chan error) {

	defer close(rows)
	defer close(errors)

	parser := iss.pool.Get()
	defer iss.pool.Put(parser)

	blockSize := 1
	start := -1

	for blockSize != 0 {
		start += blockSize

		data, err := iss.getJSON(ctx, query, start)
		if err != nil {
			errors <- err
			return
		}

		json, err := parser.ParseBytes(data)
		if err != nil {
			errors <- err
			return
		}

		json = json.Get("1")

		rawRows := json.GetArray(query.table)
		if rawRows == nil {
			errors <- err
			return
		}

		for _, rawRow := range rawRows {
			row, err := query.rowConverter(rawRow)
			if err != nil {
				errors <- err
				return
			}
			rows <- row
		}

		if !query.multipart {
			return
		}

		blockSize = len(rawRows)

		curData := json.Get("history.cursor", "0")
		if curData == nil {
			continue
		}

		if start+curData.GetInt("PAGESIZE") >= curData.GetInt("TOTAL") {
			blockSize = 0
		}
	}

}

func (iss *ISSClient) getJSON(ctx context.Context, query issQuery, start int) (data []byte, err error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, query.String(start), nil)
	if err != nil {
		return nil, err
	}

	resp, err := iss.Client.Do(req)
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
		return nil, fmt.Errorf("bad status %s", resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func (iss *ISSClient) getAll(ctx context.Context, query issQuery) (rows chan interface{}, errors chan error) {
	rows = make(chan interface{})
	errors = make(chan error, 1)

	go iss.get(ctx, query, rows, errors)

	return
}
