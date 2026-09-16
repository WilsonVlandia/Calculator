// Package config loads runtime settings from environment variables,
// applying documented defaults whenever a variable is missing or
// invalid.
package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	defaultServerPort         = "8080"
	defaultCORSAllowedOrigins = "*"
	defaultPrecision          = 2
)

// Config holds all runtime settings the API needs.
type Config struct {
	// ServerPort is the TCP port the HTTP server listens on.
	ServerPort string
	// CORSAllowedOrigins is the list of origins allowed by the global
	// CORS middleware.
	CORSAllowedOrigins []string
	// Precision is the number of decimal places every calculation
	// result is rounded to.
	Precision int
}

// Load reads SERVER_PORT, CORS_ALLOWED_ORIGINS and CALC_PRECISION from
// the environment, falling back to defaults for anything missing or
// invalid.
func Load() Config {
	return Config{
		ServerPort:         loadServerPort(),
		CORSAllowedOrigins: loadCORSAllowedOrigins(),
		Precision:          loadPrecision(),
	}
}

func loadServerPort() string {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		return port
	}
	return defaultServerPort
}

func loadCORSAllowedOrigins() []string {
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")
	if raw == "" {
		raw = defaultCORSAllowedOrigins
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, origin := range parts {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	if len(origins) == 0 {
		origins = []string{defaultCORSAllowedOrigins}
	}
	return origins
}

func loadPrecision() int {
	raw := os.Getenv("CALC_PRECISION")
	if raw == "" {
		return defaultPrecision
	}

	precision, err := strconv.Atoi(raw)
	if err != nil || precision < 0 {
		log.Printf("config: invalid CALC_PRECISION %q, falling back to default %d", raw, defaultPrecision)
		return defaultPrecision
	}
	return precision
}
