package simplerr

import "errors"

func New(msg string) error {
	return &withCode{
		err: errors.New(msg),
	}
}

func NewWithCode(msg string, code ErrCode) error {
	return &withCode{
		err:  errors.New(msg),
		code: code,
	}
}
