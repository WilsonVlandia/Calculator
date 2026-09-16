package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidationErrorMessage(t *testing.T) {
	verr := &validationError{code: "MISSING_FIELD", message: "field \"a\" is required", field: "a"}
	if verr.Error() != "field \"a\" is required" {
		t.Errorf("unexpected Error() result: %q", verr.Error())
	}
}

func TestWriteDomainErrorFallsBackToInternalError(t *testing.T) {
	rec := httptest.NewRecorder()
	writeDomainError(rec, errors.New("something unexpected"))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	body := decodeError(t, rec)
	if body.Error.Code != "INTERNAL_ERROR" {
		t.Errorf("unexpected error code: %s", body.Error.Code)
	}
}
