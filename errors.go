package gomoex

import (
	"errors"
	"fmt"
)

// ErrISSClient - базовая ошибка при работе с MOEX ISS.
var ErrISSClient = errors.New("iss client error")

func warpErrWithMsg(msg string, err error) error {
	return fmt.Errorf("%w: %s - %s", ErrISSClient, msg, err)
}

func wrapParseErr(err error) error {
	return warpErrWithMsg("can't parse", err)
}
