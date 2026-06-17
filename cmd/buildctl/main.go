package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

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

	if err := goBuild(*goBin, *out, *target, meta); err != nil {
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
	if err := runCmd(*goBin, "test", "./..."); err != nil {
		fail(err.Error())
	}

	meta, err := release.ReleaseMeta(*version)
	if err != nil {
		fail(err.Error())
	}
	if err := goBuild(*goBin, *out, *target, meta); err != nil {
		fail(err.Error())
	}
	if err := release.CreateTag(*version); err != nil {
		fail(err.Error())
	}
	if err := release.PushRelease(*version); err != nil {
		fail(err.Error())
	}
}

func goBuild(goBin string, out string, target string, meta release.Meta) error {
	return runCmd(goBin, "build", "-ldflags", release.Ldflags(meta), "-o", out, target)
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

	meta, err := release.ReleaseMeta(*version)
	if err != nil {
		fail(err.Error())
	}
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fail(err.Error())
	}

	for _, item := range release.Targets() {
		out := *outDir + string(os.PathSeparator) + release.BinaryName(item)
		if err := goBuildFor(*goBin, out, *target, meta, item); err != nil {
			fail(err.Error())
		}
	}
}

func goBuildFor(goBin string, out string, target string, meta release.Meta, targetInfo release.Target) error {
	cmd := exec.Command(goBin, "build", "-ldflags", release.Ldflags(meta), "-o", out, target)
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

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
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
