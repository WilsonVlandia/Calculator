package handler

import (
	"net/http"

	"calculator/internal/response"
)

// Percentage handles POST /api/v1/percentage: computes percentage% of
// value (e.g. value=200, percentage=15 returns 30).
func (h *Handlers) Percentage(w http.ResponseWriter, r *http.Request) {
	var req percentageRequest
	if verr := decodeJSON(r, &req); verr != nil {
		writeValidationError(w, verr)
		return
	}

	value, verr := requireField(req.Value, "value")
	if verr != nil {
		writeValidationError(w, verr)
		return
	}
	percentage, verr := requireField(req.Percentage, "percentage")
	if verr != nil {
		writeValidationError(w, verr)
		return
	}

	result, err := h.calc.Percentage(value, percentage)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.WriteSuccess(w, "percentage", result)
}
