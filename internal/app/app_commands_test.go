package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCheckRejectsInvalidStatus(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	mustInit(t, dir)
	mustWrite(t, dir, "memory/bugs/BUG-001.md", strings.Join([]string{
		"# BUG-001: 認証状態が壊れる",
		"",
		"Status: DONE",
		"Area: auth",
		"Symptom: 壊れる",
		"",
	}, "\n"))

	var out bytes.Buffer
	err := Run([]string{"check"}, dir, &out, &out)
	if err == nil {
		t.Fatal("expected check to fail")
	}

	mustContain(t, out.String(), "Invalid Status: DONE")
	mustContain(t, out.String(), "BUG-001.md")
}

func TestRunCheckRejectsMissingReference(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	mustInit(t, dir)
	mustWrite(t, dir, "memory/bugs/BUG-001.md", strings.Join([]string{
		"# BUG-001: 認証状態が壊れる",
		"",
		"Status: VERIFIED",
		"Area: auth",
		"Symptom: 壊れる",
		"ResolvedBy: CHG-999",
		"",
	}, "\n"))

	var out bytes.Buffer
	err := Run([]string{"check"}, dir, &out, &out)
	if err == nil {
		t.Fatal("expected check to fail")
	}

	mustContain(t, out.String(), "Missing reference: CHG-999")
}

func TestRunListFiltersRecords(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	mustInit(t, dir)
	writeFixtureSet(t, dir)

	var out bytes.Buffer
	if err := Run([]string{"list", "--status", "VERIFIED"}, dir, &out, &out); err != nil {
		t.Fatalf("list error: %v", err)
	}

	text := out.String()
	mustContain(t, text, "BUG-001  VERIFIED  auth")
	mustNotContain(t, text, "ADR-001")

	out.Reset()
	if err := Run([]string{"list", "--area", "auth"}, dir, &out, &out); err != nil {
		t.Fatalf("list error: %v", err)
	}

	text = out.String()
	mustContain(t, text, "BUG-001")
	mustContain(t, text, "BUG-002")
	mustNotContain(t, text, "SOL-001")

	out.Reset()
	if err := Run([]string{"list", "--type", "adr"}, dir, &out, &out); err != nil {
		t.Fatalf("list error: %v", err)
	}

	text = out.String()
	mustContain(t, text, "ADR-001")
	mustNotContain(t, text, "BUG-001")
}

func TestRunSearchFindsRecordsByText(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	mustInit(t, dir)
	writeFixtureSet(t, dir)

	var out bytes.Buffer
	if err := Run([]string{"search", "server session"}, dir, &out, &out); err != nil {
		t.Fatalf("search error: %v", err)
	}

	text := out.String()
	mustContain(t, text, "BUG-001")
	mustContain(t, text, "ADR-001")
	mustNotContain(t, text, "SOL-001")
}

func TestRunIndexGeneratesViews(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	mustInit(t, dir)
	writeFixtureSet(t, dir)

	var out bytes.Buffer
	if err := Run([]string{"index"}, dir, &out, &out); err != nil {
		t.Fatalf("index error: %v", err)
	}

	active := mustRead(t, dir, "memory/generated/active.md")
	mustContain(t, active, "<!-- This file is generated. Do not edit manually. -->")
	mustContain(t, active, "ADR-001")
	mustContain(t, active, "SOL-001")
	mustNotContain(t, active, "CHG-001")

	byArea := mustRead(t, dir, "memory/generated/by-area.md")
	mustContain(t, byArea, "## auth")
	mustContain(t, byArea, "BUG-001")

	open := mustRead(t, dir, "memory/generated/open.md")
	mustContain(t, open, "BUG-002")
	mustNotContain(t, open, "BUG-001")
}

func TestRunRelatedShowsCandidates(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	mustInit(t, dir)
	writeFixtureSet(t, dir)

	var out bytes.Buffer
	if err := Run([]string{"related", "BUG-001"}, dir, &out, &out); err != nil {
		t.Fatalf("related error: %v", err)
	}

	text := out.String()
	mustContain(t, text, "Related candidates for BUG-001:")
	mustContain(t, text, "ADR-001")
	mustContain(t, text, "CHG-001")
}

