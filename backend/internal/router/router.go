// Package router wires each calculator endpoint to its handler,
// enforces the expected HTTP method per route, and applies the global
// CORS policy.
package router

import (
	"net/http"

	"calculator/internal/handler"
	"calculator/internal/response"
)

// New builds the complete HTTP handler for the API: every calculator
// route, a JSON 404 fallback, and the CORS middleware wrapped around
// all of it.
func New(h *handler.Handlers, allowedOrigins []string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/add", withMethod(http.MethodPost, h.Add))
	mux.HandleFunc("/api/v1/subtract", withMethod(http.MethodPost, h.Subtract))
	mux.HandleFunc("/api/v1/multiply", withMethod(http.MethodPost, h.Multiply))
	mux.HandleFunc("/api/v1/divide", withMethod(http.MethodPost, h.Divide))
	mux.HandleFunc("/api/v1/power", withMethod(http.MethodPost, h.Power))
	mux.HandleFunc("/api/v1/sqrt", withMethod(http.MethodPost, h.Sqrt))
	mux.HandleFunc("/api/v1/percentage", withMethod(http.MethodPost, h.Percentage))
	mux.HandleFunc("/", notFound)

	return withCORS(mux, allowedOrigins)
}

// withMethod rejects requests that use the wrong HTTP method for a
// route, before the underlying handler ever runs.
func withMethod(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			response.WriteError(w, response.CodeMethodNotAllowed, "method not allowed for this route", "")
			return
		}
		next(w, r)
	}
}

// notFound handles any path that does not match a known route.
func notFound(w http.ResponseWriter, r *http.Request) {
	response.WriteError(w, response.CodeNotFound, "route not found", "")
}
