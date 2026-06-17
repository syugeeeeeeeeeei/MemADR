package release

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestInstallScriptContainsReleaseURLs(t *testing.T) {
	t.Parallel()

	got := InstallScript("syugeeeeeeeeeei/MemADR")

	want := []string{
		`repo="syugeeeeeeeeeei/MemADR"`,
		`asset="memadr_linux_amd64"`,
		`url="$base/latest/download/$asset"`,
		`url="$base/download/$version/$asset"`,
		`-b | --bin-dir`,
		`-v | --version`,
		`-r | --repo`,
	}

	for _, item := range want {
		if !strings.Contains(got, item) {
			t.Fatalf("InstallScript() missing %q", item)
		}
	}
}

func TestWriteInstallScriptMakesExecutable(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path, err := WriteInstallScript(dir, "syugeeeeeeeeeei/MemADR")
	if err != nil {
		t.Fatalf("WriteInstallScript() error: %v", err)
	}
	if filepath.Base(path) != InstallName() {
		t.Fatalf("WriteInstallScript() path = %q", path)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat(%q) error: %v", path, err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0o755 {
		t.Fatalf("mode = %o", info.Mode().Perm())
	}
}

func TestInstallScriptMatchesRepoScript(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "scripts", InstallName())
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error: %v", path, err)
	}
	if string(got) != InstallScript(DefaultRepo) {
		t.Fatalf("%s is out of sync with InstallScript", path)
	}
}
