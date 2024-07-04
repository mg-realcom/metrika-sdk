package metrika_sdk

import "fmt"

const (
	RequestFailedMsg           = "request failed"
	CreateRequestFailedMsg     = "can't create request"
	UnmarshalResponseFailedMsg = "can't unmarshal response"
	ReadResponseFailedMsg      = "can't read response body"
)

type APIError struct {
	Errors []struct {
		ErrorType string `json:"error_type"`
		Message   string `json:"message"`
	} `json:"errors"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("code: %v, message: %s", e.Code, e.Message)
}

type InternalError struct {
	err error
	msg string
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e *InternalError) Unwrap() error {
	return e.err
}

func newInternalError(err error, msg string) *InternalError {
	return &InternalError{err: err, msg: msg}
}
