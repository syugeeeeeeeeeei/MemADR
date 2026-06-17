package app

import (
	"fmt"
	"io"

	"memadr/internal/mem"
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
