package errors

import (
	stderrors "errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

func Wrap(err error, a ...interface{}) error {
	return pkgerrors.Wrap(err, fmt.Sprint(a...))
}

func Wrapf(err error, format string, a ...interface{}) error {
	return pkgerrors.Wrapf(err, format, a...)
}

func Is(err, target error) bool { return stderrors.Is(err, target) }

func As(err error, target interface{}) bool { return stderrors.As(err, target) }

func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}
