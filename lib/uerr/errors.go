// Package err define error type of uexky
package uerr

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type ErrorType string

const (
	UnknownError    ErrorType = ""
	ParamsError     ErrorType = "ParamsError"
	AuthError       ErrorType = "AuthError"
	PermissionError ErrorType = "PermissionError"
	NotFoundError   ErrorType = "NotFoundError"
	RateLimitError  ErrorType = "RateLimitError"
	InternalError   ErrorType = "InternalError"

	// External Services
	PostgresError ErrorType = "PostgresError"
	RedisError    ErrorType = "RedisError"
	MailgunError  ErrorType = "MailgunError"
)

var errCodes = map[ErrorType]string{
	ParamsError:     "INVALID_PARAMETER",
	AuthError:       "NOT_SIGNED_IN",
	PermissionError: "FORBIDDEN_ACTION",
	NotFoundError:   "NOT_FOUND",
	RateLimitError:  "RATE_LIMIT_EXCEEDED",

	UnknownError:  "INTERNAL_SREVER_ERROR",
	InternalError: "INTERNAL_SERVER_ERROR",
	PostgresError: "INTERNAL_SREVER_ERROR",
	RedisError:    "INTERNAL_SREVER_ERROR",
	MailgunError:  "INTERNAL_SREVER_ERROR",
}

func (t ErrorType) Code() string {
	code, ok := errCodes[t]
	if !ok {
		panic(fmt.Errorf("invalid error type: %s", t))
	}
	return code
}

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

func Wrap(t ErrorType, err error, a ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		t: t,
		e: errors.Wrap(err, fmt.Sprint(a...)),
	}
}

func Wrapf(t ErrorType, err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		t: t,
		e: errors.Wrapf(err, format, a...),
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

func (e *Error) As(target error) bool {
	_, ok := target.(*Error)
	return ok
}

func ExtractErrorType(err error) ErrorType {
	if err == nil {
		panic("error can't be nil")
	}
	for err != nil {
		if e, ok := err.(*Error); ok {
			return e.t
		}
		err = errors.Unwrap(err)
	}
	return ""
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
