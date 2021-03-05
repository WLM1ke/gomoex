package gomoex

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type cursor struct {
	Index    int `json:"INDEX"`
	Total    int `json:"TOTAL"`
	PageSize int `json:"PAGESIZE"`
}

type ISSClient struct {
	Client *http.Client
}

func (iss *ISSClient) get(ctx context.Context, query issQuery, rows chan json.RawMessage, errors chan error) {

	defer close(rows)
	defer close(errors)

	blockSize := 1
	start := -1

	var rawRows []json.RawMessage
	var cur []cursor

	for blockSize != 0 {
		start += blockSize

		rawTables, err := iss.getRawTables(ctx, query, start)
		if err != nil {
			errors <- err
			return
		}

		err = json.Unmarshal(rawTables[query.table], &rawRows)
		if err != nil {
			errors <- err
			return
		}
		for _, rawRow := range rawRows {
			rows <- rawRow
		}

		if !query.multipart {
			return
		}

		blockSize = len(rawRows)

		curData, ok := rawTables["history.cursor"]
		if !ok {
			continue
		}

		err = json.Unmarshal(curData, &cur)
		if err != nil {
			errors <- err
			return
		}

		if start+cur[0].PageSize >= cur[0].Total {
			blockSize = 0
		}

	}

}

func (iss *ISSClient) getRawTables(ctx context.Context, query issQuery, start int) (rawTable map[string]json.RawMessage, err error) {

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
			rawTable = nil
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rawData := make([]json.RawMessage, 2)
	err = json.Unmarshal(body, &rawData)
	if err != nil {
		return nil, err
	}

	// Пропускаем первый элемент с информацией о кодировке
	err = json.Unmarshal(rawData[1], &rawTable)
	if err != nil {
		return nil, err
	}

	return
}

func (iss *ISSClient) getAll(ctx context.Context, query issQuery) (rows chan json.RawMessage, errors chan error) {
	rows = make(chan json.RawMessage)
	errors = make(chan error, 1)

	go iss.get(ctx, query, rows, errors)

	return
}
