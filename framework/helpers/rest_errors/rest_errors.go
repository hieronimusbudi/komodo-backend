package resterrors

import (
	"fmt"
	"net/http"
)

//RestErr interface
type RestErr interface {
	Message() string
	Status() int
	Error() string
	Causes() string
	ErrorResponse() restErr
}

type restErr struct {
	ErrMessage string `json:"message"`
	ErrStatus  int    `json:"status"`
	ErrError   string `json:"error"`
	ErrCauses  string `json:"causes"`
}

func (e restErr) ErrorResponse() restErr {
	return e
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: %v", e.ErrMessage, e.ErrStatus, e.ErrMessage, e.ErrCauses)
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) Status() int {
	return e.ErrStatus
}

func (e restErr) Causes() string {
	return e.ErrCauses
}

//NewRestError func
func NewRestError(message string, status int, err string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
	}
}

//NewBadRequestError func
func NewBadRequestError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "bad_request",
	}
}

//NewNotFoundError func
func NewNotFoundError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not_found",
	}
}

//NewUnauthorizedError func
func NewUnauthorizedError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   "unauthorized",
	}
}

//NewInternalServerError func
func NewInternalServerError(message string, err error) RestErr {
	result := restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "internal_server_error",
	}
	if err != nil {
		result.ErrCauses = err.Error()
	}
	return result
}
