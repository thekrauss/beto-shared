package errors

// Codes globaux
const (
	CodeInternal     = "INTERNAL_ERROR"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeNotFound     = "NOT_FOUND"
	CodeConflict     = "CONFLICT"
	CodeInvalidInput = "INVALID_INPUT"
	CodeTimeout      = "TIMEOUT"
)

// Codes DB
const (
	CodeDBNotFound = "DB_NOT_FOUND"
	CodeDBConflict = "DB_CONFLICT"
	CodeDBError    = "DB_ERROR"
)

// Codes OpenStack
const (
	CodeKeystoneAuthFailed   = "KEYSTONE_AUTH_FAILED"
	CodeKeystoneUnauthorized = "KEYSTONE_UNAUTHORIZED"

	CodeNovaQuotaExceeded = "NOVA_QUOTA_EXCEEDED"
	CodeNovaInstanceError = "NOVA_INSTANCE_ERROR"

	CodeNeutronIPConflict   = "NEUTRON_IP_CONFLICT"
	CodeNeutronNetworkError = "NEUTRON_NETWORK_ERROR"
)