func TestRunCloseUpdatesStatusAndResolvedBy(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	mustInit(t, dir)
	mustWrite(t, dir, "memory/bugs/BUG-001.md", strings.Join([]string{
		"# BUG-001: 認証状態が壊れる",
		"",
		"Status: FIXED",
		"Area: auth",
		"Symptom: 壊れる",
		"",
	}, "\n"))
	mustWrite(t, dir, "memory/changes/CHG-001.md", strings.Join([]string{
		"# CHG-001: 修正",
		"",
		"Status: SHIPPED",
		"Reason: ADR-001",
		"Change: 修正した",
		"FilesChanged: internal/auth/**",
		"Verification: 実施済み",
		"",
	}, "\n"))
	mustWrite(t, dir, "memory/decisions/ADR-001.md", strings.Join([]string{
		"# ADR-001: 判断",
		"",
		"Status: ACCEPTED",
		"Context: 既存",
		"Decision: 維持",
		"Rejected: なし",
		"Consequence: 影響なし",
		"",
	}, "\n"))

	var out bytes.Buffer
	if err := Run([]string{"close", "BUG-001", "--resolved-by", "CHG-001", "--verified"}, dir, &out, &out); err != nil {
		t.Fatalf("close error: %v", err)
	}

	body := mustRead(t, dir, "memory/bugs/BUG-001.md")
	mustContain(t, body, "Status: VERIFIED")
	mustContain(t, body, "ResolvedBy: CHG-001")
}

func mustInit(t *testing.T, dir string) {
	t.Helper()

	var out bytes.Buffer
	if err := Run([]string{"init"}, dir, &out, &out); err != nil {
		t.Fatalf("init error: %v", err)
	}
}

func writeFixtureSet(t *testing.T, dir string) {
	t.Helper()

	mustWrite(t, dir, "memory/bugs/BUG-001.md", strings.Join([]string{
		"# BUG-001: 認証状態が壊れる",
		"",
		"Status: VERIFIED",
		"Area: auth",
		"Symptom: 画面遷移後にログイン状態が壊れる",
		"Cause: server session と client state が分散していた",
		"Fix: server session 中心に整理した",
		"Verification: ログインと再ログインを確認した",
		"FutureRelevance: watch",
		"Decision: ADR-001",
		"ResolvedBy: CHG-001",
		"",
	}, "\n"))
	mustWrite(t, dir, "memory/bugs/BUG-002.md", strings.Join([]string{
		"# BUG-002: 認証期限切れ表示が不足する",
		"",
		"Status: OPEN",
		"Area: auth",
		"Symptom: 期限切れ時の表示が不足する",
		"FutureRelevance: watch",
		"Related: ADR-001",
		"",
	}, "\n"))
	mustWrite(t, dir, "memory/decisions/ADR-001.md", strings.Join([]string{
		"# ADR-001: 認証責務を server session 中心に整理する",
		"",
		"Status: ACCEPTED",
		"Context: BUG-001 の原因が責務分散だった",
		"Decision: server session を中心に責務を集約する",
		"Rejected: 局所パッチ継続",
		"Consequence: auth 依存を整理する",
		"Related: BUG-001",
		"",
	}, "\n"))
	mustWrite(t, dir, "memory/changes/CHG-001.md", strings.Join([]string{
		"# CHG-001: 認証関連ディレクトリを再編成",
		"",
		"Status: SHIPPED",
		"Reason: ADR-001",
		"Change: server session 中心に責務を集約した",
		"FilesChanged: internal/auth/**",
		"Verification: 主要導線を確認した",
		"",
	}, "\n"))
	mustWrite(t, dir, "memory/solutions/SOL-001.md", strings.Join([]string{
		"# SOL-001: UTCで期限判定する",
		"",
		"Status: ACTIVE",
		"ProblemPattern: JWT期限判定が環境差で壊れる",
		"Solution: UTCで統一する",
		"AppliesTo: auth, session",
		"FutureRelevance: reusable",
		"Related: BUG-002",
		"",
	}, "\n"))
}

func mustWrite(t *testing.T, dir string, rel string, body string) {
	t.Helper()

	path := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error: %v", path, err)
	}
}

func mustRead(t *testing.T, dir string, rel string) string {
	t.Helper()

	path := filepath.Join(dir, filepath.FromSlash(rel))
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error: %v", path, err)
	}
	return string(body)
}

func mustNotContain(t *testing.T, got string, want string) {
	t.Helper()

	if strings.Contains(got, want) {
		t.Fatalf("expected %q to not contain %q", got, want)
	}
}
