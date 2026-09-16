package handler

import (
	"net/http"

	"calculator/internal/response"
)

// Power handles POST /api/v1/power: computes base raised to exponent.
func (h *Handlers) Power(w http.ResponseWriter, r *http.Request) {
	var req powerRequest
	if verr := decodeJSON(r, &req); verr != nil {
		writeValidationError(w, verr)
		return
	}

	base, verr := requireField(req.Base, "base")
	if verr != nil {
		writeValidationError(w, verr)
		return
	}
	exponent, verr := requireField(req.Exponent, "exponent")
	if verr != nil {
		writeValidationError(w, verr)
		return
	}

	result, err := h.calc.Power(base, exponent)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.WriteSuccess(w, "power", result)
}
