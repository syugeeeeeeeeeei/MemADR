package app

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"memadr/internal/buildinfo"
	"memadr/internal/mem"
	"memadr/internal/record"
	"memadr/internal/template"
)

func Run(args []string, wd string, out io.Writer, errOut io.Writer) error {
	if len(args) == 0 {
		return runHelp(nil, wd, out, errOut)
	}

	if args[0] == "-h" || args[0] == "--help" {
		return runHelp(nil, wd, out, errOut)
	}

	cmd, ok := findCommand(args[0])
	if !ok {
		return fmt.Errorf("unknown command: %s", args[0])
	}

	return cmd.Run(args[1:], wd, out, errOut)
}

func runInit(_ []string, wd string, out io.Writer, _ io.Writer) error {
	if err := mem.Init(wd); err != nil {
		return err
	}
	if err := mem.WriteFileIfAbsent(filepath.Join(wd, mem.WorkflowGuideFile), template.RenderWorkflowGuide()); err != nil {
		return err
	}

	_, err := fmt.Fprintln(out, "initialized memory/")
	return err
}

func runNew(args []string, wd string, out io.Writer, _ io.Writer) error {
	if len(args) == 0 {
		cmd, _ := findCommand("new")
		_, err := io.WriteString(out, renderCommandHelp(cmd))
		return err
	}

	kind, ok := mem.ParseKind(args[0])
	if !ok {
		return fmt.Errorf("unknown record type: %s", args[0])
	}

	title := "タイトル未設定"
	if len(args) > 1 && args[1] != "" {
		title = args[1]
	}

	path, id, err := mem.NextPath(wd, kind)
	if err != nil {
		return err
	}

	body := template.Render(kind, id, title)
	if err := mem.WriteFile(path, body); err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "created %s\n", path)
	return err
}

func runCheck(_ []string, wd string, out io.Writer, _ io.Writer) error {
	recs, err := record.LoadAll(wd)
	if err != nil {
		return err
	}

	issues := record.Validate(wd, recs)
	if len(issues) == 0 {
		_, err := fmt.Fprintln(out, "ok")
		return err
	}

	for _, issue := range issues {
		if _, err := fmt.Fprintf(out, "ERROR: %s\n%s\n", issue.Path, issue.Msg); err != nil {
			return err
		}
	}
	return fmt.Errorf("check failed")
}

func runList(args []string, wd string, out io.Writer, _ io.Writer) error {
	f, err := parseListArgs(args)
	if err != nil {
		return err
	}

	recs, err := record.LoadAll(wd)
	if err != nil {
		return err
	}

	return writeRows(out, record.Filtered(recs, f))
}

func runSearch(args []string, wd string, out io.Writer, _ io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("search query is required")
	}

	recs, err := record.LoadAll(wd)
	if err != nil {
		return err
	}

	return writeRows(out, record.Search(recs, strings.Join(args, " ")))
}

func runIndex(_ []string, wd string, out io.Writer, _ io.Writer) error {
	recs, err := record.LoadAll(wd)
	if err != nil {
		return err
	}
	if err := record.WriteIndexes(wd, recs); err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, "generated memory/indexes")
	return err
}

func runRelated(args []string, wd string, out io.Writer, _ io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("record id is required")
	}

	recs, err := record.LoadAll(wd)
	if err != nil {
		return err
	}

	cands, err := record.Related(recs, args[0])
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "Related candidates for %s:\n\n", args[0]); err != nil {
		return err
	}
	return writeCandRows(out, cands)
}

func runClose(args []string, wd string, out io.Writer, _ io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("record id is required")
	}

	id := args[0]
	vals, err := parseCloseArgs(args[1:])
	if err != nil {
		return err
	}

	rec, err := record.LoadByID(wd, id)
	if err != nil {
		return err
	}

	if _, ok := vals["Status"]; !ok {
		vals["Status"] = "CLOSED"
	}
	if err := record.UpdateFields(rec.Path, vals); err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "updated %s\n", rec.Path)
	return err
}

func runVersion(args []string, _ string, out io.Writer, _ io.Writer) error {
	info := buildinfo.Current()
	text := buildinfo.Render(info)
	for _, arg := range args {
		switch arg {
		case "--verbose":
			text = buildinfo.RenderVerbose(info)
		default:
			return fmt.Errorf("unknown option: %s", arg)
		}
	}

	_, err := io.WriteString(out, text)
	return err
}

func parseListArgs(args []string) (record.Filter, error) {
	var f record.Filter
	for i := 0; i < len(args); i++ {
		if i+1 >= len(args) {
			return f, fmt.Errorf("missing value for %s", args[i])
		}
		switch args[i] {
		case "--type":
			f.Type = args[i+1]
		case "--status":
			f.Status = args[i+1]
		case "--area":
			f.Area = args[i+1]
		case "--future":
			f.Future = args[i+1]
		default:
			return f, fmt.Errorf("unknown option: %s", args[i])
		}
		i++
	}
	return f, nil
}

func parseCloseArgs(args []string) (map[string]string, error) {
	vals := map[string]string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--verified":
			vals["Status"] = "VERIFIED"
		case "--resolved-by":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for %s", args[i])
			}
			vals["ResolvedBy"] = args[i+1]
			i++
		default:
			return nil, fmt.Errorf("unknown option: %s", args[i])
		}
	}
	return vals, nil
}

func writeRows(out io.Writer, recs []record.Rec) error {
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	for _, rec := range recs {
		if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", rec.ID, rec.Status(), blank(rec.Area()), rec.Title); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func writeCandRows(out io.Writer, cands []record.Candidate) error {
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	for _, cand := range cands {
		rec := cand.Rec
		if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", rec.ID, rec.Status(), blank(rec.Area()), rec.Title); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func blank(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}
