package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteSuccess(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteSuccess(rec, "add", 5)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json content type, got %q", ct)
	}

	var body successBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if !body.Success || body.Operation != "add" || body.Result != 5 {
		t.Errorf("unexpected body: %+v", body)
	}
}

func TestWriteErrorStatusMapping(t *testing.T) {
	tests := []struct {
		code           string
		expectedStatus int
	}{
		{CodeInvalidJSON, http.StatusBadRequest},
		{CodeMissingField, http.StatusBadRequest},
		{CodeInvalidNumber, http.StatusBadRequest},
		{CodeDivisionByZero, http.StatusUnprocessableEntity},
		{CodeNegativeSqrtInput, http.StatusUnprocessableEntity},
		{CodeResultNotFinite, http.StatusUnprocessableEntity},
		{CodeNotFound, http.StatusNotFound},
		{CodeMethodNotAllowed, http.StatusMethodNotAllowed},
		{CodeInternalError, http.StatusInternalServerError},
		{"SOME_UNKNOWN_CODE", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			rec := httptest.NewRecorder()
			WriteError(rec, tt.code, "message", "field")

			if rec.Code != tt.expectedStatus {
				t.Errorf("code %s: expected status %d, got %d", tt.code, tt.expectedStatus, rec.Code)
			}

			var body errorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("failed to decode body: %v", err)
			}
			if body.Success {
				t.Errorf("expected success=false, got true")
			}
			if body.Error.Code != tt.code || body.Error.Message != "message" || body.Error.Field != "field" {
				t.Errorf("unexpected error body: %+v", body)
			}
		})
	}
}
