package record

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"memadr/internal/mem"
)

var generatedMark = "<!-- This file is generated. Do not edit manually. -->"

func Validate(root string, recs []Rec) []Issue {
	byID := make(map[string]Rec, len(recs))
	var issues []Issue

	for _, rec := range recs {
		if _, ok := byID[rec.ID]; ok {
			issues = append(issues, Issue{Path: rec.Path, Msg: "Duplicate ID: " + rec.ID})
			continue
		}
		byID[rec.ID] = rec
	}

	for _, rec := range recs {
		issues = append(issues, validateRec(rec, byID)...)
	}

	issues = append(issues, validateGenerated(root)...)
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Path == issues[j].Path {
			return issues[i].Msg < issues[j].Msg
		}
		return issues[i].Path < issues[j].Path
	})
	return issues
}

func validateRec(rec Rec, byID map[string]Rec) []Issue {
	var issues []Issue

	wantName := rec.ID + ".md"
	if rec.ID == "" {
		issues = append(issues, Issue{Path: rec.Path, Msg: "Invalid heading: expected `# ID: title`"})
	} else if rec.FileName() != wantName {
		issues = append(issues, Issue{Path: rec.Path, Msg: "File name does not match ID: " + rec.ID})
	}

	for _, name := range requiredFields(rec.Kind.Name) {
		if strings.TrimSpace(rec.Fields[name]) == "" {
			issues = append(issues, Issue{Path: rec.Path, Msg: "Missing field: " + name})
		}
	}

	status := rec.Status()
	if status != "" && !isAllowed(status, mem.Statuses()) {
		issues = append(issues, Issue{Path: rec.Path, Msg: "Invalid Status: " + status})
	}

	future := rec.Future()
	if future != "" && !isAllowed(future, mem.FutureValues()) {
		issues = append(issues, Issue{Path: rec.Path, Msg: "Invalid FutureRelevance: " + future})
	}

	for _, ref := range refsInRec(rec) {
		if ref == rec.ID {
			continue
		}
		if _, ok := byID[ref]; !ok {
			issues = append(issues, Issue{Path: rec.Path, Msg: "Missing reference: " + ref})
		}
	}

	return issues
}

func validateGenerated(root string) []Issue {
	dir := filepath.Join(root, "memory", "generated")
	items, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return []Issue{{Path: dir, Msg: err.Error()}}
	}

	var issues []Issue
	for _, item := range items {
		if item.IsDir() || filepath.Ext(item.Name()) != ".md" {
			continue
		}
		path := filepath.Join(dir, item.Name())
		body, err := os.ReadFile(path)
		if err != nil {
			issues = append(issues, Issue{Path: path, Msg: err.Error()})
			continue
		}
		if !strings.HasPrefix(string(body), generatedMark) {
			issues = append(issues, Issue{
				Path: path,
				Msg:  fmt.Sprintf("Generated file is missing comment: %s", generatedMark),
			})
		}
	}
	return issues
}

func requiredFields(kind string) []string {
	switch kind {
	case "bug":
		return []string{"Status", "Area", "Symptom"}
	case "adr":
		return []string{"Status", "Context", "Decision", "Rejected", "Consequence"}
	case "chg":
		return []string{"Status", "Reason", "Change", "FilesChanged", "Verification"}
	case "rev":
		return []string{"Status", "Reverted", "Reason", "Result", "FollowUp"}
	case "sol":
		return []string{"Status", "ProblemPattern", "Solution", "AppliesTo"}
	case "sup":
		return []string{"Status", "Supersedes", "Reason", "NewBaseline"}
	case "prob":
		return []string{"Status", "Area", "Impact", "Finding", "CurrentDirection"}
	default:
		return []string{"Status"}
	}
}

func refsInRec(rec Rec) []string {
	set := map[string]struct{}{}
	for _, val := range rec.Fields {
		for _, ref := range FindRefs(val) {
			set[ref] = struct{}{}
		}
	}

	out := make([]string, 0, len(set))
	for ref := range set {
		out = append(out, ref)
	}
	sort.Strings(out)
	return out
}

func isAllowed(val string, allowed []string) bool {
	for _, item := range allowed {
		if item == val {
			return true
		}
	}
	return false
}
