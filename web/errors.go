package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error  string            `json:"error"`
	Errors map[string]string `json:"errors,omitempty"`
}

// FieldError is an error that is returned when a field fails validation
type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

type FieldErrors []FieldError

func NewFieldError(field string, err error) FieldErrors {
	var fe FieldErrors
	if err != nil {
		fe = append(fe, FieldError{field, err.Error()})
	}

	return fe
}

func (fe FieldErrors) Error() string {
	return "Validation Error"
}

func (fe FieldErrors) Fields() map[string]string {
	m := make(map[string]string)
	for _, fld := range fe {
		m[fld.Field] = fld.Err
	}
	return m
}

func IsFieldErrors(err error) bool {
	_, ok := err.(FieldErrors)
	return ok
}

func GetFieldErrors(err error) FieldErrors {
	var fe FieldErrors
	if !errors.As(err, &fe) {
		return FieldErrors{}
	}
	return fe
}

// trustedError is an error that is trusted to be returned to the client
type trustedError struct {
	Err    error
	Status int
}

func NewTrustedError(err error, status int) error {
	if err != nil {
		return &trustedError{
			err,
			status,
		}
	}
	return &trustedError{err, status}
}

func (te *trustedError) Error() string {
	return te.Err.Error()
}

func IsTrustedError(err error) bool {
	var te *trustedError
	return errors.As(err, &te)
}

func GetTrustedError(err error) *trustedError {
	var te *trustedError
	if !errors.As(err, &te) {
		return nil
	}
	return te
}

func IsEchoHTTPError(err error) bool {
	_, ok := err.(*echo.HTTPError)
	return ok
}

func GetEchoHTTPError(err error) *echo.HTTPError {
	if e, ok := err.(*echo.HTTPError); ok {
		return e
	}
	return &echo.HTTPError{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}
}

func HttpErrorHandler(err error, ctx echo.Context) {

	var er ErrorResponse
	var status int

	switch {
	case IsFieldErrors(err):
		er = ErrorResponse{
			Error:  "Invalid request",
			Errors: GetFieldErrors(err).Fields(),
		}
		status = http.StatusBadRequest

	case IsTrustedError(err):
		status = GetTrustedError(err).Status
		er = ErrorResponse{
			Error: GetTrustedError(err).Error(),
		}

	case IsEchoHTTPError(err):
		status = GetEchoHTTPError(err).Code

		message := GetEchoHTTPError(err).Message

		msg, ok := message.(string)

		if ok {
			er = ErrorResponse{
				Error: msg,
			}
		} else {
			er = ErrorResponse{
				Error: http.StatusText(status),
			}
		}
	default:
		isDevEnvironment := GetEnv("ENV", "development") == "development"
		status = http.StatusInternalServerError
		if isDevEnvironment {
			er = ErrorResponse{
				Error: err.Error(),
			}
		} else {
			er = ErrorResponse{
				Error: http.StatusText(status),
			}
		}
	}

	json.NewEncoder(os.Stdout).Encode(er)

	ctx.JSON(status, er)
}

type ErrorNotFound struct{}

func (e ErrorNotFound) Error() string {
	return "not found"
}
