package errors

import (
	sterrors "errors"
	"fmt"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"` // erreur technique interne
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// crée une nouvelle erreur simple
func New(code string, message string) *Error {
	return &Error{Code: code, Message: message}
}

// crée une nouvelle erreur avec message formaté
func Newf(code string, format string, args ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}

// enrobe une erreur existante
func Wrap(cause error, code string, message string) *Error {
	return &Error{Code: code, Message: message, Cause: cause}
}

// permet de comparer un code d’erreur
func Is(err error, code string) bool {
	var e *Error
	if sterrors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// compatibilité avec errors.Is / errors.As
func (e *Error) Unwrap() error {
	return e.Cause
}
