// Package simplerr implement error with code
package simplerr

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type Call struct {
	Line     int
	File     string
	FuncName string
}

type withCode struct {
	err   error
	msg   string
	code  int64
	stack []Call
}

func (e *withCode) Unwrap() error {
	return e.err
}

func (e *withCode) Error() string {
	if e.err == nil {
		return ""
	}

	return strings.Join([]string{e.msg, e.err.Error()}, ": ")
}

func GetStack(err error) []Call {
	if e, ok := err.(*withCode); ok {
		return e.stack
	}

	return nil
}

func GetCode(err error) int64 {
	if e, ok := err.(*withCode); ok {
		return e.code
	}

	return 0
}

func GetText(err error) string {
	if e, ok := err.(*withCode); ok {
		return e.msg
	}

	return ""
}

func WrapWithCode(err error, code int64, msg string) error {
	return &withCode{
		err:   err,
		msg:   msg,
		code:  code,
		stack: makeStack(),
	}
}

func WrapfWithCode(err error, code int64, tmpl string, args ...interface{}) error {
	return &withCode{
		err:   err,
		msg:   fmt.Sprintf(tmpl, args...),
		code:  code,
		stack: makeStack(),
	}
}

func Wrap(err error, msg string) error {
	return WrapWithCode(err, 0, msg)
}

func Wrapf(err error, tmpl string, args ...interface{}) error {
	return WrapfWithCode(err, 0, tmpl, args...)
}

func Is(err, target error) bool {
	if e, ok := target.(interface{ Unwrap() error }); ok {
		return errors.Is(err, e.Unwrap())
	}

	return errors.Is(err, target)
}

func HasCode(err error, code int64) bool {
	for {
		if e, ok := err.(*withCode); ok {
			if e.code == code {
				return true
			}
		}

		if err = errors.Unwrap(err); err == nil {
			break
		}
	}

	return false
}

func makeStack() []Call {
	stack := make([]Call, 0)

	for i := 2; i < 34; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fname := "unknown"

		if fn := runtime.FuncForPC(pc); fn != nil {
			fname = fn.Name()
		}

		stack = append(stack, Call{
			Line:     line,
			File:     file,
			FuncName: fname,
		})
	}

	return stack
}
