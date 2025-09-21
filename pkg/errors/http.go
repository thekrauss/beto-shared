package errors

import (
	"encoding/json"
	sterrors "errors"
	"net/http"
)

// la réponse JSON standardisée
type HTTPErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// convertit une erreur en (statusCode, body JSON)
func ToHTTPError(err error) (int, []byte) {

	var e *Error
	if sterrors.As(err, &e) {
		status := mapCodeToHTTPStatus(e.Code)
		resp := HTTPErrorResponse{
			Code:    e.Code,
			Message: e.Message,
		}
		body, _ := json.Marshal(resp)
		return status, body
	}

	// fallback si erreur inconnue
	resp := HTTPErrorResponse{
		Code:    CodeInternal,
		Message: "Internal server error",
	}
	body, _ := json.Marshal(resp)
	return http.StatusInternalServerError, body
}

// mappe nos codes vers des status HTTP
func mapCodeToHTTPStatus(code string) int {
	switch code {
	case CodeUnauthorized, CodeKeystoneUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound, CodeDBNotFound:
		return http.StatusNotFound
	case CodeConflict, CodeDBConflict, CodeNeutronIPConflict:
		return http.StatusConflict
	case CodeInvalidInput:
		return http.StatusBadRequest
	case CodeTimeout:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}
