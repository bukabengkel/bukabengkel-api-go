package utils

import (
	"errors"
	"fmt"
	"net/http"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	ErrBadRequest            = errors.New("bad request")
	ErrWrongCredentials      = errors.New("wrong credentials")
	ErrNotFound              = errors.New("not found")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrPermissionDenied      = errors.New("permission denied")
	ErrExpiredCSRFError      = errors.New("expired CSRF token")
	ErrWrongCSRFToken        = errors.New("wrong CSRF token")
	ErrCSRFNotPresented      = errors.New("CSRF not presented")
	ErrNotRequiredFields     = errors.New("no such required fields")
	ErrBadQueryParams        = errors.New("invalid query params")
	ErrInternalServerError   = errors.New("internal server error")
	ErrRequestTimeoutError   = errors.New("request timeout")
	ErrExistsEmailError      = errors.New("user with given email already exists")
	ErrInvalidJWTToken       = errors.New("invalid JWT token")
	ErrInvalidJWTClaims      = errors.New("invalid JWT claims")
	ErrNotAllowedImageHeader = errors.New("not allowed image header")
	ErrNoCookie              = errors.New("not found cookie header")
	ErrUnprocessableEntity   = errors.New("unprocessable entity")
	ErrAuthenticationFailed  = errors.New("authentication vailed")
)

// HttpErr interface
type HttpErr interface {
	Status() int
	Error() string
	Details() any
}

// HttpError struct
type HttpError struct {
	ErrStatus  int         `json:"status"`
	ErrError   string      `json:"error"`
	ErrDetails any `json:"details"`
}

// Error  Error() interface method
func (e HttpError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - details: %v", e.ErrStatus, e.ErrError, e.ErrDetails)
}

// Error status
func (e HttpError) Status() int {
	return e.ErrStatus
}

// HttpError Details
func (e HttpError) Details() any {
	return e.ErrDetails
}

// New Http Error
func NewHttpError(status int, err string, details any) HttpErr {
	return HttpError{
		ErrStatus:  status,
		ErrError:   err,
		ErrDetails: details,
	}
}

// New Authentication Failed Error
func NewAuthenticationFailedError(details any) HttpErr {
	return HttpError{
		ErrStatus:  401,
		ErrError:   ErrAuthenticationFailed.Error(),
		ErrDetails: details,
	}
}

// New Bad Request Error
func NewBadRequestError(details any) HttpErr {
	return HttpError{
		ErrStatus:  http.StatusBadRequest,
		ErrError:   ErrBadRequest.Error(),
		ErrDetails: details,
	}
}

// New Not Found Error
func NewNotFoundError(details any) HttpErr {
	return HttpError{
		ErrStatus:  http.StatusNotFound,
		ErrError:   ErrNotFound.Error(),
		ErrDetails: details,
	}
}

// New Unauthorized Error
func NewUnauthorizedError(details any) HttpErr {
	return HttpError{
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   ErrUnauthorized.Error(),
		ErrDetails: details,
	}
}

// New Forbidden Error
func NewForbiddenError(details any) HttpErr {
	return HttpError{
		ErrStatus:  http.StatusForbidden,
		ErrError:   ErrForbidden.Error(),
		ErrDetails: details,
	}
}

// New Internal Server Error
func NewInternalServerError(err error) HttpErr {
	log.Error(err.Error())

	return HttpError{
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   ErrInternalServerError.Error(),
		ErrDetails: err,
	}
}

// New Internal Server Error
func NewDomainError(details any) HttpErr {
	log.Error(details.(error).Error())

	return HttpError{
		ErrStatus:  http.StatusBadRequest,
		ErrError:   ErrBadRequest.Error(),
		ErrDetails: details,
	}
}

// New Unprocessable Entity Error
func NewUnprocessableEntityError(details any) HttpErr {
	return HttpError{
		ErrStatus:  http.StatusUnprocessableEntity,
		ErrError:   ErrUnprocessableEntity.Error(),
		ErrDetails: details,
	}
}

// New Invalid Input Error - Validation
func NewHTTPValidationError(c echo.Context, e validator.ValidationErrors) error {
	errorMessages := make(map[string]string)

	for _, err := range e {
		// Base Assumption the json tag always in snake case and represent the field name
		fieldName := toSnakeCase(err.StructField())

		errorMessages[fieldName] = fmt.Sprintf("Failed on the '%s' tag", err.Tag())
	}

	return c.JSON(http.StatusBadRequest, Error{
		Code:    http.StatusBadRequest,
		Message: "Validation Error",
		Error:   errorMessages,
	})
}

// Parse Http Error
func ParseHttpError(err error) (int, HttpErr) {
	if httpErr, ok := err.(HttpErr); ok {
		return httpErr.Status(), httpErr
	}
	return http.StatusInternalServerError, NewInternalServerError(err)
}

// PanicIfNeeded is panic if needed
func PanicIfNeeded(err any) {
	if err != nil {
		panic(err)
	}
}

func toSnakeCase(s string) string {
	var result string
	for i, r := range s {
			if i > 0 && unicode.IsUpper(r) {
					result += "_"
			}
			result += string(unicode.ToLower(r))
	}
	return result
}
