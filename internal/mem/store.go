package mem

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var idRe = regexp.MustCompile(`^([A-Z]+)-(\d{3})\.md$`)

func Init(root string) error {
	for _, rel := range Dirs() {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(path, 0o755); err != nil {
			return err
		}
	}
	return nil
}

func NextPath(root string, kind Kind) (string, string, error) {
	dir := filepath.Join(root, "memory", kind.Dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", "", err
	}

	next, err := nextNum(dir, kind.Prefix)
	if err != nil {
		return "", "", err
	}

	id := fmt.Sprintf("%s-%03d", kind.Prefix, next)
	path := filepath.Join(dir, id+".md")
	return path, id, nil
}

func nextNum(dir string, prefix string) (int, error) {
	items, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	max := 0
	for _, item := range items {
		if item.IsDir() {
			continue
		}

		m := idRe.FindStringSubmatch(item.Name())
		if len(m) != 3 || m[1] != prefix {
			continue
		}

		n, err := strconv.Atoi(m[2])
		if err != nil {
			return 0, err
		}
		if n > max {
			max = n
		}
	}

	return max + 1, nil
}

func WriteFile(path string, body string) error {
	return os.WriteFile(path, []byte(body), 0o644)
}
