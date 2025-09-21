package errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPCError convertit une erreur en gRPC status.Error
func ToGRPCError(err error) error {
	var e *Error
	if errors.As(err, &e) {
		return status.Error(mapCodeToGRPCStatus(e.Code), e.Message)
	}
	return status.Error(codes.Internal, "Internal server error")
}

// mapCodeToGRPCStatus mappe nos codes vers des status gRPC
func mapCodeToGRPCStatus(code string) codes.Code {
	switch code {
	case CodeUnauthorized, CodeKeystoneUnauthorized:
		return codes.Unauthenticated
	case CodeForbidden:
		return codes.PermissionDenied
	case CodeNotFound, CodeDBNotFound:
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
