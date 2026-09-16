package handler

import (
	"net/http"

	"calculator/internal/response"
)

// Sqrt handles POST /api/v1/sqrt: computes the square root of value.
func (h *Handlers) Sqrt(w http.ResponseWriter, r *http.Request) {
	var req sqrtRequest
	if verr := decodeJSON(r, &req); verr != nil {
		writeValidationError(w, verr)
		return
	}

	value, verr := requireField(req.Value, "value")
	if verr != nil {
		writeValidationError(w, verr)
		return
	}

	result, err := h.calc.Sqrt(value)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.WriteSuccess(w, "sqrt", result)
}
