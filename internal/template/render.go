package template

import (
	"strings"

	"memadr/internal/mem"
)

func Render(kind mem.Kind, id string, title string) string {
	lines := []string{
		"# " + id + ": " + title,
		"",
	}

	lines = append(lines, fields(kind)...)
	lines = append(lines, "")
	return strings.Join(lines, "\n")
}

func fields(kind mem.Kind) []string {
	switch kind.Name {
	case "bug":
		return []string{
			"Status: OPEN",
			"Area: ",
			"Symptom: ",
			"Cause: ",
			"Fix: ",
			"Verification: ",
			"FutureRelevance: watch",
			"Decision: ",
			"ResolvedBy: ",
		}
	case "prob":
		return []string{
			"Status: UNRESOLVED",
			"Area: ",
			"Impact: ",
			"Finding: ",
			"CurrentDirection: ",
			"Related: ",
		}
	case "adr":
		return []string{
			"Status: PROPOSED",
			"Context: ",
			"Decision: ",
			"Rejected: ",
			"Consequence: ",
			"Related: ",
		}
	case "chg":
		return []string{
			"Status: SHIPPED",
			"Reason: ",
			"Change: ",
			"FilesChanged: ",
			"Verification: ",
		}
	case "rev":
		return []string{
			"Status: CLOSED",
			"Reverted: ",
			"Reason: ",
			"Result: ",
			"FollowUp: ",
		}
	case "sol":
		return []string{
			"Status: ACTIVE",
			"ProblemPattern: ",
			"Solution: ",
			"AppliesTo: ",
			"FutureRelevance: reusable",
			"Related: ",
		}
	case "sup":
		return []string{
			"Status: ACTIVE",
			"Supersedes: ",
			"Reason: ",
			"NewBaseline: ",
		}
	default:
		return []string{"Status: OPEN"}
	}
}
