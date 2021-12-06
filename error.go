package gomoex

import "fmt"

const _errName = "iss client error"

// ISSClientError - базовая ошибка при работе с MOEX ISS.
type ISSClientError struct {
	msg string
	err error
}

func newErrWithMsg(msg string) error {
	return ISSClientError{
		msg: msg,
	}
}

func newWarpedErr(msg string, err error) error {
	return ISSClientError{
		msg: msg,
		err: err,
	}
}

func newParseErr(err error) error {
	return newWarpedErr("can't parse", err)
}

// Error возвращает текстовое представление ошибки.
func (iss ISSClientError) Error() string {
	if iss.err == nil {
		return fmt.Sprintf("%s: %s", _errName, iss.msg)
	}

	return fmt.Sprintf("%s: %s - %s", _errName, iss.msg, iss.err)
}

// Unwrap возвращает стороннюю причину возникновения ошибки при ее наличии.
func (iss ISSClientError) Unwrap() error {
	return iss.err
}
