# Gomoex

[![Go test](https://github.com/WLM1ke/gomoex/actions/workflows/test.yml/badge.svg)](https://github.com/WLM1ke/gomoex/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/WLM1ke/gomoex)](https://goreportcard.com/report/github.com/WLM1ke/gomoex)
[![codecov](https://codecov.io/gh/WLM1ke/gomoex/branch/main/graph/badge.svg?token=74YYC7H8EC)](https://codecov.io/gh/WLM1ke/gomoex)
[![Go Reference](https://pkg.go.dev/badge/github.com/WLM1ke/gomoex.svg)](https://pkg.go.dev/github.com/WLM1ke/gomoex)

Реализация части запросов к [MOEX Informational & Statistical Server](https://www.moex.com/a2193).

# Основные возможности

Реализованы запросы получения:

* списка торгуемых инструментов
* интервалов дат с доступными свечками и историческими котировками
* свечек и исторических котировок
* дивидендов

При необходимости перечень запросов может быть расширен. [Документация](https://pkg.go.dev/github.com/WLM1ke/gomoex).

# Пример использования

Получение дневных свечек для AKRN:

```
package main

import (
	"context"
	"fmt"
	"github.com/WLM1ke/gomoex"
	"net/http"
)

func main() {
	cl := gomoex.NewISSClient(http.DefaultClient)
	rows, _ := cl.MarketCandles(context.Background(), gomoex.EngineStock, gomoex.MarketShares, "AKRN", "2021-03-01", "2021-03-11", gomoex.IntervalDay)
	for _, row := range rows {
		fmt.Printf("%+v\n", row)
	}
}
```

Результат:

```
{Begin:2021-03-01 00:00:00 +0000 UTC End:2021-03-01 23:59:59 +0000 UTC Open:6006 Close:5992 High:6018 Low:5990 Value:5.138208e+06 Volume:856}
{Begin:2021-03-02 00:00:00 +0000 UTC End:2021-03-02 23:59:59 +0000 UTC Open:6006 Close:6032 High:6046 Low:5990 Value:1.2557102e+07 Volume:2087}
{Begin:2021-03-03 00:00:00 +0000 UTC End:2021-03-03 23:59:59 +0000 UTC Open:6048 Close:6000 High:6048 Low:5990 Value:7.280306e+06 Volume:1209}
{Begin:2021-03-04 00:00:00 +0000 UTC End:2021-03-04 23:59:59 +0000 UTC Open:6000 Close:5982 High:6008 Low:5964 Value:8.168796e+06 Volume:1365}
{Begin:2021-03-05 00:00:00 +0000 UTC End:2021-03-05 23:59:59 +0000 UTC Open:5968 Close:5996 High:6010 Low:5968 Value:4.505082e+06 Volume:752}
{Begin:2021-03-09 00:00:00 +0000 UTC End:2021-03-09 23:59:59 +0000 UTC Open:6018 Close:6010 High:6018 Low:5960 Value:9.577078e+06 Volume:1597}
{Begin:2021-03-10 00:00:00 +0000 UTC End:2021-03-10 23:59:59 +0000 UTC Open:6008 Close:6004 High:6010 Low:5982 Value:5.505522e+06 Volume:918}
{Begin:2021-03-11 00:00:00 +0000 UTC End:2021-03-11 23:59:59 +0000 UTC Open:6006 Close:6000 High:6010 Low:5992 Value:3.228186e+06 Volume:538}
```