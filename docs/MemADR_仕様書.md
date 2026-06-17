# Go製CLIツール仕様書：MemADR

## 1. 概要

`MemADR` は、LLM向けの開発知識レコードを管理するためのGo製CLIツールである。

CLIコマンド名は `memadr` とする。

本ツールは、作業ログそのものではなく、次回以降の開発判断に必要な「問題」「原因」「設計判断」「変更」「解決策」「無効化された前提」を、短く構造化されたMarkdownレコードとして管理する。

正本は常に `memory/` 配下のMarkdownファイルとし、`memadr` はその上に薄く乗る補助ツールとする。

## 2. 目的

本ツールの目的は、以下である。

* LLMに渡すコンテキスト量を削減する
* 開発判断に必要な知識を短く残す
* Git管理可能なMarkdownレコードを生成・検査する
* レコードの一覧、検索、関連候補提示を補助する
* `generated/` 配下の集約ファイルを再生成する

本ツールは、Issue管理システム、Git代替、タスク管理ツール、完全自動分類システムではない。

## 3. 基本方針

### 3.1 Markdownを正本にする

すべてのレコードの正本はMarkdownファイルである。

SQLite、検索インデックス、一覧ファイル、集約ファイルはすべて生成物とし、Markdownから再生成できるものとする。

### 3.2 Gitを再発明しない

以下はGitに任せる。

* 履歴管理
* 差分確認
* ブランチ管理
* マージ
* レビュー
* いつ変更したかの追跡

`memadr` はGitの上に乗る補助ツールに限定する。

### 3.3 CLIは小さく保つ

初期版では、以下に責務を限定する。

* 初期ディレクトリ生成
* テンプレート生成
* 形式チェック
* 一覧表示
* 簡易検索
* 集約ファイル生成
* 参照関係の検査

LLM連携、ベクトル検索、GUI、大規模なDB管理は初期スコープに含めない。

## 4. 提供形態

### 4.1 実装言語

Goを使用する。

理由は以下である。

* 単体バイナリとして配布しやすい
* Windows、macOS、Linuxに対応しやすい
* 利用者側でGo環境を必要としない
* ファイル操作、CLI、文字列処理に向いている
* 小さな保守可能なCLIとして実装しやすい

### 4.2 配布形式

GitHub ReleasesでOS別バイナリを配布する。

想定する配布物は以下である。

```text
memadr_windows_amd64.exe
memadr_darwin_arm64
memadr_darwin_amd64
memadr_linux_amd64
```

利用者は対象OSのバイナリを取得し、PATHに配置するだけで利用できる。

### 4.3 利用者側の必要環境

必須：

```text
memadr
```

推奨：

```text
Git
```

任意：

```text
ripgrep
```

`ripgrep` が存在する場合は検索に利用してもよい。ただし、存在しない場合でも `memadr` 内蔵検索で動作する。

## 5. ディレクトリ構成

`memadr init` により、以下の構成を作成する。

```text
memory/
├─ bugs/
├─ problems/
├─ decisions/
├─ changes/
├─ reversions/
├─ solutions/
├─ supersessions/
└─ generated/
```

各ディレクトリの意味は以下である。

| ディレクトリ           | 種別   | 内容          |
| ---------------- | ---- | ----------- |
| `bugs/`          | BUG  | バグ、不具合、失敗   |
| `problems/`      | PROB | 根本問題、構造的課題  |
| `decisions/`     | ADR  | 設計判断        |
| `changes/`       | CHG  | 実際に行った変更    |
| `reversions/`    | REV  | 巻き戻し、撤回     |
| `solutions/`     | SOL  | 再利用可能な解決策   |
| `supersessions/` | SUP  | 過去判断や実装の無効化 |
| `generated/`     | 生成物  | 一覧、索引、検索DB  |

`generated/` 配下は手動編集しない。

## 6. レコード種別

### 6.1 BUG

バグ、不具合、失敗を記録する。

必須フィールド：

```text
Status
Area
Symptom
```

推奨フィールド：

```text
Cause
Fix
Verification
FutureRelevance
Decision
ResolvedBy
Related
```

例：

