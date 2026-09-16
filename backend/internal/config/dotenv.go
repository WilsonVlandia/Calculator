package config

import (
	"os"
	"strings"
)

// LoadDotEnv reads a .env file at path (one KEY=VALUE per line; blank
// lines and lines starting with "#" are ignored) and applies each
// entry to the process environment. A variable that is already set in
// the real environment is left untouched, so real environment
// variables always take precedence over the file. A missing file is
// not an error: Load's own defaults still apply.
func LoadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}
