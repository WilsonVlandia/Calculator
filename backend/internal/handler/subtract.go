package handler

import (
	"net/http"

	"calculator/internal/response"
)

// Subtract handles POST /api/v1/subtract: computes a - b.
func (h *Handlers) Subtract(w http.ResponseWriter, r *http.Request) {
	var req binaryRequest
	if verr := decodeJSON(r, &req); verr != nil {
		writeValidationError(w, verr)
		return
	}

	a, verr := requireField(req.A, "a")
	if verr != nil {
		writeValidationError(w, verr)
		return
	}
	b, verr := requireField(req.B, "b")
	if verr != nil {
		writeValidationError(w, verr)
		return
	}

	result, err := h.calc.Subtract(a, b)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.WriteSuccess(w, "subtract", result)
}
