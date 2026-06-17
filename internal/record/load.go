package record

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"memadr/internal/mem"
)

var (
	headRe  = regexp.MustCompile(`^# ([A-Z]+-\d{3}): (.+)$`)
	fieldRe = regexp.MustCompile(`^([A-Za-z][A-Za-z0-9]*):\s*(.*)$`)
	idRe    = regexp.MustCompile(`\b[A-Z]+-\d{3}\b`)
)

func LoadAll(root string) ([]Rec, error) {
	var out []Rec
	for _, kind := range mem.Kinds() {
		dir := filepath.Join(root, "memory", kind.Dir)
		items, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			if item.IsDir() || filepath.Ext(item.Name()) != ".md" {
				continue
			}
			path := filepath.Join(dir, item.Name())
			rec, err := LoadFile(path, kind)
			if err != nil {
				return nil, err
			}
			out = append(out, rec)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out, nil
}

func LoadByID(root string, id string) (Rec, error) {
	prefix := idPrefix(id)
	kind, ok := mem.KindByPrefix(prefix)
	if !ok {
		return Rec{}, fmt.Errorf("unknown record id: %s", id)
	}

	path := filepath.Join(root, "memory", kind.Dir, id+".md")
	return LoadFile(path, kind)
}

func LoadFile(path string, kind mem.Kind) (Rec, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return Rec{}, err
	}

	rec := Rec{
		Path:   path,
		Kind:   kind,
		Fields: map[string]string{},
		Raw:    string(body),
	}

	sc := bufio.NewScanner(strings.NewReader(rec.Raw))
	lineNo := 0
	for sc.Scan() {
		line := sc.Text()
		if lineNo == 0 {
			m := headRe.FindStringSubmatch(line)
			if len(m) == 3 {
				rec.ID = m[1]
				rec.Title = strings.TrimSpace(m[2])
			}
			lineNo++
			continue
		}

		m := fieldRe.FindStringSubmatch(line)
		if len(m) == 3 {
			rec.Fields[m[1]] = strings.TrimSpace(m[2])
		}
		lineNo++
	}
	if err := sc.Err(); err != nil {
		return Rec{}, err
	}

	return rec, nil
}

func FindRefs(text string) []string {
	return idRe.FindAllString(text, -1)
}

func idPrefix(id string) string {
	idx := strings.IndexByte(id, '-')
	if idx < 0 {
		return id
	}
	return id[:idx]
}
