package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"memadr/internal/release"
)

func main() {
	if len(os.Args) < 2 {
		fail("usage: buildctl <build|release>")
	}

	switch os.Args[1] {
	case "build":
		runBuild(os.Args[2:])
	case "dist":
		runDist(os.Args[2:])
	case "release":
		runRelease(os.Args[2:])
	default:
		fail("unknown subcommand: " + os.Args[1])
	}
}

func runBuild(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	out := fs.String("o", binaryName(), "output binary")
	target := fs.String("target", "./cmd/memadr", "build target")
	goBin := fs.String("go", "go", "go binary")
	fs.Parse(args)

	meta, err := release.DevMeta()
	if err != nil {
		fail(err.Error())
	}

	root, err := repoRoot()
	if err != nil {
		fail(err.Error())
	}

	outPath := absPath(root, *out)

	if err := goBuild(root, *goBin, outPath, *target, meta); err != nil {
		fail(err.Error())
	}
}

func runRelease(args []string) {
	fs := flag.NewFlagSet("release", flag.ExitOnError)
	version := fs.String("version", "", "release version")
	out := fs.String("o", binaryName(), "output binary")
	target := fs.String("target", "./cmd/memadr", "build target")
	goBin := fs.String("go", "go", "go binary")
	fs.Parse(args)

	if *version == "" {
		fail("release version is required")
	}
	if err := release.ValidateVersion(*version); err != nil {
		fail(err.Error())
	}
	if err := release.EnsureClean(); err != nil {
		fail(err.Error())
	}
	if err := release.EnsureTagMissing(*version); err != nil {
		fail(err.Error())
	}

	root, err := repoRoot()
	if err != nil {
		fail(err.Error())
	}

	if err := runCmd(root, *goBin, "test", "./..."); err != nil {
		fail(err.Error())
	}

	meta, err := release.ReleaseMeta(*version)
	if err != nil {
		fail(err.Error())
	}
	outPath := absPath(root, *out)
	if err := goBuild(root, *goBin, outPath, *target, meta); err != nil {
		fail(err.Error())
	}
	if err := release.CreateTag(*version); err != nil {
		fail(err.Error())
	}
	if err := release.PushRelease(*version); err != nil {
		fail(err.Error())
	}
}

func runDist(args []string) {
	fs := flag.NewFlagSet("dist", flag.ExitOnError)
	version := fs.String("version", "", "release version")
	outDir := fs.String("out-dir", "dist", "output directory")
	target := fs.String("target", "./cmd/memadr", "build target")
	goBin := fs.String("go", "go", "go binary")
	fs.Parse(args)

	if *version == "" {
		fail("dist version is required")
	}

	root, err := repoRoot()
	if err != nil {
		fail(err.Error())
	}

	meta, err := release.ReleaseMeta(*version)
	if err != nil {
		fail(err.Error())
	}
	distDir := absPath(root, *outDir)
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		fail(err.Error())
	}

	for _, item := range release.Targets() {
		out := filepath.Join(distDir, release.BinaryName(item))
		if err := goBuildFor(root, *goBin, out, *target, meta, item); err != nil {
			fail(err.Error())
		}
	}
}

func goBuild(root string, goBin string, out string, target string, meta release.Meta) error {
	return runCmd(root, goBin, "build", "-ldflags", release.Ldflags(meta), "-o", out, target)
}

func goBuildFor(root string, goBin string, out string, target string, meta release.Meta, targetInfo release.Target) error {
	cmd := exec.Command(goBin, "build", "-ldflags", release.Ldflags(meta), "-o", out, target)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"CGO_ENABLED=0",
		"GOOS="+targetInfo.GOOS,
		"GOARCH="+targetInfo.GOARCH,
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s build %s/%s failed: %w", goBin, targetInfo.GOOS, targetInfo.GOARCH, err)
	}
	return nil
}

func runCmd(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %v failed: %w", name, args, err)
	}
	return nil
}

func binaryName() string {
	if os.PathSeparator == '\\' {
		return "memadr.exe"
	}
	return "memadr"
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func repoRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git rev-parse --show-toplevel failed: %s", string(out))
	}
	return filepath.Clean(stringTrim(string(out))), nil
}

func absPath(root string, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

func stringTrim(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}
