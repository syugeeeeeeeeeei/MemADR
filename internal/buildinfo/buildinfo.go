package buildinfo

import (
	"runtime/debug"
	"strconv"
	"strings"
)

var (
	version  string
	commit   string
	date     string
	dirty    string
	readInfo = debug.ReadBuildInfo
)

type Info struct {
	Version string
	Commit  string
	Date    string
	Dirty   bool
}

func Current() Info {
	info := Info{
		Version: strings.TrimSpace(version),
		Commit:  strings.TrimSpace(commit),
		Date:    strings.TrimSpace(date),
		Dirty:   parseBool(dirty),
	}

	if bi, ok := readInfo(); ok {
		applyBuildInfo(&info, bi)
	}

	if info.Version == "" {
		info.Version = fallbackDevVersion(info.Commit)
	}

	return info
}

func Render(info Info) string {
	return info.Version + "\n"
}

func RenderVerbose(info Info) string {
	var lines []string
	lines = append(lines, info.Version)
	if info.Commit != "" {
		lines = append(lines, "commit: "+info.Commit)
	}
	if info.Date != "" {
		lines = append(lines, "date: "+info.Date)
	}
	lines = append(lines, "state: "+stateLabel(info.Dirty))

	return strings.Join(lines, "\n") + "\n"
}

func applyBuildInfo(info *Info, bi *debug.BuildInfo) {
	if info.Version == "" && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		info.Version = bi.Main.Version
	}

	for _, item := range bi.Settings {
		switch item.Key {
		case "vcs.revision":
			if info.Commit == "" {
				info.Commit = item.Value
			}
		case "vcs.time":
			if info.Date == "" {
				info.Date = item.Value
			}
		case "vcs.modified":
			if dirty == "" {
				info.Dirty = parseBool(item.Value)
			}
		}
	}
}

func fallbackDevVersion(commit string) string {
	short := shortCommit(commit)
	if short == "" {
		return "v0.0.0-dev"
	}
	return "v0.0.0-dev+" + short
}

func shortCommit(commit string) string {
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}

func parseBool(raw string) bool {
	ok, err := strconv.ParseBool(strings.TrimSpace(raw))
	return err == nil && ok
}

func stateLabel(dirty bool) string {
	if dirty {
		return "modified"
	}
	return "clean"
}
