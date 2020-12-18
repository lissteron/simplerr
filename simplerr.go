// Package simplerr implement error with code
package simplerr

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type withCode struct {
	err  error
	msg  string
	code ErrCodeInterface
}

func (e *withCode) Unwrap() error {
	return e.err
}

func (e *withCode) Error() string {
	if e.err == nil {
		return ""
	}

	if e.msg == "" {
		return e.err.Error()
	}

	return strings.Join([]string{e.msg, e.err.Error()}, ": ")
}

func GetStack(err error) []Call {
	if e, ok := err.(*withStack); ok {
		return e.stack
	}

	return nil
}

func GetCode(err error) ErrCodeInterface {
	if e, ok := err.(*withCode); ok {
		return e.code
	}

	return new(ErrCode)
}

func GetText(err error) string {
	if e, ok := err.(*withCode); ok {
		return e.msg
	} else if err != nil {
		return err.Error()
	}

	return ""
}

func WithCode(err error, code ErrCodeInterface) error {
	return &withCode{
		err:  err,
		code: code,
	}
}

func WrapWithCode(err error, code ErrCodeInterface, msg string) error {
	return &withCode{
		err:  err,
		msg:  msg,
		code: code,
	}
}

func WrapfWithCode(err error, code ErrCodeInterface, tmpl string, args ...interface{}) error {
	return &withCode{
		err:  err,
		msg:  fmt.Sprintf(tmpl, args...),
		code: code,
	}
}

func Wrap(err error, msg string) error {
	return WrapWithCode(err, new(ErrCode), msg)
}

func Wrapf(err error, tmpl string, args ...interface{}) error {
	return WrapfWithCode(err, new(ErrCode), tmpl, args...)
}

func Is(err, target error) bool {
	if e, ok := target.(interface{ Unwrap() error }); ok {
		return errors.Is(err, e.Unwrap())
	}

	return errors.Is(err, target)
}

func HasCode(err error, code ErrCodeInterface) bool {
	for {
		if e, ok := err.(*withCode); ok {
			if e.code.Int() == code.Int() {
				return true
			}
		}

		if err = errors.Unwrap(err); err == nil {
			break
		}
	}

	return false
}

func GetWithCode(err error) error {
	for {
		if e, ok := err.(*withCode); ok {
			if e.code.Int() > 0 {
				return e
			}
		}

		if err = errors.Unwrap(err); err == nil {
			break
		}
	}

	return nil
}

type Call struct {
	Line     int
	File     string
	FuncName string
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	return &withStack{
		error: err,
		stack: makeStack(),
	}
}

type withStack struct {
	error
	stack []Call
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withStack) Unwrap() error { return w.error }

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
