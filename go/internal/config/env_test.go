package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fenneh/reddit-stream-console/internal/config"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return f
}

func TestLoadDotEnvMissingFile(t *testing.T) {
	if err := config.LoadDotEnv("/nonexistent/.env"); err != nil {
		t.Errorf("missing file should return nil, got: %v", err)
	}
}

func TestLoadDotEnvSetsVars(t *testing.T) {
	key := "TEST_DOTENV_SETS_" + t.Name()
	os.Unsetenv(key)
	t.Cleanup(func() { os.Unsetenv(key) })

	f := writeTempEnv(t, key+"=hello\n")
	if err := config.LoadDotEnv(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv(key); got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
}

func TestLoadDotEnvDoesNotOverwriteExisting(t *testing.T) {
	key := "TEST_DOTENV_NOOVER_" + t.Name()
	os.Setenv(key, "original")
	t.Cleanup(func() { os.Unsetenv(key) })

	f := writeTempEnv(t, key+"=overwritten\n")
	if err := config.LoadDotEnv(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv(key); got != "original" {
		t.Errorf("existing var overwritten: got %q, want %q", got, "original")
	}
}

func TestLoadDotEnvStripsQuotes(t *testing.T) {
	key := "TEST_DOTENV_QUOTES_" + t.Name()
	os.Unsetenv(key)
	t.Cleanup(func() { os.Unsetenv(key) })

	f := writeTempEnv(t, key+`="quoted value"`+"\n")
	if err := config.LoadDotEnv(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv(key); got != "quoted value" {
		t.Errorf("got %q, want %q", got, "quoted value")
	}
}

func TestLoadDotEnvSkipsCommentsAndBlanks(t *testing.T) {
	key := "TEST_DOTENV_SKIP_" + t.Name()
	os.Unsetenv(key)
	t.Cleanup(func() { os.Unsetenv(key) })

	f := writeTempEnv(t, "# comment\n\n"+key+"=value\n# another\n")
	if err := config.LoadDotEnv(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv(key); got != "value" {
		t.Errorf("got %q, want %q", got, "value")
	}
}

func TestLoadDotEnvSkipsLineWithoutEquals(t *testing.T) {
	key := "TEST_DOTENV_NOEQUALS_" + t.Name()
	os.Unsetenv(key)
	t.Cleanup(func() { os.Unsetenv(key) })

	// "NOEQUALS" has no '=', should be silently skipped
	f := writeTempEnv(t, "JUSTAKEYNOEQUALS\n"+key+"=ok\n")
	if err := config.LoadDotEnv(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv(key); got != "ok" {
		t.Errorf("got %q, want %q", got, "ok")
	}
}
