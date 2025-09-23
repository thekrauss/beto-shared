package errors

import (
	sterrors "errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// converts an error to gRPC status.Error
func ToGRPCError(err error) error {
	var e *Error
	if sterrors.As(err, &e) {
		return status.Error(mapCodeToGRPCStatus(e.Code), e.Message)
	}
	return status.Error(codes.Internal, "internal server error")
}

// maps our codes to gRPC statuses
func mapCodeToGRPCStatus(code string) codes.Code {
	switch code {
	case CodeUnauthorized, CodeKeystoneAuthFailed, CodeKeystoneTokenInvalid:
		return codes.Unauthenticated
	case CodeForbidden, CodeKeystoneForbidden:
		return codes.PermissionDenied
	case CodeNotFound, CodeDBNotFound, CodeNovaNotFound, CodeNeutronNotFound:
		return codes.NotFound
	case CodeConflict, CodeDBConflict, CodeNeutronIPConflict:
		return codes.AlreadyExists
	case CodeInvalidInput:
		return codes.InvalidArgument
	case CodeTimeout:
		return codes.DeadlineExceeded
	default:
		return codes.Internal
	}
}
