package release

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var verRe = regexp.MustCompile(`^v\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)

type Meta struct {
	Version string
	Commit  string
	Date    string
	Dirty   bool
}

type Target struct {
	GOOS   string
	GOARCH string
	Ext    string
}

func ValidateVersion(version string) error {
	if !verRe.MatchString(version) {
		return fmt.Errorf("version must look like v0.1.0: %s", version)
	}
	return nil
}

func DevMeta() (Meta, error) {
	commit, err := git("rev-parse", "HEAD")
	if err != nil {
		return Meta{}, err
	}

	dirty, err := isDirty()
	if err != nil {
		return Meta{}, err
	}

	base, err := latestTag()
	if err != nil {
		return Meta{}, err
	}

	exact, err := exactTag()
	if err != nil {
		return Meta{}, err
	}

	version := DevVersion(base, commit)
	if exact != "" && !dirty {
		version = exact
	}

	return Meta{
		Version: version,
		Commit:  commit,
		Date:    time.Now().UTC().Format(time.RFC3339),
		Dirty:   dirty,
	}, nil
}

func ReleaseMeta(version string) (Meta, error) {
	if err := ValidateVersion(version); err != nil {
		return Meta{}, err
	}

	commit, err := git("rev-parse", "HEAD")
	if err != nil {
		return Meta{}, err
	}

	dirty, err := isDirty()
	if err != nil {
		return Meta{}, err
	}

	return Meta{
		Version: version,
		Commit:  commit,
		Date:    time.Now().UTC().Format(time.RFC3339),
		Dirty:   dirty,
	}, nil
}

func EnsureClean() error {
	dirty, err := isDirty()
	if err != nil {
		return err
	}
	if dirty {
		return fmt.Errorf("working tree must be clean")
	}
	return nil
}

func EnsureTagMissing(version string) error {
	err := exec.Command("git", "rev-parse", "--verify", version).Run()
	if err == nil {
		return fmt.Errorf("tag already exists: %s", version)
	}
	return nil
}

func CreateTag(version string) error {
	cmd := exec.Command("git", "tag", "-a", version, "-m", "release "+version)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git tag failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func PushRelease(version string) error {
	if _, err := git("push", "origin", "HEAD"); err != nil {
		return err
	}
	if _, err := git("push", "origin", version); err != nil {
		return err
	}
	return nil
}

func Ldflags(meta Meta) string {
	dirty := "false"
	if meta.Dirty {
		dirty = "true"
	}

	items := []string{
		"-X", "memadr/internal/buildinfo.version=" + meta.Version,
		"-X", "memadr/internal/buildinfo.commit=" + meta.Commit,
		"-X", "memadr/internal/buildinfo.date=" + meta.Date,
		"-X", "memadr/internal/buildinfo.dirty=" + dirty,
	}
	return strings.Join(items, " ")
}

func DevVersion(base string, commit string) string {
	if base == "" {
		base = "v0.0.0"
	}

	short := commit
	if len(short) > 7 {
		short = short[:7]
	}
	if short == "" {
		return base + "-dev"
	}
	return base + "-dev+" + short
}

func Targets() []Target {
	return []Target{
		{GOOS: "windows", GOARCH: "amd64", Ext: ".exe"},
		{GOOS: "darwin", GOARCH: "amd64"},
		{GOOS: "darwin", GOARCH: "arm64"},
		{GOOS: "linux", GOARCH: "amd64"},
	}
}

func BinaryName(target Target) string {
	return fmt.Sprintf("memadr_%s_%s%s", target.GOOS, target.GOARCH, target.Ext)
}

func isDirty() (bool, error) {
	out, err := exec.Command("git", "status", "--porcelain").CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("git status failed: %s", strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)) != "", nil
}

func git(args ...string) (string, error) {
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s failed: %s", strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func latestTag() (string, error) {
	out, err := exec.Command("git", "tag", "--list", "v*", "--sort=-version:refname").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git latest tag failed: %s", strings.TrimSpace(string(out)))
	}
	lines := strings.Fields(strings.TrimSpace(string(out)))
	if len(lines) == 0 {
		return "", nil
	}
	return lines[0], nil
}

func exactTag() (string, error) {
	out, err := exec.Command("git", "tag", "--points-at", "HEAD", "--list", "v*", "--sort=-version:refname").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git exact tag failed: %s", strings.TrimSpace(string(out)))
	}
	lines := strings.Fields(strings.TrimSpace(string(out)))
	if len(lines) == 0 {
		return "", nil
	}
	return lines[0], nil
}