```markdown
# BUG-012: 認証状態が画面遷移後に壊れる

Status: VERIFIED
Area: auth
Symptom: 画面遷移後にログイン状態が不安定になる
Cause: 認証責務がclient/server/storageに分散していた
Fix: 認証責務をserver session中心に再整理した
Verification: ログイン、ログアウト、期限切れ、再ログインを確認
FutureRelevance: watch
Decision: ADR-007
ResolvedBy: CHG-021
```

### 6.2 PROB

構造的な問題や未解決課題を記録する。

例：

```markdown
# PROB-001: 認証状態管理が不安定

Status: UNRESOLVED
Area: auth
Impact: 認証バグが複数箇所で再発する
Finding: localStorage、cookie、server sessionの責務が混在している
CurrentDirection: 部分修正ではなく認証基盤を再設計する
Related: ADR-007
```

### 6.3 ADR

設計判断を記録する。

例：

```markdown
# ADR-007: 認証責務をserver session中心に整理する

Status: ACCEPTED
Context: BUG-012の原因が単一不具合ではなく、認証責務の分散だった
Decision: 認証状態管理をserver session中心に統一する
Rejected: 既存構造への局所パッチ継続
Consequence: 認証関連の依存方向とディレクトリ構造を整理する
Related: BUG-012
```

### 6.4 CHG

実際に行った変更を記録する。

例：

```markdown
# CHG-021: 認証関連ディレクトリを再編成

Status: SHIPPED
Reason: ADR-007
Change: 認証責務をserver session中心に集約した
FilesChanged: src/auth/**, src/session/**, src/middleware/**
Verification: 認証フローの主要ケースを確認
```

### 6.5 REV

過去の変更を戻した場合に記録する。

例：

```markdown
# REV-003: 新キャッシュ方式をロールバック

Status: CLOSED
Reverted: CHG-008
Reason: 一部環境でキャッシュ不整合が発生した
Result: 旧キャッシュ方式へ戻した
FollowUp: 再導入する場合はSOL-007の検証条件を満たす
```

### 6.6 SOL

再利用可能な解決パターンを記録する。

例：

```markdown
# SOL-004: JWT期限判定はUTCに統一する

Status: ACTIVE
ProblemPattern: JWT期限判定が環境差で壊れる
Solution: 時刻処理をUTCに統一する
AppliesTo: auth, session, token-expiry
FutureRelevance: reusable
Related: BUG-001, BUG-018
```

### 6.7 SUP

過去の実装、判断、変更を無効化または置き換えたことを記録する。

例：

```markdown
# SUP-003: 旧認証実装を廃止

Status: ACTIVE
Supersedes: CHG-004, CHG-009, BUG-012, BUG-018
Reason: 同系統の問題が再発し、部分修正では保守性が悪化した
NewBaseline: ADR-007
```

## 7. ID規則

各レコードは、以下の形式のIDを持つ。

```text
BUG-001
PROB-001
ADR-001
CHG-001
REV-001
SOL-001
SUP-001
```

ファイル名はIDに一致させる。

例：

```text
memory/bugs/BUG-012.md
memory/decisions/ADR-007.md
memory/changes/CHG-021.md
```

IDは種別ごとに連番とする。

## 8. Status規則

使用可能なStatusは以下とする。

```text
OPEN
INVESTIGATING
UNRESOLVED
PROPOSED
ACCEPTED
FIXED
VERIFIED
SHIPPED
CLOSED
SUPERSEDED
ARCHIVED
ACTIVE
```

`FIXED` と `VERIFIED` は区別する。

| Status       | 意味                 |
| ------------ | ------------------ |
| `FIXED`      | 修正は入ったが、確認は完了していない |
| `VERIFIED`   | 修正後の挙動確認が完了している    |
| `SUPERSEDED` | 後続の判断や変更により置き換えられた |
| `ARCHIVED`   | 通常参照しない過去情報        |
| `ACTIVE`     | 現在も有効な解決策または置き換え情報 |

## 9. FutureRelevance規則

`FutureRelevance` には以下の値を使用する。

```text
ignore
watch
reusable
```

| 値          | 意味                  |
| ---------- | ------------------- |
| `ignore`   | 解決済みで通常は考慮不要        |
| `watch`    | 再発や副作用に注意           |
| `reusable` | 同種問題の解決策として再利用価値がある |

