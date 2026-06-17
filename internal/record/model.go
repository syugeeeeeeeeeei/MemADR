package record

import (
	"path/filepath"
	"strings"

	"memadr/internal/mem"
)

type Rec struct {
	ID     string
	Title  string
	Path   string
	Kind   mem.Kind
	Fields map[string]string
	Raw    string
}

func (r Rec) Status() string {
	return r.Fields["Status"]
}

func (r Rec) Area() string {
	return r.Fields["Area"]
}

func (r Rec) Future() string {
	return r.Fields["FutureRelevance"]
}

func (r Rec) Type() string {
	return r.Kind.Name
}

func (r Rec) FileName() string {
	return filepath.Base(r.Path)
}

func (r Rec) SearchText() string {
	var b strings.Builder
	b.WriteString(r.ID)
	b.WriteString("\n")
	b.WriteString(r.Title)
	b.WriteString("\n")
	for key, val := range r.Fields {
		b.WriteString(key)
		b.WriteString(": ")
		b.WriteString(val)
		b.WriteString("\n")
	}
	b.WriteString(r.Raw)
	return strings.ToLower(b.String())
}

type Filter struct {
	Type   string
	Status string
	Area   string
	Future string
}

type Issue struct {
	Path string
	Msg  string
}
