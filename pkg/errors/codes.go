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
	CodeDBError    = "DB_ERROR"
	CodeDBNotFound = "DB_NOT_FOUND"
	CodeDBConflict = "DB_CONFLICT"
)

// Codes OpenStack - Keystone
const (
	CodeKeystoneAuthFailed   = "KEYSTONE_AUTH_FAILED"
	CodeKeystoneTokenInvalid = "KEYSTONE_TOKEN_INVALID"
	CodeKeystoneForbidden    = "KEYSTONE_FORBIDDEN"
)

// Codes OpenStack - Nova
const (
	CodeNovaError       = "NOVA_ERROR"
	CodeNovaNotFound    = "NOVA_NOT_FOUND"
	CodeNovaQuotaExceed = "NOVA_QUOTA_EXCEEDED"
)

// Codes OpenStack - Neutron
const (
	CodeNeutronError      = "NEUTRON_ERROR"
	CodeNeutronNotFound   = "NEUTRON_NOT_FOUND"
	CodeNeutronIPConflict = "NEUTRON_IP_CONFLICT"
)
