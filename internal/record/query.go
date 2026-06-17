package record

import "strings"

func Filtered(recs []Rec, f Filter) []Rec {
	var out []Rec
	for _, rec := range recs {
		if f.Type != "" && rec.Type() != strings.ToLower(f.Type) {
			continue
		}
		if f.Status != "" && !strings.EqualFold(rec.Status(), f.Status) {
			continue
		}
		if f.Area != "" && !strings.EqualFold(rec.Area(), f.Area) {
			continue
		}
		if f.Future != "" && !strings.EqualFold(rec.Future(), f.Future) {
			continue
		}
		out = append(out, rec)
	}
	return out
}

func Search(recs []Rec, q string) []Rec {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return nil
	}

	var out []Rec
	for _, rec := range recs {
		if strings.Contains(rec.SearchText(), q) {
			out = append(out, rec)
		}
	}
	return out
}
