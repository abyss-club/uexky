// Package err define error type of uexky
package uerr

import (
	"errors"
	"fmt"
)

type ErrorType string

const (
	ParamsError     ErrorType = "ParamsError"
	AuthError       ErrorType = "AuthError"
	PermissionError ErrorType = "PermissionError"
	NotFoundError   ErrorType = "NotFoundError"
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
