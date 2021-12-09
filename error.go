package gomoex

import (
	"errors"
	"fmt"
)

// ErrISSClient - базовая ошибка при работе с MOEX ISS.
var ErrISSClient = errors.New("iss client error")

func newErrWithMsg(msg string) error {
	return fmt.Errorf("%w: %s", ErrISSClient, msg)
}

func newWarpedErr(msg string, err error) error {
	return fmt.Errorf("%w: %s - %s", ErrISSClient, msg, err)
}

func newParseErr(err error) error {
	return newWarpedErr("can't parse", err)
}
