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

	guide := filepath.Join(dir, "MEMADR_WORKFLOW.md")
	body, err := os.ReadFile(guide)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error: %v", guide, err)
	}

	text := string(body)
	mustContain(t, text, "MemADR Workflow Guide")
	mustContain(t, text, "人間向け導線")
	mustContain(t, text, "LLM向け導線")
	mustContain(t, text, "memadr init")
	mustContain(t, text, "memadr new bug")

	outText := out.String()
	mustContain(t, outText, "initialized memory/")
	mustContain(t, outText, "LLMがMemADRを使用するために、以下をAGENTS.mdに貼り付けてください。")
	mustContain(t, outText, "## MemADR運用ポリシー")
	mustContain(t, outText, "開発知識レコードの管理には `memadr` を使用する。")
}

func TestRunInitDoesNotOverwriteExistingGuide(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	guide := filepath.Join(dir, "MEMADR_WORKFLOW.md")
	if err := os.WriteFile(guide, []byte("custom guide\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(%q) error: %v", guide, err)
	}

	var out bytes.Buffer
	if err := Run([]string{"init"}, dir, &out, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	body, err := os.ReadFile(guide)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error: %v", guide, err)
	}
	if string(body) != "custom guide\n" {
		t.Fatalf("guide was overwritten: %q", string(body))
	}
}

func TestRunWithoutArgsShowsHelp(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	var out bytes.Buffer

	if err := Run(nil, dir, &out, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	text := out.String()
	mustContain(t, text, "MemADR")
	mustContain(t, text, "Usage:")
	mustContain(t, text, "memadr init")
	mustContain(t, text, "memadr new <type> [title]")
	mustContain(t, text, "Record types:")
	mustContain(t, text, "BUG")
	mustContain(t, text, "Workflow guide:")
	mustContain(t, text, "MEMADR_WORKFLOW.md")
}

func TestRunHelpNewShowsDetailedHelp(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	var out bytes.Buffer

	if err := Run([]string{"help", "new"}, dir, &out, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	text := out.String()
	mustContain(t, text, "memadr new <type> [title]")
	mustContain(t, text, "Supported types:")
	mustContain(t, text, "bug")
	mustContain(t, text, "adr")
	mustContain(t, text, `memadr new bug "認証状態が壊れる"`)
	mustContain(t, text, "MEMADR_WORKFLOW.md")
}

func TestRunHelpListShowsOptionDetailsAndValues(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	var out bytes.Buffer

	if err := Run([]string{"help", "list"}, dir, &out, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	text := out.String()
	mustContain(t, text, "Options:")
	mustContain(t, text, "--type <TYPE>")
	mustContain(t, text, "値: bug, prob, adr, chg, rev, sol, sup")
	mustContain(t, text, "--status <STATUS>")
	mustContain(t, text, "値: OPEN, INVESTIGATING, UNRESOLVED")
	mustContain(t, text, "--future <FUTURE>")
	mustContain(t, text, "値: ignore, watch, reusable")
	mustContain(t, text, "--area <AREA>")
	mustContain(t, text, "任意のArea文字列")
}

func TestRunHelpCloseShowsOptionKinds(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	var out bytes.Buffer

	if err := Run([]string{"help", "close"}, dir, &out, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	text := out.String()
	mustContain(t, text, "Options:")
	mustContain(t, text, "--verified")
	mustContain(t, text, "値なし")
	mustContain(t, text, "--resolved-by <CHG-ID>")
	mustContain(t, text, "値: `CHG-001` のような変更ID")
}

func TestRunVersionShowsBuildVersion(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	var out bytes.Buffer

	if err := Run([]string{"version"}, dir, &out, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	text := out.String()
	mustContain(t, text, "v")
	mustNotContain(t, text, "dirty:")
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
