package template

import "strings"

func RenderAgentsSnippet() string {
	lines := []string{
		"# AGENTS.md",
		"",
		"このリポジトリでは、MemADRを開発判断メモリとして使用する。",
		"",
		"作業するエージェントは、必ず `MEMADR_WORKFLOW.md` を読み、その内容に従うこと。",
		"",
		"特に、次を守ること。",
		"",
		"- `memory/` に作業ログを書かない",
		"- 現在も価値がある情報と、過去情報を区別する",
		"- 古い判断、削除済み機能、無効化済み情報を現在有効な前提として扱わない",
		"- MemADRレコードを追加または更新した場合は、完了前に `memadr check` と `memadr index` を実行する",
	}

	return strings.Join(lines, "\n")
}
