package errors

import (
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

type Error struct {
	t ErrorType
	e error
}

type ErrorType string

const (
	TypeUnknown    ErrorType = ""
	TypeParams     ErrorType = "ParamsError"
	TypeAuth       ErrorType = "AuthError"
	TypePermission ErrorType = "PermissionError"
	TypeNotFound   ErrorType = "NotFoundError"
	TypeComplexity ErrorType = "ComplexityLimitError"
	TypeDuplicated ErrorType = "DuplicatedError"
	TypeInternal   ErrorType = "InternalError"

	// External Services
	TypePostgres ErrorType = "PostgresError"
	TypeRedis    ErrorType = "RedisError"
	TypeMailgun  ErrorType = "MailgunError"
)

var errCodes = map[ErrorType]string{
	TypeParams:     "INVALID_PARAMETER",
	TypeAuth:       "NOT_SIGNED_IN",
	TypePermission: "FORBIDDEN_ACTION",
	TypeNotFound:   "NOT_FOUND",
	TypeComplexity: "COMPLEXITY_LIMIT_EXCEEDED",
	TypeDuplicated: "DUPLICATED_CONTENT",

	TypeUnknown:  "INTERNAL_SREVER_ERROR",
	TypeInternal: "INTERNAL_SERVER_ERROR",
	TypePostgres: "INTERNAL_SREVER_ERROR",
	TypeRedis:    "INTERNAL_SREVER_ERROR",
	TypeMailgun:  "INTERNAL_SREVER_ERROR",
}

var (
	BadParams  = &Error{t: TypeParams}
	NoAuth     = &Error{t: TypeAuth}
	Permission = &Error{t: TypePermission}
	NotFound   = &Error{t: TypeNotFound}
	Complexity = &Error{t: TypeComplexity}
	Duplicated = &Error{t: TypeDuplicated}
	Internal   = &Error{t: TypeInternal}
	Postgres   = &Error{t: TypePostgres}
	Redis      = &Error{t: TypeRedis}
	Mailgun    = &Error{t: TypeMailgun}
)

func (t ErrorType) Code() string {
	code, ok := errCodes[t]
	if !ok {
		panic(fmt.Errorf("invalid error type: %s", t))
	}
	return code
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.t, e.e.Error())
}

func (e *Error) Code() string {
	return e.t.Code()
}

func (e *Error) New(a ...interface{}) error {
	return &Error{
		t: e.t,
		e: pkgerrors.New(fmt.Sprint(a...)),
	}
}

func (e *Error) Errorf(format string, a ...interface{}) error {
	return &Error{
		t: e.t,
		e: pkgerrors.New(fmt.Sprintf(format, a...)),
	}
}

func (e *Error) Handle(err error, a ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		t: e.t,
		e: pkgerrors.Wrap(err, fmt.Sprint(a...)),
	}
}

func (e *Error) Handlef(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		t: e.t,
		e: pkgerrors.Wrapf(err, format, a...),
	}
}

func (e *Error) Unwrap() error { return e.e }

func (e *Error) Cause() error { return e.e }

func (e *Error) Is(target error) bool {
	te, ok := target.(*Error)
	return ok && te.t == e.t
}
