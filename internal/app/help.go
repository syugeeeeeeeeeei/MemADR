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
	Args     []inputDef
	Options  []optionDef
	Examples []string
	Run      func(args []string, wd string, out io.Writer, errOut io.Writer) error
}

type inputDef struct {
	Name      string
	Summary   string
	ValueNote string
	Values    func() []string
}

type optionDef struct {
	Name      string
	ValueName string
	Summary   string
	ValueNote string
	Values    func() []string
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
			Args: []inputDef{
				{
					Name:      "[command]",
					Summary:   "詳細を表示したいコマンド名",
					Values:    commandNames,
					ValueNote: "省略時は総合ヘルプを表示する",
				},
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
			Args: []inputDef{
				{
					Name:    "<type>",
					Summary: "作成するレコード種別",
					Values:  kindNames,
				},
				{
					Name:      "[title]",
					Summary:   "レコードのタイトル",
					ValueNote: "任意のタイトル文字列。省略時は `タイトル未設定`",
				},
			},
			Examples: []string{
				`memadr new bug "認証状態が壊れる"`,
				`memadr new adr "認証責務を server session 中心に整理する"`,
			},
			Run: runNew,
		},
		{
			Name:    "check",
			Summary: "レコード形式と参照整合性を検査する",
			Usage:   "memadr check",
			Details: []string{
				"Status と FutureRelevance の許可値を検査する",
				"必須フィールド不足と存在しない参照IDを検査する",
				"generated 配下の生成コメントも検査する",
			},
			Examples: []string{
				"memadr check",
			},
			Run: runCheck,
		},
		{
			Name:    "list",
			Summary: "レコードを条件付きで一覧表示する",
			Usage:   "memadr list [--type kind] [--status STATUS] [--area area] [--future value]",
			Details: []string{
				"未指定なら全レコードを表示する",
				"複数条件を同時に指定できる",
			},
			Options: []optionDef{
				{
					Name:      "--type",
					ValueName: "TYPE",
					Summary:   "表示対象のレコード種別で絞り込む",
					Values:    kindNames,
				},
				{
					Name:      "--status",
					ValueName: "STATUS",
					Summary:   "Status で絞り込む",
					Values:    mem.Statuses,
				},
				{
					Name:      "--area",
					ValueName: "AREA",
					Summary:   "Area で絞り込む",
					ValueNote: "任意のArea文字列",
				},
				{
					Name:      "--future",
					ValueName: "FUTURE",
					Summary:   "FutureRelevance で絞り込む",
					Values:    mem.FutureValues,
				},
			},
			Examples: []string{
				"memadr list",
				"memadr list --status VERIFIED",
				"memadr list --area auth",
			},
			Run: runList,
		},
		{
			Name:    "search",
			Summary: "タイトルとフィールドから文字列検索する",
			Usage:   "memadr search <query>",
			Details: []string{
				"タイトル、Status、Area、本文フィールド、関連IDを検索する",
			},
			Args: []inputDef{
				{
					Name:      "<query>",
					Summary:   "検索文字列",
					ValueNote: "任意の検索語句",
				},
			},
			Examples: []string{
				`memadr search "認証"`,
				`memadr search "server session"`,
			},
			Run: runSearch,
		},
		{
			Name:    "index",
			Summary: "generated 配下の集約Markdownを再生成する",
			Usage:   "memadr index",
			Details: []string{
				"active by-status by-area by-type open を再生成する",
			},
			Examples: []string{
				"memadr index",
			},
			Run: runIndex,
		},
		{
			Name:    "related",
			Summary: "指定レコードに関連しそうな候補を表示する",
			Usage:   "memadr related <record-id>",
			Details: []string{
				"Area 一致、明示リンク、単語一致をもとに候補を並べる",
			},
			Args: []inputDef{
				{
					Name:      "<record-id>",
					Summary:   "関連候補を見たいレコードID",
					ValueNote: "BUG-001 や ADR-001 のような既存レコードID",
				},
			},
			Examples: []string{
				"memadr related BUG-012",
			},
			Run: runRelated,
		},
		{
			Name:    "close",
			Summary: "レコードの Status と関連フィールドを更新する",
			Usage:   "memadr close <record-id> [--resolved-by CHG-001] [--verified]",
			Details: []string{
				"`--verified` を付けると Status を VERIFIED にする",
				"未指定時は Status を CLOSED にする",
			},
			Args: []inputDef{
				{
					Name:      "<record-id>",
					Summary:   "更新対象のレコードID",
					ValueNote: "BUG-001 や ADR-001 のような既存レコードID",
				},
			},
			Options: []optionDef{
				{
					Name:      "--resolved-by",
					ValueName: "CHG-ID",
					Summary:   "ResolvedBy フィールドを更新する",
					ValueNote: "`CHG-001` のような変更ID",
				},
				{
					Name:      "--verified",
					Summary:   "Status を VERIFIED にする",
					ValueNote: "値なし",
				},
			},
			Examples: []string{
				"memadr close BUG-012 --resolved-by CHG-021 --verified",
			},
			Run: runClose,
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
	b.WriteString("  memadr help list\n")
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

	if len(cmd.Args) > 0 {
		b.WriteString("Arguments:\n")
		b.WriteString(renderInputTable(cmd.Args))
		b.WriteString("\n")
	}

	if len(cmd.Options) > 0 {
		b.WriteString("Options:\n")
		b.WriteString(renderOptionTable(cmd.Options))
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

func renderInputTable(args []inputDef) string {
	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	for _, arg := range args {
		fmt.Fprintf(tw, "  %s\t%s\n", arg.Name, arg.Summary)
		if line := renderValueLine(arg.ValueNote, arg.Values); line != "" {
			fmt.Fprintf(tw, "  \t%s\n", line)
		}
	}
	_ = tw.Flush()
	return b.String()
}

func renderOptionTable(opts []optionDef) string {
	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	for _, opt := range opts {
		fmt.Fprintf(tw, "  %s\t%s\n", opt.Label(), opt.Summary)
		if line := renderValueLine(opt.ValueNote, opt.Values); line != "" {
			fmt.Fprintf(tw, "  \t%s\n", line)
		}
	}
	_ = tw.Flush()
	return b.String()
}

func renderValueLine(note string, vals func() []string) string {
	if vals != nil {
		items := vals()
		if len(items) > 0 {
			return "値: " + strings.Join(items, ", ")
		}
	}
	if note != "" {
		return "値: " + note
	}
	return ""
}

func (opt optionDef) Label() string {
	if opt.ValueName == "" {
		return opt.Name
	}
	return opt.Name + " <" + opt.ValueName + ">"
}

func kindNames() []string {
	kinds := mem.Kinds()
	out := make([]string, 0, len(kinds))
	for _, kind := range kinds {
		out = append(out, kind.Name)
	}
	return out
}

func commandNames() []string {
	cmds := commandDefs()
	out := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		out = append(out, cmd.Name)
	}
	return out
}
