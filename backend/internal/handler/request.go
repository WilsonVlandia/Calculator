// Package handler implements the controllers: it decodes and validates
// each request, calls the calculator model, and hands the outcome to
// the response package. It never builds a JSON body itself.
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"calculator/internal/calculator"
	"calculator/internal/response"
)

// Handlers holds the shared dependencies every endpoint needs.
type Handlers struct {
	calc *calculator.Calculator
}

// New creates the Handlers for the given calculator model.
func New(calc *calculator.Calculator) *Handlers {
	return &Handlers{calc: calc}
}

// binaryRequest is the payload shape for two-operand operations: add,
// subtract, multiply and divide.
type binaryRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

// powerRequest is the payload shape for the power operation.
type powerRequest struct {
	Base     *float64 `json:"base"`
	Exponent *float64 `json:"exponent"`
}

// sqrtRequest is the payload shape for the square root operation.
type sqrtRequest struct {
	Value *float64 `json:"value"`
}

// percentageRequest is the payload shape for the percentage operation.
type percentageRequest struct {
	Value      *float64 `json:"value"`
	Percentage *float64 `json:"percentage"`
}

// validationError describes a request that failed input validation,
// carrying enough detail for the caller to build a response.
type validationError struct {
	code    string
	message string
	field   string
}

func (e *validationError) Error() string {
	return e.message
}

// decodeJSON decodes the request body into dst, distinguishing
// malformed JSON from a field whose value is not a number. Missing
// fields are detected by the caller via requireField, since a nil
// pointer decodes successfully.
func decodeJSON(r *http.Request, dst any) *validationError {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) && typeErr.Field != "" {
			return &validationError{
				code:    response.CodeInvalidNumber,
				message: fmt.Sprintf("field %q must be a number", typeErr.Field),
				field:   typeErr.Field,
			}
		}
		return &validationError{
			code:    response.CodeInvalidJSON,
			message: "request body is not valid JSON",
		}
	}
	return nil
}

// requireField dereferences a decoded optional field, reporting a
// MISSING_FIELD validation error when it was absent from the payload.
func requireField(value *float64, name string) (float64, *validationError) {
	if value == nil {
		return 0, &validationError{
			code:    response.CodeMissingField,
			message: fmt.Sprintf("field %q is required", name),
			field:   name,
		}
	}
	return *value, nil
}

// writeValidationError reports a request validation failure.
func writeValidationError(w http.ResponseWriter, verr *validationError) {
	response.WriteError(w, verr.code, verr.message, verr.field)
}

// writeDomainError maps a calculator domain error to its API error
// code and reports it. Any error that is not one of the known domain
// sentinels is treated as an unexpected internal error.
func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calculator.ErrDivisionByZero):
		response.WriteError(w, response.CodeDivisionByZero, err.Error(), "")
	case errors.Is(err, calculator.ErrNegativeSqrtInput):
		response.WriteError(w, response.CodeNegativeSqrtInput, err.Error(), "")
	case errors.Is(err, calculator.ErrResultNotFinite):
		response.WriteError(w, response.CodeResultNotFinite, err.Error(), "")
	default:
		response.WriteError(w, response.CodeInternalError, "an unexpected error occurred", "")
	}
}
