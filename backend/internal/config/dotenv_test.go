package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvSetsUnsetVariables(t *testing.T) {
	t.Setenv("DOTENV_TEST_PORT", "")
	os.Unsetenv("DOTENV_TEST_PORT")

	path := writeTempEnvFile(t, "DOTENV_TEST_PORT=9999\n")
	LoadDotEnv(path)

	if got := os.Getenv("DOTENV_TEST_PORT"); got != "9999" {
		t.Errorf("expected DOTENV_TEST_PORT=9999, got %q", got)
	}
}

func TestLoadDotEnvDoesNotOverrideRealEnvironment(t *testing.T) {
	t.Setenv("DOTENV_TEST_ORIGIN", "http://real-value")

	path := writeTempEnvFile(t, "DOTENV_TEST_ORIGIN=http://from-file\n")
	LoadDotEnv(path)

	if got := os.Getenv("DOTENV_TEST_ORIGIN"); got != "http://real-value" {
		t.Errorf("expected real environment value to win, got %q", got)
	}
}

func TestLoadDotEnvSkipsBlankLinesAndComments(t *testing.T) {
	t.Setenv("DOTENV_TEST_PRECISION", "")
	os.Unsetenv("DOTENV_TEST_PRECISION")

	path := writeTempEnvFile(t, "# a comment\n\nDOTENV_TEST_PRECISION=4\n")
	LoadDotEnv(path)

	if got := os.Getenv("DOTENV_TEST_PRECISION"); got != "4" {
		t.Errorf("expected DOTENV_TEST_PRECISION=4, got %q", got)
	}
}

func TestLoadDotEnvIgnoresMissingFile(t *testing.T) {
	// Should not panic or error when the file does not exist.
	LoadDotEnv(filepath.Join(t.TempDir(), "does-not-exist.env"))
}

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp .env file: %v", err)
	}
	return path
}
