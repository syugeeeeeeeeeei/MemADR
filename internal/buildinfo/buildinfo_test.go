package buildinfo

import (
	"runtime/debug"
	"strings"
	"testing"
)

func TestCurrentFallsBackToDevVersionFromVCS(t *testing.T) {
	t.Parallel()

	restore := snapshot()
	defer restore()

	version = ""
	commit = ""
	date = ""
	dirty = ""
	readInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{Version: "(devel)"},
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "abcdef1234567890"},
				{Key: "vcs.time", Value: "2026-06-17T12:34:56Z"},
				{Key: "vcs.modified", Value: "true"},
			},
		}, true
	}

	got := Current()
	if got.Version != "v0.0.0-dev+abcdef1" {
		t.Fatalf("Version = %q", got.Version)
	}
	if got.Commit != "abcdef1234567890" {
		t.Fatalf("Commit = %q", got.Commit)
	}
	if got.Date != "2026-06-17T12:34:56Z" {
		t.Fatalf("Date = %q", got.Date)
	}
	if !got.Dirty {
		t.Fatal("Dirty = false")
	}
}

func TestCurrentPrefersInjectedVersion(t *testing.T) {
	t.Parallel()

	restore := snapshot()
	defer restore()

	version = "v0.1.0"
	commit = "abc1234"
	date = "2026-06-17T00:00:00Z"
	dirty = "false"
	readInfo = func() (*debug.BuildInfo, bool) {
		return nil, false
	}

	got := Current()
	if got.Version != "v0.1.0" {
		t.Fatalf("Version = %q", got.Version)
	}
	if got.Commit != "abc1234" {
		t.Fatalf("Commit = %q", got.Commit)
	}
	if got.Date != "2026-06-17T00:00:00Z" {
		t.Fatalf("Date = %q", got.Date)
	}
	if got.Dirty {
		t.Fatal("Dirty = true")
	}
}

func TestRenderIncludesDetails(t *testing.T) {
	t.Parallel()

	text := Render(Info{
		Version: "v0.1.0",
		Commit:  "abc1234",
		Date:    "2026-06-17T00:00:00Z",
		Dirty:   true,
	})

	if text != "v0.1.0\n" {
		t.Fatalf("Render = %q", text)
	}
}

func TestRenderVerboseIncludesState(t *testing.T) {
	t.Parallel()

	text := RenderVerbose(Info{
		Version: "v0.1.0-dev+abc1234",
		Commit:  "abc1234",
		Date:    "2026-06-17T00:00:00Z",
		Dirty:   true,
	})

	for _, want := range []string{
		"v0.1.0-dev+abc1234",
		"commit: abc1234",
		"date: 2026-06-17T00:00:00Z",
		"state: modified",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected %q to contain %q", text, want)
		}
	}
}

func snapshot() func() {
	oldVersion := version
	oldCommit := commit
	oldDate := date
	oldDirty := dirty
	oldReadInfo := readInfo

	return func() {
		version = oldVersion
		commit = oldCommit
		date = oldDate
		dirty = oldDirty
		readInfo = oldReadInfo
	}
}