## 10. コマンド仕様

### 10.1 `memadr init`

`memory/` ディレクトリ構成を作成する。

実行例：

```bash
memadr init
```

作成対象：

```text
memory/
memory/bugs/
memory/problems/
memory/decisions/
memory/changes/
memory/reversions/
memory/solutions/
memory/supersessions/
memory/generated/
```

既に存在するディレクトリは上書きしない。

### 10.2 `memadr new`

新規レコードを作成する。

実行例：

```bash
memadr new bug
memadr new adr
memadr new chg
memadr new sol
```

タイトル指定：

```bash
memadr new bug "認証状態が画面遷移後に壊れる"
```

生成例：

```text
memory/bugs/BUG-012.md
```

IDは既存ファイルを走査して自動採番する。

対応種別：

```text
bug
prob
adr
chg
rev
sol
sup
```

### 10.3 `memadr check`

レコード形式を検査する。

実行例：

```bash
memadr check
```

検査内容：

* ファイル名とIDの一致
* ID形式の妥当性
* 必須フィールドの存在
* Statusの許可値
* FutureRelevanceの許可値
* 参照IDの存在
* `generated/` 配下の手動編集禁止コメント
* 同一IDの重複
* 不明なレコード種別

異常がある場合は、終了コード `1` を返す。

### 10.4 `memadr list`

レコードを一覧表示する。

実行例：

```bash
memadr list
memadr list --type bug
memadr list --status VERIFIED
memadr list --area auth
memadr list --future watch
```

標準出力例：

```text
BUG-012  VERIFIED  auth  認証状態が画面遷移後に壊れる
ADR-007  ACCEPTED  auth  認証責務をserver session中心に整理する
CHG-021  SHIPPED   auth  認証関連ディレクトリを再編成
```

### 10.5 `memadr search`

Markdownレコードを検索する。

実行例：

```bash
memadr search "認証"
memadr search "server session"
memadr search "JWT 期限"
```

検索対象：

* タイトル
* Status
* Area
* 本文フィールド
* 関連ID

初期実装では、Markdownファイルを直接走査する。

`ripgrep` が利用可能な環境では、内部的に利用してもよい。ただし、`ripgrep` は必須依存にしない。

### 10.6 `memadr index`

`generated/` 配下の集約ファイルを再生成する。

実行例：

```bash
memadr index
```

生成対象：

```text
memory/generated/active.md
memory/generated/by-status.md
memory/generated/by-area.md
memory/generated/by-type.md
memory/generated/open.md
```

各生成ファイルの先頭には以下を付与する。

```markdown
<!-- This file is generated. Do not edit manually. -->
```

### 10.7 `memadr related`

指定レコードに関連しそうな候補を表示する。

実行例：

```bash
memadr related BUG-012
```

初期実装では、以下をもとに候補を表示する。

* 同じArea
* タイトル内の共通語
* Related、ResolvedBy、Decisionなどの明示リンク
* Symptom、Cause、Fix、Decision内の単語一致

出力例：

```text
Related candidates for BUG-012:

ADR-007  ACCEPTED  auth  認証責務をserver session中心に整理する
CHG-021  SHIPPED   auth  認証関連ディレクトリを再編成
SOL-004  ACTIVE    auth  JWT期限判定はUTCに統一する
```

初期版では、関連候補を自動で書き込まない。

### 10.8 `memadr close`

BUGなどのStatusを更新する補助コマンド。

実行例：

```bash
memadr close BUG-012 --resolved-by CHG-021 --verified
```

処理内容：

* `Status: VERIFIED` に更新
* `ResolvedBy: CHG-021` を追加または更新
* 既存本文は可能な限り保持する

自動更新対象は、明確に指定されたフィールドに限定する。

## 11. 設定ファイル

設定ファイルは任意とする。

配置場所：

```text
.memadr.yml
```

例：

```yaml
memory_dir: memory
language: ja
id_padding: 3
generated_dir: memory/generated
```

初期版では、設定ファイルがなくても動作する。

## 12. 出力形式

標準出力は、人間が読みやすいテキスト形式を基本とする。

将来拡張として、JSON出力を追加できる。

例：

```bash
memadr list --json
memadr search "auth" --json
```

ただし、初期版ではJSON出力は必須ではない。

