package errors

// Keystone
func NewKeystoneAuthFailed(cause error) *Error {
	return Wrap(CodeKeystoneAuthFailed, "Authentication with Keystone failed", cause)
}
func NewKeystoneUnauthorized() *Error {
	return New(CodeKeystoneUnauthorized, "Unauthorized request to Keystone")
}

// Nova
func NewNovaQuotaExceeded(cause error) *Error {
	return Wrap(CodeNovaQuotaExceeded, "Nova quota exceeded", cause)
}
func NewNovaInstanceError(cause error) *Error {
	return Wrap(CodeNovaInstanceError, "Nova instance operation failed", cause)
}

// Neutron
func NewNeutronIPConflict(cause error) *Error {
	return Wrap(CodeNeutronIPConflict, "Neutron IP address conflict", cause)
}
func NewNeutronNetworkError(cause error) *Error {
	return Wrap(CodeNeutronNetworkError, "Neutron network operation failed", cause)
}
