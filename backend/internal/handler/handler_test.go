package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"calculator/internal/calculator"
)

type errorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field"`
	} `json:"error"`
}

type successResponse struct {
	Success   bool    `json:"success"`
	Operation string  `json:"operation"`
	Result    float64 `json:"result"`
}

func newTestHandlers() *Handlers {
	return New(calculator.New(2))
}

func doRequest(t *testing.T, handlerFunc http.HandlerFunc, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	handlerFunc(rec, req)
	return rec
}

func decodeSuccess(t *testing.T, rec *httptest.ResponseRecorder) successResponse {
	t.Helper()
	var body successResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode success response: %v (body: %s)", err, rec.Body.String())
	}
	return body
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) errorResponse {
	t.Helper()
	var body errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode error response: %v (body: %s)", err, rec.Body.String())
	}
	return body
}

func TestAddHandler(t *testing.T) {
	h := newTestHandlers()

	t.Run("success", func(t *testing.T) {
		rec := doRequest(t, h.Add, `{"a": 2, "b": 3}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		body := decodeSuccess(t, rec)
		if !body.Success || body.Operation != "add" || body.Result != 5 {
			t.Errorf("unexpected body: %+v", body)
		}
	})

	t.Run("missing field", func(t *testing.T) {
		rec := doRequest(t, h.Add, `{"a": 2}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "MISSING_FIELD" || body.Error.Field != "b" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("non-numeric field", func(t *testing.T) {
		rec := doRequest(t, h.Add, `{"a": "x", "b": 3}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_NUMBER" || body.Error.Field != "a" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		rec := doRequest(t, h.Add, `{"a": 2, "b":`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_JSON" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("unknown field is rejected as malformed", func(t *testing.T) {
		rec := doRequest(t, h.Add, `{"a": 2, "b": 3, "c": 4}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_JSON" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})
}

func TestSubtractHandler(t *testing.T) {
	h := newTestHandlers()

	t.Run("success", func(t *testing.T) {
		rec := doRequest(t, h.Subtract, `{"a": 5, "b": 3}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		body := decodeSuccess(t, rec)
		if body.Operation != "subtract" || body.Result != 2 {
			t.Errorf("unexpected body: %+v", body)
		}
	})

	t.Run("missing field", func(t *testing.T) {
		rec := doRequest(t, h.Subtract, `{"a": 5}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "MISSING_FIELD" || body.Error.Field != "b" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("non-numeric field", func(t *testing.T) {
		rec := doRequest(t, h.Subtract, `{"a": 5, "b": "x"}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_NUMBER" || body.Error.Field != "b" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		rec := doRequest(t, h.Subtract, `{"a":`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_JSON" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})
}

func TestMultiplyHandler(t *testing.T) {
	h := newTestHandlers()

	t.Run("success", func(t *testing.T) {
		rec := doRequest(t, h.Multiply, `{"a": 4, "b": 5}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		body := decodeSuccess(t, rec)
		if body.Operation != "multiply" || body.Result != 20 {
			t.Errorf("unexpected body: %+v", body)
		}
	})

	t.Run("overflow maps to RESULT_NOT_FINITE", func(t *testing.T) {
		rec := doRequest(t, h.Multiply, `{"a": 1.7976931348623157e+308, "b": 2}`)
		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "RESULT_NOT_FINITE" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("missing field", func(t *testing.T) {
		rec := doRequest(t, h.Multiply, `{"a": 4}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "MISSING_FIELD" || body.Error.Field != "b" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("non-numeric field", func(t *testing.T) {
		rec := doRequest(t, h.Multiply, `{"a": "x", "b": 4}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_NUMBER" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})
}

func TestDivideHandler(t *testing.T) {
	h := newTestHandlers()

	t.Run("success", func(t *testing.T) {
		rec := doRequest(t, h.Divide, `{"a": 10, "b": 2}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		body := decodeSuccess(t, rec)
		if body.Operation != "divide" || body.Result != 5 {
			t.Errorf("unexpected body: %+v", body)
		}
	})

	t.Run("division by zero maps to 422", func(t *testing.T) {
		rec := doRequest(t, h.Divide, `{"a": 10, "b": 0}`)
		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "DIVISION_BY_ZERO" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("missing field", func(t *testing.T) {
		rec := doRequest(t, h.Divide, `{"a": 10}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "MISSING_FIELD" || body.Error.Field != "b" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("non-numeric field", func(t *testing.T) {
		rec := doRequest(t, h.Divide, `{"a": 10, "b": "x"}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_NUMBER" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})
}

func TestPowerHandler(t *testing.T) {
	h := newTestHandlers()

	t.Run("success", func(t *testing.T) {
		rec := doRequest(t, h.Power, `{"base": 2, "exponent": 3}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		body := decodeSuccess(t, rec)
		if body.Operation != "power" || body.Result != 8 {
			t.Errorf("unexpected body: %+v", body)
		}
	})

	t.Run("zero base with negative exponent maps to DIVISION_BY_ZERO", func(t *testing.T) {
		rec := doRequest(t, h.Power, `{"base": 0, "exponent": -1}`)
		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "DIVISION_BY_ZERO" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("missing exponent", func(t *testing.T) {
		rec := doRequest(t, h.Power, `{"base": 2}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "MISSING_FIELD" || body.Error.Field != "exponent" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("non-numeric field", func(t *testing.T) {
		rec := doRequest(t, h.Power, `{"base": "x", "exponent": 2}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_NUMBER" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		rec := doRequest(t, h.Power, `{"base":`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_JSON" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})
}

func TestSqrtHandler(t *testing.T) {
	h := newTestHandlers()

	t.Run("success", func(t *testing.T) {
		rec := doRequest(t, h.Sqrt, `{"value": 4}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		body := decodeSuccess(t, rec)
		if body.Operation != "sqrt" || body.Result != 2 {
			t.Errorf("unexpected body: %+v", body)
		}
	})

	t.Run("negative input maps to 422", func(t *testing.T) {
		rec := doRequest(t, h.Sqrt, `{"value": -4}`)
		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "NEGATIVE_SQRT_INPUT" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("missing field", func(t *testing.T) {
		rec := doRequest(t, h.Sqrt, `{}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "MISSING_FIELD" || body.Error.Field != "value" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("non-numeric field", func(t *testing.T) {
		rec := doRequest(t, h.Sqrt, `{"value": "x"}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_NUMBER" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		rec := doRequest(t, h.Sqrt, `{"value":`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_JSON" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})
}

func TestPercentageHandler(t *testing.T) {
	h := newTestHandlers()

	t.Run("success", func(t *testing.T) {
		rec := doRequest(t, h.Percentage, `{"value": 200, "percentage": 15}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		body := decodeSuccess(t, rec)
		if body.Operation != "percentage" || body.Result != 30 {
			t.Errorf("unexpected body: %+v", body)
		}
	})

	t.Run("missing percentage field", func(t *testing.T) {
		rec := doRequest(t, h.Percentage, `{"value": 200}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "MISSING_FIELD" || body.Error.Field != "percentage" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("non-numeric field", func(t *testing.T) {
		rec := doRequest(t, h.Percentage, `{"value": 200, "percentage": "x"}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_NUMBER" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		rec := doRequest(t, h.Percentage, `{"value":`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		body := decodeError(t, rec)
		if body.Error.Code != "INVALID_JSON" {
			t.Errorf("unexpected error body: %+v", body)
		}
	})
}
