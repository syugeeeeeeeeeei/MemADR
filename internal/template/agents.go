package template

import "strings"

func RenderAgentsSnippet() string {
	lines := []string{
		"## MemADR Policy",
		"",
		"Use `memadr` for development knowledge records.",
		"",
		"- Use Japanese for memory records by default.",
		"- Keep `memory/` records short and outcome-focused.",
		"- Do not write work logs under `memory/`.",
		"- Before creating a new BUG, PROB, ADR, CHG, REV, SOL, or SUP record, search existing records with `memadr search` and `memadr related`.",
		"- When starting a repository, run `memadr init` and read `MEMADR_WORKFLOW.md`.",
		"- Create or update memory records when a bug, structural problem, architectural decision, change, rollback, reusable solution, or supersession is found.",
		"- Use `memadr check` before finishing work when records were added or updated.",
		"- Use `memadr index` to regenerate aggregate files instead of editing `memory/generated/` manually.",
		"- Treat each record file as the source of truth for its own status.",
	}

	return strings.Join(lines, "\n")
}
