package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// ApiError

type ApiError struct {
	statusCode int
	msg        string
	err        error
}

func NewApiError(statusCode int, msg string, err error) error {
	if err != nil {
		msg = fmt.Sprintf("%s: %s", msg, err.Error())
	}
	return &ApiError{statusCode: statusCode, msg: msg, err: err}
}

func (e *ApiError) Error() string {
	return e.msg
}

func (e *ApiError) Unwrap() error {
	return e.err
}

func (e *ApiError) StatusCode() int {
	return e.statusCode
}

func WrapApiError(template error, err error) error {
	apiErrorTemplate, ok := template.(*ApiError)
	if ok {
		err = NewApiError(apiErrorTemplate.StatusCode(), apiErrorTemplate.Error(), err)
	} else {
		err = WrapError(template.Error(), err)
	}
	return err
}

// http

var (
	ErrInvalidRequestBody = NewApiError(http.StatusBadRequest, "invalid request body", nil)
	ErrInvalidPrices      = NewApiError(http.StatusBadRequest, "invalid prices", nil)
)

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

func WriteError(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	errorResponse := ErrorResponse{StatusCode: statusCode, Error: errorMessage}
	err := encoder.Encode(errorResponse)
	if err != nil {
		panic(err)
	}
}

type HandleE func(http.ResponseWriter, *http.Request) error

func ErrAwareHandle(h HandleE) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			// respond error with simplified message
			errorMessage := GetErrorMessageTitle(err)
			var statusCode int
			switch e := err.(type) {
			case *ApiError:
				statusCode = e.StatusCode()
			default:
				statusCode = http.StatusInternalServerError
			}
			WriteError(w, errorMessage, statusCode)

			// log detailed error
			messageChain := ExtractErrorMessageChain(err)
			fmt.Fprintf(os.Stderr, "error:\n%s====\n", messageChain)
		}
	}
}
