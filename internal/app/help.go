package app

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"memadr/internal/mem"
)

type commandDef struct {
	Name     string
	Summary  string
	Usage    string
	Details  []string
	Examples []string
	Run      func(args []string, wd string, out io.Writer, errOut io.Writer) error
}

func commandDefs() []commandDef {
	return []commandDef{
		{
			Name:    "help",
			Summary: "CLIヘルプを表示する",
			Usage:   "memadr help [command]",
			Details: []string{
				"引数なしでは全体ヘルプを表示する",
				"`memadr help new` のように個別コマンドの詳細も表示できる",
			},
			Examples: []string{
				"memadr help",
				"memadr help init",
				"memadr help new",
			},
			Run: runHelp,
		},
		{
			Name:    "init",
			Summary: "memory/ の標準ディレクトリ構成を作成する",
			Usage:   "memadr init",
			Details: []string{
				"既存ディレクトリは壊さずに不足分だけ作成する",
				"`memory/generated/` も同時に用意する",
			},
			Examples: []string{
				"memadr init",
			},
			Run: runInit,
		},
		{
			Name:    "new",
			Summary: "新しいメモリレコードをテンプレート付きで作成する",
			Usage:   "memadr new <type> [title]",
			Details: []string{
				"種別ごとに連番IDを自動採番する",
				"タイトルを省略すると `タイトル未設定` を使う",
				"作成先ディレクトリは種別から自動で決まる",
			},
			Examples: []string{
				`memadr new bug "認証状態が壊れる"`,
				`memadr new adr "認証責務を server session 中心に整理する"`,
			},
			Run: runNew,
		},
	}
}

func findCommand(name string) (commandDef, bool) {
	for _, cmd := range commandDefs() {
		if cmd.Name == name {
			return cmd, true
		}
	}
	return commandDef{}, false
}

func runHelp(args []string, _ string, out io.Writer, _ io.Writer) error {
	if len(args) == 0 {
		_, err := io.WriteString(out, renderGeneralHelp())
		return err
	}

	cmd, ok := findCommand(args[0])
	if !ok {
		return fmt.Errorf("unknown command for help: %s", args[0])
	}

	_, err := io.WriteString(out, renderCommandHelp(cmd))
	return err
}

func renderGeneralHelp() string {
	var b strings.Builder

	b.WriteString("MemADR\n")
	b.WriteString("LLM向けの開発知識レコードを、短いMarkdown正本として管理するCLIです。\n\n")

	b.WriteString("Usage:\n")
	b.WriteString("  memadr <command> [arguments]\n\n")

	b.WriteString("Commands:\n")
	b.WriteString(renderCommandTable(commandDefs()))
	b.WriteString("\n")

	b.WriteString("Quick start:\n")
	b.WriteString("  1. memadr init\n")
	b.WriteString("  2. memadr new bug \"認証状態が壊れる\"\n")
	b.WriteString("  3. memadr new adr \"認証責務を整理する\"\n\n")

	b.WriteString("Record types:\n")
	b.WriteString(renderKindTable())
	b.WriteString("\n")

	b.WriteString("Examples:\n")
	b.WriteString("  memadr help new\n")
	b.WriteString("  memadr init\n")
	b.WriteString("  memadr new bug \"認証状態が壊れる\"\n")
	b.WriteString("  memadr new adr \"認証責務を整理する\"\n")

	return b.String()
}

func renderCommandHelp(cmd commandDef) string {
	var b strings.Builder

	b.WriteString(cmd.Name)
	b.WriteString("\n")
	b.WriteString(cmd.Summary)
	b.WriteString("\n\n")

	b.WriteString("Usage:\n")
	b.WriteString("  ")
	b.WriteString(cmd.Usage)
	b.WriteString("\n\n")

	if len(cmd.Details) > 0 {
		b.WriteString("Details:\n")
		for _, line := range cmd.Details {
			b.WriteString("  - ")
			b.WriteString(line)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	if cmd.Name == "new" {
		b.WriteString("Supported types:\n")
		b.WriteString(renderKindTable())
		b.WriteString("\n")
	}

	if len(cmd.Examples) > 0 {
		b.WriteString("Examples:\n")
		for _, ex := range cmd.Examples {
			b.WriteString("  ")
			b.WriteString(ex)
			b.WriteString("\n")
		}
	}

	return b.String()
}

func renderCommandTable(cmds []commandDef) string {
	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	for _, cmd := range cmds {
		fmt.Fprintf(tw, "  %s\t%s\t%s\n", cmd.Name, cmd.Usage, cmd.Summary)
	}
	_ = tw.Flush()
	return b.String()
}

func renderKindTable() string {
	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	for _, kind := range mem.Kinds() {
		fmt.Fprintf(tw, "  %s (%s)\t%s\n", kind.Name, kind.Prefix, kind.Desc)
	}
	_ = tw.Flush()
	return b.String()
}
