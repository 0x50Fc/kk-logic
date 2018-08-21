package logic

import (
	"fmt"
)

const (
	ERROR_UNKNOWN = 400
	ERROR_INPUT   = 401
	ERROR_OK      = 200
)

type Error struct {
	Errno  int
	Errmsg string
}

func (E *Error) Error() string {
	return fmt.Sprintf("[%d] %s", E.Errno, E.Errmsg)
}

func NewError(errno int, errmsg string) *Error {
	return &Error{errno, errmsg}
}

func GetError(err error) (int, string) {
	e, ok := err.(*Error)
	if ok {
		return e.Errno, e.Errmsg
	}
	return ERROR_UNKNOWN, err.Error()
}

func GetErrorObject(err error) interface{} {
	n, m := GetError(err)
	return map[string]interface{}{"errno": n, "errmsg": m}
}
