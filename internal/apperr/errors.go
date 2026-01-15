package apperr

import "net/http"

func BadRequest(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusBadRequest, code, msg, err)
}

func Unauthorized(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusUnauthorized, code, msg, err)
}

func PaymentRequired(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusPaymentRequired, code, msg, err)
}

func Forbidden(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusForbidden, code, msg, err)
}

func NotFound(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusNotFound, code, msg, err)
}

func MethodNotAllowed(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusMethodNotAllowed, code, msg, err)
}

func NotAcceptable(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusNotAcceptable, code, msg, err)
}

func RequestTimeout(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusRequestTimeout, code, msg, err)
}

func Conflict(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusConflict, code, msg, err)
}

func Gone(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusGone, code, msg, err)
}

func LengthRequired(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusLengthRequired, code, msg, err)
}

func PreconditionFailed(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusPreconditionFailed, code, msg, err)
}

func RequestEntityTooLarge(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusRequestEntityTooLarge, code, msg, err)
}

func RequestURITooLong(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusRequestURITooLong, code, msg, err)
}

func UnsupportedMediaType(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusUnsupportedMediaType, code, msg, err)
}

func RequestedRangeNotSatisfiable(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusRequestedRangeNotSatisfiable, code, msg, err)
}

func UnprocessableEntity(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusUnprocessableEntity, code, msg, err)
}

func Locked(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusLocked, code, msg, err)
}

func TooManyRequests(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusTooManyRequests, code, msg, err)
}

func Internal(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusInternalServerError, code, msg, err)
}

func NotImplemented(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusNotImplemented, code, msg, err)
}

func BadGateway(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusBadGateway, code, msg, err)
}

func ServiceUnavailable(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusServiceUnavailable, code, msg, err)
}

func GatewayTimeout(code AppErrorCode, msg string, err error) *AppError {
	return New(http.StatusGatewayTimeout, code, msg, err)
}