## 13. 終了コード

| 終了コード | 意味               |
| ----: | ---------------- |
|   `0` | 正常終了             |
|   `1` | 検査エラー、入力エラー      |
|   `2` | ファイルシステムエラー      |
|   `3` | 未対応コマンド、未対応オプション |

## 14. エラーメッセージ方針

エラーメッセージは、以下を満たす。

* 対象ファイルを示す
* 問題のあるフィールドを示す
* 修正方針を短く示す

例：

```text
ERROR: memory/bugs/BUG-012.md
Invalid Status: DONE
Allowed values: OPEN, INVESTIGATING, FIXED, VERIFIED, CLOSED, ARCHIVED
```

## 15. 初期実装スコープ

初期版で実装する機能は以下である。

```text
memadr init
memadr new
memadr check
memadr list
memadr search
memadr index
memadr related
```

初期版で実装しない機能は以下である。

```text
GUI
Web UI
ベクトル検索
LLM API連携
SQLite FTS
自動分類
Git操作の代替
Issue管理
```

## 16. 将来拡張

運用上の必要性が確認できた場合、以下を追加する。

### 16.1 SQLite FTS

生成物として以下を作成する。

```text
memory/generated/memory.sqlite
```

用途：

* 高速全文検索
* 状態別集計
* Area別集計
* 逆引き
* 関連候補提示

SQLiteは正本にしない。

### 16.2 JSON出力

LLMや他ツールと連携しやすくするため、JSON出力を追加する。

例：

```bash
memadr list --json
memadr related BUG-012 --json
```

### 16.3 LLM向けコンテキスト出力

指定条件に一致するレコードだけを、LLMに渡しやすい短いMarkdownとして出力する。

例：

```bash
memadr export-context --area auth --status ACTIVE
memadr export-context --related BUG-012
```

## 17. 非目標

本ツールでは、以下を行わない。

* Git履歴の代替
* Issue管理の代替
* タスク管理
* チャットログ保存
* コマンド実行履歴の保存
* LLMの推論過程の保存
* 大規模なナレッジグラフ管理
* 全自動の分類・要約

## 18. 実装構成案

Goプロジェクトの構成は以下とする。

```text
memadr/
├─ cmd/
│  └─ memadr/
│     └─ main.go
├─ internal/
│  ├─ app/
│  ├─ record/
│  ├─ parser/
│  ├─ template/
│  ├─ checker/
│  ├─ indexer/
│  ├─ searcher/
│  └─ fsutil/
├─ templates/
│  ├─ bug.md
│  ├─ prob.md
│  ├─ adr.md
│  ├─ chg.md
│  ├─ rev.md
│  ├─ sol.md
│  └─ sup.md
├─ testdata/
├─ go.mod
├─ README.md
└─ .goreleaser.yaml
```

## 19. パーサー仕様

Markdownの先頭見出しをタイトルとして扱う。

例：

```markdown
# BUG-012: 認証状態が画面遷移後に壊れる
```

以下の形式を構造化フィールドとして扱う。

```markdown
Status: VERIFIED
Area: auth
Symptom: 画面遷移後にログイン状態が不安定になる
```

フィールド値は1行を基本とする。

初期版では、複数行フィールドは必須対応しない。

## 20. 受け入れ条件

初期版は、以下を満たした場合に完了とする。

* `memadr init` で標準ディレクトリを作成できる
* `memadr new bug` で連番付きBUGファイルを作成できる
* `memadr new adr` で連番付きADRファイルを作成できる
* `memadr check` で不正なStatusを検出できる
* `memadr check` で存在しない参照IDを検出できる
* `memadr list --status VERIFIED` が動作する
* `memadr list --area auth` が動作する
* `memadr search "認証"` が動作する
* `memadr index` で `generated/` 配下の一覧を再生成できる
* Windows、macOS、Linux向けのバイナリを作成できる
* 利用者側にGoのインストールを要求しない

## 21. 暫定結論

`MemADR` は、Go製の単体CLIとして実装する。

初期版では、Markdown正本、短いテンプレート、形式チェック、一覧生成、検索補助に限定する。

SQLite、LLM連携、ベクトル検索は初期版では導入せず、実際の運用で必要性が明確になった段階で追加する。
