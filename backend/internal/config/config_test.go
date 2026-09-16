package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("SERVER_PORT", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	t.Setenv("CALC_PRECISION", "")

	cfg := Load()

	if cfg.ServerPort != defaultServerPort {
		t.Errorf("expected default port %q, got %q", defaultServerPort, cfg.ServerPort)
	}
	if len(cfg.CORSAllowedOrigins) != 1 || cfg.CORSAllowedOrigins[0] != defaultCORSAllowedOrigins {
		t.Errorf("expected default origins [%q], got %v", defaultCORSAllowedOrigins, cfg.CORSAllowedOrigins)
	}
	if cfg.Precision != defaultPrecision {
		t.Errorf("expected default precision %d, got %d", defaultPrecision, cfg.Precision)
	}
}

func TestLoadReadsEnvironmentValues(t *testing.T) {
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://a.com, http://b.com")
	t.Setenv("CALC_PRECISION", "4")

	cfg := Load()

	if cfg.ServerPort != "9090" {
		t.Errorf("expected port 9090, got %q", cfg.ServerPort)
	}
	expectedOrigins := []string{"http://a.com", "http://b.com"}
	if len(cfg.CORSAllowedOrigins) != len(expectedOrigins) {
		t.Fatalf("expected %v, got %v", expectedOrigins, cfg.CORSAllowedOrigins)
	}
	for i, origin := range expectedOrigins {
		if cfg.CORSAllowedOrigins[i] != origin {
			t.Errorf("expected origin %q at index %d, got %q", origin, i, cfg.CORSAllowedOrigins[i])
		}
	}
	if cfg.Precision != 4 {
		t.Errorf("expected precision 4, got %d", cfg.Precision)
	}
}

func TestLoadFallsBackOnInvalidPrecision(t *testing.T) {
	t.Setenv("CALC_PRECISION", "not-a-number")

	cfg := Load()

	if cfg.Precision != defaultPrecision {
		t.Errorf("expected fallback to default precision %d, got %d", defaultPrecision, cfg.Precision)
	}
}

func TestLoadFallsBackOnNegativePrecision(t *testing.T) {
	t.Setenv("CALC_PRECISION", "-1")

	cfg := Load()

	if cfg.Precision != defaultPrecision {
		t.Errorf("expected fallback to default precision %d, got %d", defaultPrecision, cfg.Precision)
	}
}
