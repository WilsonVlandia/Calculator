// Command api starts the calculator HTTP server. It only wires the
// pieces together: it has no business logic of its own.
package main

import (
	"log"
	"net/http"

	"calculator/internal/calculator"
	"calculator/internal/config"
	"calculator/internal/handler"
	"calculator/internal/router"
)

func main() {
	config.LoadDotEnv(".env")
	cfg := config.Load()

	calc := calculator.New(cfg.Precision)
	handlers := handler.New(calc)
	mux := router.New(handlers, cfg.CORSAllowedOrigins)

	addr := ":" + cfg.ServerPort
	log.Printf("calculator API listening on %s (precision=%d, allowed origins=%v)", addr, cfg.Precision, cfg.CORSAllowedOrigins)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
