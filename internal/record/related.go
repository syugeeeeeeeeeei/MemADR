package record

import (
	"regexp"
	"sort"
	"strings"
)

var wordRe = regexp.MustCompile(`[\pL\pN]+`)

type Candidate struct {
	Rec   Rec
	Score int
}

func Related(recs []Rec, targetID string) ([]Candidate, error) {
	var target Rec
	found := false
	for _, rec := range recs {
		if rec.ID == targetID {
			target = rec
			found = true
			break
		}
	}
	if !found {
		return nil, errNotFound(targetID)
	}

	targetRefs := refSet(target)
	targetWords := wordSet(target)
	var out []Candidate

	for _, rec := range recs {
		if rec.ID == target.ID {
			continue
		}

		score := 0
		if target.Area() != "" && strings.EqualFold(target.Area(), rec.Area()) {
			score += 3
		}
		if targetRefs[rec.ID] || refSet(rec)[target.ID] {
			score += 5
		}
		for word := range wordSet(rec) {
			if targetWords[word] {
				score++
			}
		}
		if score == 0 {
			continue
		}

		out = append(out, Candidate{Rec: rec, Score: score})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].Rec.ID < out[j].Rec.ID
		}
		return out[i].Score > out[j].Score
	})
	return out, nil
}

func wordSet(rec Rec) map[string]bool {
	set := map[string]bool{}
	text := rec.Title + "\n"
	for _, val := range rec.Fields {
		text += val + "\n"
	}
	for _, word := range wordRe.FindAllString(strings.ToLower(text), -1) {
		if len([]rune(word)) < 2 {
			continue
		}
		set[word] = true
	}
	return set
}

func refSet(rec Rec) map[string]bool {
	set := map[string]bool{}
	for _, ref := range refsInRec(rec) {
		set[ref] = true
	}
	return set
}

type errMissing struct {
	id string
}

func (e errMissing) Error() string {
	return "record not found: " + e.id
}

func errNotFound(id string) error {
	return errMissing{id: id}
}
