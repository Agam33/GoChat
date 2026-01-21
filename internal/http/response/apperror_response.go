package response

import "net/http"

type AppErr struct {
	Code    int    // HTTP status code
	Message string // Client message
	Err     error  // Original err
}

func (e *AppErr) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Message
}

func NewAppErr(code int, msg string) *AppErr {
	return &AppErr{
		Code:    code,
		Message: msg,
	}
}

func WrapAppErr(code int, msg string, err error) *AppErr {
	return &AppErr{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func NewInternalServerErr(msg string, err error) *AppErr {
	return &AppErr{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Err:     err,
	}
}

func NewNotFoundErr(message string, err error) *AppErr {
	return &AppErr{
		Code:    http.StatusNotFound,
		Message: message,
		Err:     err,
	}
}

func NewBadRequestErr(message string, err error) *AppErr {
	return &AppErr{
		Code:    http.StatusBadRequest,
		Message: message,
		Err:     err,
	}
}

func NewAlreadyExistErr(message string, err error) *AppErr {
	return &AppErr{
		Code:    http.StatusConflict,
		Message: message,
		Err:     err,
	}
}

func NewUnauthorized() *AppErr {
	return &AppErr{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
	}
}
