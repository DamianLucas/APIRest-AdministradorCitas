package utils

type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequest(msg string) *AppError {
	return &AppError{Status: 400, Message: msg}
}

func NewUnauthorized(msg string) *AppError {
	return &AppError{Status: 401, Message: msg}
}

func NewForbidden(msg string) *AppError {
	return &AppError{Status: 403, Message: msg}
}

func NewNotFound(msg string) *AppError {
	return &AppError{Status: 404, Message: msg}
}

func NewConflict(msg string) *AppError {
	return &AppError{Status: 409, Message: msg}
}

func NewInternal(msg string) *AppError {
	return &AppError{Status: 500, Message: msg}
}
