package response

import "net/http"

// Error codes shared across the whole API. Callers pick one of these
// when reporting a failure; the HTTP status is derived from the code
// itself, so callers never have to remember which status goes with
// which case.
const (
	CodeInvalidJSON       = "INVALID_JSON"
	CodeMissingField      = "MISSING_FIELD"
	CodeInvalidNumber     = "INVALID_NUMBER"
	CodeDivisionByZero    = "DIVISION_BY_ZERO"
	CodeNegativeSqrtInput = "NEGATIVE_SQRT_INPUT"
	CodeResultNotFinite   = "RESULT_NOT_FINITE"
	CodeNotFound          = "NOT_FOUND"
	CodeMethodNotAllowed  = "METHOD_NOT_ALLOWED"
	CodeInternalError     = "INTERNAL_ERROR"
)

// statusByCode maps every known error code to the HTTP status it is
// reported with.
var statusByCode = map[string]int{
	CodeInvalidJSON:       http.StatusBadRequest,
	CodeMissingField:      http.StatusBadRequest,
	CodeInvalidNumber:     http.StatusBadRequest,
	CodeDivisionByZero:    http.StatusUnprocessableEntity,
	CodeNegativeSqrtInput: http.StatusUnprocessableEntity,
	CodeResultNotFinite:   http.StatusUnprocessableEntity,
	CodeNotFound:          http.StatusNotFound,
	CodeMethodNotAllowed:  http.StatusMethodNotAllowed,
	CodeInternalError:     http.StatusInternalServerError,
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type errorBody struct {
	Success bool        `json:"success"`
	Error   errorDetail `json:"error"`
}

// WriteError writes a JSON error response using the HTTP status
// associated with code. field is optional and identifies the request
// field the error is about, when applicable.
func WriteError(w http.ResponseWriter, code, message, field string) {
	status, ok := statusByCode[code]
	if !ok {
		status = http.StatusInternalServerError
	}

	writeJSON(w, status, errorBody{
		Success: false,
		Error: errorDetail{
			Code:    code,
			Message: message,
			Field:   field,
		},
	})
}
