package template

import (
	"strings"

	"memadr/internal/mem"
)

func RenderWorkflowGuide() string {
	lines := []string{
		"# MemADR Workflow Guide",
		"",
		"このファイルは `memadr init` が生成する、人間およびLLM向けの利用導線です。",
		"",
		"## 目的",
		"",
		"- `memory/` を短いMarkdown正本として維持する",
		"- 作業ログではなく、問題、判断、変更、解決策を残す",
		"- 人間とLLMが次の操作に迷わないように、基本ワークフローを揃える",
		"",
		"## 人間向け導線",
		"",
		"1. `memadr init` で標準構成とこの案内を作成する",
		"2. 問題を見つけたら `memadr new bug \"...\"` か `memadr new prob \"...\"` で起票する",
		"3. 設計判断が必要なら `memadr new adr \"...\"` を作る",
		"4. 実装後は `memadr new chg \"...\"` で変更記録を残す",
		"5. `memadr check` で形式と参照を検査する",
		"6. `memadr list` `memadr search` `memadr related` で既存記録を確認する",
		"7. 必要なら `memadr close BUG-001 --resolved-by CHG-001 --verified` で状態を更新する",
		"8. `memadr index` で `memory/generated/` を再生成する",
		"",
		"## LLM向け導線",
		"",
		"- 新しい記録を作る前に `memadr search` と `memadr related` で重複や近い判断を探す",
		"- `memory/` には作業ログではなく、短い結論だけを書く",
		"- BUGは症状、原因、修正、検証に絞る",
		"- ADRはなぜその方針を選んだかを書く",
		"- CHGは何を変えたかと検証結果を書く",
		"- 状態集約は手書きせず、`memadr index` に任せる",
		"",
		"## 代表ユースケース",
		"",
		"### バグ修正",
		"",
		"1. `memadr new bug \"認証状態が壊れる\"`",
		"2. 必要なら `memadr new adr \"認証責務を整理する\"`",
		"3. 実装後に `memadr new chg \"認証関連ディレクトリを再編成\"`",
		"4. `memadr close BUG-001 --resolved-by CHG-001 --verified`",
		"",
		"### 既存知識の確認",
		"",
		"1. `memadr list --status OPEN` で未完了項目を確認する",
		"2. `memadr search \"認証\"` で関連記録を探す",
		"3. `memadr related ADR-001` で近い判断を探す",
		"",
		"## 補足",
		"",
		"- この案内は必要なら人間が追記してよいが、`memory/generated/` のような生成物ではない",
		"- 詳細なCLIオプションは `memadr help` と `memadr help <command>` を使う",
		"- ワークフローの参照先ファイル名は `" + mem.WorkflowGuideFile + "`",
		"",
	}

	return strings.Join(lines, "\n")
}
