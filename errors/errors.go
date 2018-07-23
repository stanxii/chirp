package errors

import (
	"errors"
	"log"
	"net/http"
)

// type validationError struct {
// 	Field string `json:"field"`
// 	Error string `json:"error"`
// }}

// InternalServerError creates a new API error representing an internal server error (HTTP 500)
func InternalServerError(err error) *APIError {
	return NewAPIError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", Params{"error": err.Error()})
}

// NotFound creates a new API error representing a resource-not-found error (HTTP 404)
func NotFound(resource string) *APIError {
	return NewAPIError(http.StatusNotFound, "NOT_FOUND", Params{"resource": resource})
}

// Unauthorized creates a new API error representing an authentication failure (HTTP 401)
func Unauthorized(msg ...string) *APIError {
	var errMsg string
	if len(msg) > 0 {
		errMsg = msg[0]
	} else {
		//default msg
		errMsg = "You are not authorized to modify this resource"
	}
	return NewAPIError(http.StatusUnauthorized, "UNAUTHORIZED", Params{"error": errMsg})
}

// InvalidData converts a data validation error into an API error (HTTP 400)
func InvalidData(err error) *APIError {
	return NewAPIError(http.StatusBadRequest, "INVALID_DATA", Params{"message": err.Error()})
}

// GeneralErrorMsg is displayed when any random error
// is encountered by our backend.
const GeneralErrorMsg = "Something went wrong. Please try again, and contact us if the problem persists."

type PublicError interface {
	error
	Public() string
}

func SetCustomError(err error, data ...interface{}) *APIError {
	var apiErr *APIError
	if pErr, ok := err.(PublicError); ok {
		apiErr = NewAPIError(http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", Params{"message": pErr.Public()})
	} else {
		log.Println(err)
		apiErr = InvalidData(err)
	}

	if len(data) >= 1 {
		apiErr.Details = data[0]
	}
	if len(data) >= 2 {
		apiErr.Message = data[1].(string)
	}
	return apiErr
}

func NewStdError() error {
	return errors.New("")
}

// InvalidData converts a data validation error into an API error (HTTP 400)
// func InvalidData(errs validation.Errors) *APIError {
// 	result := []validationError{}}
// 	fields := []string{}
// 	for field := range errs {
// 		fields = append(fields, field)
// 	}
// 	sort.Strings(fields)
// 	for _, field := range fields {
// 		err := errs[field]
// 		result = append(result, validationError{
// 			Field: field,
// 			Error: err.Error(),
// 		})
// 	}

// 	err := NewAPIError(http.StatusBadRequest, "INVALID_DATA", nil)
// 	err.Details = result

// 	return err
// }
