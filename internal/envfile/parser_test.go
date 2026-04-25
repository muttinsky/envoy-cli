package envfile

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := env.ToMap()
	if m["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", m["APP_ENV"])
	}
	if m["DEBUG"] != "false" {
		t.Errorf("expected DEBUG=false, got %q", m["DEBUG"])
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `DB_URL="postgres://localhost/mydb"
SECRET='mysecret'
`)
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := env.ToMap()
	if m["DB_URL"] != "postgres://localhost/mydb" {
		t.Errorf("unexpected DB_URL: %q", m["DB_URL"])
	}
	if m["SECRET"] != "mysecret" {
		t.Errorf("unexpected SECRET: %q", m["SECRET"])
	}
}

func TestParse_CommentsAndBlankLines(t *testing.T) {
	path := writeTempEnv(t, "# This is a comment\n\nKEY=value\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := env.ToMap()
	if _, ok := m[""]; ok {
		t.Error("map should not contain empty key from comment")
	}
	if m["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", m["KEY"])
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestStripQuotes(t *testing.T) {
	cases := []struct{ input, want string }{
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{`plain`, "plain"},
		{`"`, `"`},
	}
	for _, c := range cases {
		if got := stripQuotes(c.input); got != c.want {
			t.Errorf("stripQuotes(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
