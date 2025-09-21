package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"` // erreur interne
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// crée une nouvelle erreur sans cause
func New(code, message string) *Error {
	return &Error{Code: code, Message: message}
}

// crée une nouvelle erreur formatée
func Newf(code, format string, args ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}

// enrobe une erreur existante
func Wrap(code, message string, cause error) *Error {
	return &Error{Code: code, Message: message, Cause: cause}
}

// Is permet de comparer les erreurs
func Is(err error, code string) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// pour compatibilité avec errors.Is / errors.As
func (e *Error) Unwrap() error {
	return e.Cause
}
