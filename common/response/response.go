// Package response Dictionaries and functions to build responses for HTTP requests
package response

import (
	"encoding/json"
	"errors"
	"go-boilerplate/common"
	"go-boilerplate/repository"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	unknownErrorCode        = "GEN001"
	unprocessableEntityCode = "GEN002"
	unauthorizedErrorCode   = "GEN003"
	forbiddenErrorCode      = "GEN004"

	validationErrorCode = "VLD001"
)

var errorCodes = [...]errorCode{}

type errorCode struct {
	err    error
	code   string
	status int
}

// Error is the default API error format
type Error struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

// Success is the default API success format
type Success struct {
	ID int `json:"id,omitempty"`
}

// Write writes needed headers and content to response
func Write(w http.ResponseWriter, body interface{}, status int) {
	if body == nil {
		w.WriteHeader(status)
		return
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		common.HandleError("error marshaling json body", err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(bytes)
}

// WriteServerError writes the given error to response
func WriteServerError(w http.ResponseWriter, err error, message string) {
	common.HandleError(message, err)
	Write(w, Error{
		Code:  unknownErrorCode,
		Error: err.Error(),
	}, http.StatusInternalServerError)
}

// WriteError writes the given error to response
func WriteError(w http.ResponseWriter, r *http.Request, err error, message string) {
	if !errors.Is(err, repository.ErrNotFound) {
		if span, ok := tracer.SpanFromContext(r.Context()); ok {
			span.SetTag(ext.Error, err)
			span.SetTag("cause.error", err.Error())
			span.SetTag("cause.message", message)
		}
	}
	for _, errorCode := range errorCodes {
		if errors.Is(err, errorCode.err) {
			Write(w, Error{
				Code:  errorCode.code,
				Error: err.Error(),
			}, errorCode.status)
			return
		}
	}
	WriteServerError(w, err, message)
}

// WriteUnauthorizedError writes the given error to response
func WriteUnauthorizedError(w http.ResponseWriter) {
	Write(w, Error{
		Code:  unauthorizedErrorCode,
		Error: http.StatusText(http.StatusUnauthorized),
	}, http.StatusUnauthorized)
}

// WriteForbiddenError writes the given error to response
func WriteForbiddenError(w http.ResponseWriter) {
	Write(w, Error{
		Code:  forbiddenErrorCode,
		Error: http.StatusText(http.StatusForbidden),
	}, http.StatusForbidden)
}

// WriteUnprocessableEntity writes a unprocessable entity response
func WriteUnprocessableEntity(w http.ResponseWriter, err error) {
	Write(w, Error{
		Code:  unprocessableEntityCode,
		Error: err.Error(),
	}, http.StatusUnprocessableEntity)
}

// WriteValidationError writes a vlidation error to response
func WriteValidationError(w http.ResponseWriter, err error) {
	Write(w, Error{
		Code:  validationErrorCode,
		Error: err.Error(),
	}, http.StatusBadRequest)
}
