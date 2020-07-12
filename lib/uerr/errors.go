// Package err define error type of uexky
package uerr

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type ErrorType string

const (
	ParamsError     ErrorType = "ParamsError"
	AuthError       ErrorType = "AuthError"
	PermissionError ErrorType = "PermissionError"
	NotFoundError   ErrorType = "NotFoundError"
	DBError         ErrorType = "DBError"
	InternalError   ErrorType = "InternalError"
)

type Error struct {
	t ErrorType
	e error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.t, e.e.Error())
}

func New(t ErrorType, a ...interface{}) *Error {
	return &Error{
		t: t,
		e: errors.New(fmt.Sprint(a...)),
	}
}

func Errorf(t ErrorType, format string, a ...interface{}) *Error {
	return &Error{
		t: t,
		e: fmt.Errorf(format, a...),
	}
}

func (e *Error) Unwrap() error {
	return errors.Unwrap(e.e)
}

func (e *Error) Is(target error) bool {
	ue, ok := target.(*Error)
	if !ok {
		return false
	}
	return ue.t == e.t
}

type ErrSlice []error

func (e ErrSlice) Error() string {
	if len(e) == 0 {
		return "empty error"
	}
	if len(e) == 1 {
		return e[0].Error()
	}
	var msg []string
	for _, err := range e {
		msg = append(msg, err.Error())
	}
	return fmt.Sprintf("multiple errors: %s", strings.Join(msg, "; "))
}
