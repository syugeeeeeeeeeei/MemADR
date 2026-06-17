package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInitCreatesMemoryDirs(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	var out bytes.Buffer

	if err := Run([]string{"init"}, dir, &out, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	want := []string{
		"memory",
		"memory/bugs",
		"memory/problems",
		"memory/decisions",
		"memory/changes",
		"memory/reversions",
		"memory/solutions",
		"memory/supersessions",
		"memory/generated",
	}

	for _, rel := range want {
		path := filepath.Join(dir, filepath.FromSlash(rel))
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("os.Stat(%q) error: %v", path, err)
		}
		if !info.IsDir() {
			t.Fatalf("%q is not a directory", path)
		}
	}
}

func TestRunNewBugCreatesFirstRecord(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	var out bytes.Buffer

	if err := Run([]string{"init"}, dir, &out, &out); err != nil {
		t.Fatalf("init error: %v", err)
	}

	out.Reset()
	if err := Run([]string{"new", "bug", "認証状態が壊れる"}, dir, &out, &out); err != nil {
		t.Fatalf("new bug error: %v", err)
	}

	path := filepath.Join(dir, "memory", "bugs", "BUG-001.md")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error: %v", path, err)
	}

	text := string(body)
	mustContain(t, text, "# BUG-001: 認証状態が壊れる")
	mustContain(t, text, "Status: OPEN")
	mustContain(t, text, "Area: ")
	mustContain(t, text, "Symptom: ")
}

func TestRunNewAdrUsesNextID(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	var out bytes.Buffer

	if err := Run([]string{"init"}, dir, &out, &out); err != nil {
		t.Fatalf("init error: %v", err)
	}

	existing := filepath.Join(dir, "memory", "decisions", "ADR-001.md")
	body := strings.Join([]string{
		"# ADR-001: 既存判断",
		"",
		"Status: ACCEPTED",
		"Context: 既存",
		"Decision: 維持",
		"Rejected: なし",
		"Consequence: 影響なし",
		"",
	}, "\n")
	if err := os.WriteFile(existing, []byte(body), 0o644); err != nil {
		t.Fatalf("os.WriteFile(%q) error: %v", existing, err)
	}

	out.Reset()
	if err := Run([]string{"new", "adr", "新しい判断"}, dir, &out, &out); err != nil {
		t.Fatalf("new adr error: %v", err)
	}

	path := filepath.Join(dir, "memory", "decisions", "ADR-002.md")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error: %v", path, err)
	}

	text := string(got)
	mustContain(t, text, "# ADR-002: 新しい判断")
	mustContain(t, text, "Status: PROPOSED")
	mustContain(t, text, "Context: ")
	mustContain(t, text, "Decision: ")
}

func mustContain(t *testing.T, got string, want string) {
	t.Helper()

	if !strings.Contains(got, want) {
		t.Fatalf("expected %q to contain %q", got, want)
	}
}
