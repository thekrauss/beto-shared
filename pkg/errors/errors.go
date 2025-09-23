package errors

import (
	sterrors "errors"
	"fmt"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"` //internal technical error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// creates a new simple error
func New(code string, message string) *Error {
	return &Error{Code: code, Message: message}
}

// creates a new error with formatted message
func Newf(code string, format string, args ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}

// coats an existing error
func Wrap(cause error, code string, message string) *Error {
	return &Error{Code: code, Message: message, Cause: cause}
}

// allows you to compare an error code
func Is(err error, code string) bool {
	var e *Error
	if sterrors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// compatibility with errors.Is /errors.As
func (e *Error) Unwrap() error {
	return e.Cause
}
