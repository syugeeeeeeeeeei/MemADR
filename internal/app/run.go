package app

import (
	"errors"
	"fmt"
	"io"

	"memadr/internal/mem"
	"memadr/internal/template"
)

func Run(args []string, wd string, out io.Writer, errOut io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: memadr <command>")
	}

	switch args[0] {
	case "init":
		return runInit(wd, out)
	case "new":
		return runNew(args[1:], wd, out)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func runInit(wd string, out io.Writer) error {
	if err := mem.Init(wd); err != nil {
		return err
	}

	_, err := fmt.Fprintln(out, "initialized memory/")
	return err
}

func runNew(args []string, wd string, out io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: memadr new <type> [title]")
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
