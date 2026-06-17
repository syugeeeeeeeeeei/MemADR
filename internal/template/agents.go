package template

import "strings"

func RenderAgentsSnippet() string {
	lines := []string{
		"## MemADR運用ポリシー",
		"",
		"開発知識レコードの管理には `memadr` を使用する。",
		"",
		"- `memory/` 配下のレコードは原則として日本語で記述する。",
		"- `memory/` のレコードは短く、結果と判断が分かる内容を優先する。",
		"- `memory/` に作業ログを書かない。",
		"- BUG、PROB、ADR、CHG、REV、SOL、SUP を新規作成する前に、`memadr search` と `memadr related` で既存レコードを確認する。",
		"- リポジトリ開始時は `memadr init` を実行し、`MEMADR_WORKFLOW.md` を読む。",
		"- バグ、構造問題、設計判断、変更、巻き戻し、再利用可能な解決策、無効化が見つかったら、対応するレコードを作成または更新する。",
		"- レコードを追加または更新した作業では、完了前に `memadr check` を実行する。",
		"- `memory/generated/` は手編集せず、`memadr index` で再生成する。",
		"- 各レコードファイルの状態は、そのファイル自身を正本として扱う。",
	}

	return strings.Join(lines, "\n")
}
