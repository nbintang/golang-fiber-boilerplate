package apperr

type AppError struct {
	Code    AppErrorCode `json:"code"`
	Message string       `json:"message"`
	Status  int          `json:"-"`
	Err     error        `json:"-"`
}

func (e *AppError) Error() string {
	if e.Message == "" {
		return e.Message
	}
	if e.Err == nil {
		return e.Err.Error()
	}
	return "error"
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(status int, code AppErrorCode, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Status: status, Err: err}
}
