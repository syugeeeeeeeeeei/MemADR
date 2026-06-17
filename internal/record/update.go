package record

import (
	"os"
	"sort"
	"strings"
)

func UpdateFields(path string, vals map[string]string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(body), "\n")
	seen := map[string]bool{}
	keys := make([]string, 0, len(vals))
	for key := range vals {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for i, line := range lines {
		for _, key := range keys {
			val := vals[key]
			prefix := key + ":"
			if strings.HasPrefix(line, prefix) {
				lines[i] = prefix + " " + val
				seen[key] = true
			}
		}
	}

	insertAt := len(lines)
	for insertAt > 0 && lines[insertAt-1] == "" {
		insertAt--
	}
	for _, key := range keys {
		if seen[key] {
			continue
		}
		val := vals[key]
		lines = append(lines[:insertAt], append([]string{key + ": " + val}, lines[insertAt:]...)...)
		insertAt++
	}

	out := strings.Join(lines, "\n")
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	return os.WriteFile(path, []byte(out), 0o644)
}
