package errors

func NewDBNotFound(entity string) *Error {
	return Newf(CodeDBNotFound, "%s not found", entity)
}

func NewDBConflict(entity string) *Error {
	return Newf(CodeDBConflict, "conflict detected on %s", entity)
}

func NewDBError(cause error) *Error {
	return Wrap(cause, CodeDBError, "database error")
}
