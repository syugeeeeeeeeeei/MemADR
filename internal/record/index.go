package record

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func WriteIndexes(root string, recs []Rec) error {
	dir := filepath.Join(root, "memory", "generated")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	files := map[string]string{
		"active.md":    renderListDoc("Active", Filtered(recs, Filter{}), isActive),
		"by-status.md": renderGroupedDoc("By Status", recs, func(rec Rec) string { return rec.Status() }),
		"by-area.md":   renderGroupedDoc("By Area", recs, func(rec Rec) string { return rec.Area() }),
		"by-type.md":   renderGroupedDoc("By Type", recs, func(rec Rec) string { return strings.ToUpper(rec.Type()) }),
		"open.md":      renderListDoc("Open", Filtered(recs, Filter{}), isOpen),
	}

	for name, body := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func renderListDoc(title string, recs []Rec, keep func(Rec) bool) string {
	var b strings.Builder
	b.WriteString(generatedMark)
	b.WriteString("\n\n# ")
	b.WriteString(title)
	b.WriteString("\n")

	for _, rec := range recs {
		if keep != nil && !keep(rec) {
			continue
		}
		fmt.Fprintf(&b, "- %s | %s | %s | %s\n", rec.ID, rec.Status(), blankAsDash(rec.Area()), rec.Title)
	}
	return b.String()
}

func renderGroupedDoc(title string, recs []Rec, group func(Rec) string) string {
	buckets := map[string][]Rec{}
	var keys []string

	for _, rec := range recs {
		key := blankAsDash(group(rec))
		if _, ok := buckets[key]; !ok {
			keys = append(keys, key)
		}
		buckets[key] = append(buckets[key], rec)
	}

	sort.Strings(keys)

	var b strings.Builder
	b.WriteString(generatedMark)
	b.WriteString("\n\n# ")
	b.WriteString(title)
	b.WriteString("\n")

	for _, key := range keys {
		b.WriteString("\n## ")
		b.WriteString(key)
		b.WriteString("\n")
		for _, rec := range buckets[key] {
			fmt.Fprintf(&b, "- %s | %s | %s | %s\n", rec.ID, rec.Status(), blankAsDash(rec.Area()), rec.Title)
		}
	}
	return b.String()
}

func isActive(rec Rec) bool {
	switch rec.Status() {
	case "OPEN", "INVESTIGATING", "UNRESOLVED", "PROPOSED", "ACCEPTED", "FIXED", "VERIFIED", "ACTIVE":
		return true
	default:
		return false
	}
}

func isOpen(rec Rec) bool {
	switch rec.Status() {
	case "OPEN", "INVESTIGATING", "UNRESOLVED", "PROPOSED", "FIXED":
		return true
	default:
		return false
	}
}

func blankAsDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}
