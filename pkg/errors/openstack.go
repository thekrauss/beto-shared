package errors

// Keystone
func NewKeystoneAuthFailed(cause error) *Error {
	return Wrap(cause, CodeKeystoneAuthFailed, "keystone authentication failed")
}
func NewKeystoneTokenInvalid() *Error {
	return New(CodeKeystoneTokenInvalid, "invalid Keystone token")
}
func NewKeystoneForbidden() *Error {
	return New(CodeKeystoneForbidden, "forbidden Keystone request")
}

// Nova
func NewNovaError(cause error) *Error {
	return Wrap(cause, CodeNovaError, "nova API error")
}
func NewNovaNotFound(resource string) *Error {
	return Newf(CodeNovaNotFound, "nova resource not found: %s", resource)
}
func NewNovaQuotaExceeded(resource string) *Error {
	return Newf(CodeNovaQuotaExceed, "nova quota exceeded for %s", resource)
}

// Neutron
func NewNeutronError(cause error) *Error {
	return Wrap(cause, CodeNeutronError, "neutron API error")
}
func NewNeutronNotFound(resource string) *Error {
	return Newf(CodeNeutronNotFound, "neutron resource not found: %s", resource)
}
func NewNeutronIPConflict(ip string) *Error {
	return Newf(CodeNeutronIPConflict, "neutron IP conflict: %s", ip)
}
