# 企画書：MemADR

## LLM向けコンテキスト削減を目的とした開発知識レコード運用

## 1. 背景

現在、LLMエージェントに開発作業を行わせるたびに、以下のような情報をMarkdownログとして逐次保存している。

* 依頼したタスク
* LLMが行った作業の流れ
* 調査内容
* 実行したコマンド
* 修正した内容
* 試行錯誤
* 作業後の要約

この方式では作業履歴は残るものの、開発が進むほどログが膨大になり、LLMに渡すコンテキスト量が増え続ける。

その結果、以下の問題が発生する。

* トークン消費が増える
* LLMに渡す文脈が肥大化する
* 古い情報や無効な前提が混入する
* 必要な情報がログの中に埋もれる
* 要約処理のたびに情報欠落や誤解が発生する
* 人間が要約の正確性を検証しにくい

## 2. 主目的

本企画の主目的は、以下である。

> LLMに渡すトークン量とコンテキスト量を削減しながら、開発判断に必要な知識を失わないこと。

単にログを整理することではない。

目的は、LLMと人間が次回以降の開発で必要とする情報だけを、短く、構造化され、Git管理可能なMarkdownレコードとして残すことである。

## 3. ツール名称

本企画で扱うツール名は `MemADR` とする。

CLIコマンド名は `memadr` とする。

名称の意味は、開発知識の記録であるMemoryと、設計判断を記録するADRの考え方を組み合わせたものである。

`MemADR` は、ADRだけを管理するツールではない。BUG、PROB、ADR、CHG、REV、SOL、SUPを含む、LLM向けの開発知識レコードを管理する薄いCLIツールである。

## 4. 管理したい情報

管理対象は、作業手順ではなく、開発上の判断や結論である。

具体的には以下を管理する。

* どのような問題があったか
* その問題の原因は何だったか
* どのように修正したか
* 解決済みか、未解決か
* 今後も注意すべきか
* どの設計判断が現在も有効か
* どの設計判断が無効化されたか
* なぜ大規模リファクタリングに至ったか
* どの変更がどの問題を解決したか
* どの解決策が再利用可能か

## 5. 使用言語

基本的に、記録本文の使用言語は日本語とする。

ただし、以下のようなフィールド名は英語を許容する。

* Status
* Area
* Symptom
* Cause
* Fix
* Verification
* FutureRelevance
* Decision
* Context
* Rejected
* Consequence
* Related
* ResolvedBy
* SupersededBy

理由は、LLM、CLI、検索、静的解析との相性を考慮するためである。

本文、説明、判断理由、原因、修正内容は日本語で記述する。

## 6. 現状方式の課題

### 6.1 ログの肥大化

LLMのタスクログやウォークスルーをすべて保存すると、Markdownファイルが大量に増える。

開発期間が長くなるほど、以下が混在する。

* 現在も有効な情報
* すでに無効な前提
* 途中で捨てた案
* 失敗した試行
* 解決済みのバグ
* 未解決の問題
* 古い設計判断

これにより、LLMに渡すコンテキストが膨張し、必要な情報に到達しにくくなる。

### 6.2 要約による情報欠落

定期的にログを要約して現在状態ドキュメントを作る方法には限界がある。

問題点は以下である。

* 要約自体に大量のトークンを消費する
* LLMがログを誤解する可能性がある
* 要約時に重要情報が抜け落ちる可能性がある
* 人間が元ログとの差分を確認しにくい
* 要約を繰り返すほど情報が劣化する

### 6.3 文脈の抽出困難

ログには作業の流れは書かれているが、後から必要になるのは作業手順そのものではない。

後から知りたいのは以下である。

* なぜこの設計になったのか
* どの問題が原因でこの変更をしたのか
* なぜ部分修正ではなく作り直しになったのか
* どの判断が失敗だったのか
* 現在どの判断が有効なのか

単純なログ保存では、これらの文脈が埋もれる。

## 7. 基本方針

### 7.1 作業ログではなく知識レコードを残す

残すべきものは、作業手順ではなく結論である。

悪い例：

* npm testを実行した
* auth.tsを開いた
* 仮説を検証した
* middlewareを修正した
* 再度テストした

良い例：

* Cause: 認証状態の責務がclient/server/storageに分散していた
* Fix: server session中心に責務を整理した
* Verification: ログイン、ログアウト、期限切れ、再ログインを確認した
* FutureRelevance: watch

### 7.2 トークン削減を優先する

各レコードは、LLMにそのまま渡しても負担にならない長さにする。

目安は以下である。

* BUG: 5〜10行
* PROB: 5〜10行
* ADR: 5〜10行
* CHG: 5〜10行
* REV: 5〜10行
* SOL: 5〜15行
* SUP: 5〜10行

作業ログ全文、長いコード差分、コマンド履歴、試行錯誤の詳細は原則として `memory/` 配下に残さない。

### 7.3 1ファイルは1つの意味単位にする

1ファイルに何でも書かない。

以下のように分ける。

* BUG: バグ、不具合、失敗
* PROB: 根本問題、未解決課題、構造的制約
* ADR: 設計判断
* CHG: 実際に行った変更
* REV: 巻き戻し、撤回
* SOL: 再利用可能な解決パターン
* SUP: 過去の実装や判断の無効化、置き換え

### 7.4 状態は各ファイル自身に持たせる

`active.md` のような一元管理ファイルを手動更新しない。

状態は各レコードに書く。

```markdown
Status: VERIFIED
```

一覧やactive表示は、必要に応じて `memadr` や検索で生成する。

これにより、状態変更時に複数ファイルを修正する必要を避ける。

### 7.5 Gitを再発明しない

履歴管理、分岐、マージ、差分、レビューはGitに任せる。

`MemADR` を作るとしても、Gitの上に薄く乗る補助ツールに留める。

`MemADR` の責務は以下に限定する。

* テンプレート生成
* 形式チェック
* 一覧表示
* 検索
* 関連候補の提示
* レコード種別の分割補助
* 生成ファイルの再生成

## 8. レコード種別

### 8.1 BUG

バグそのものを記録する。

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

BUGには詳細な作業手順を書かない。

書くのは以下に限定する。

* 症状
* 原因
* 修正
* 検証
* 今後の関連性
* 関連する判断や変更

### 8.2 PROB

根本問題や構造的な未解決課題を記録する。

BUGより抽象度が高い。

```markdown
# PROB-001: 認証状態管理が不安定

Status: UNRESOLVED
Area: auth
Impact: 認証バグが複数箇所で再発する
Finding: localStorage、cookie、server sessionの責務が混在している
CurrentDirection: 部分修正ではなく認証基盤を再設計する
Related: ADR-007
```

PROBは、単発バグではなく、複数の問題の背景にある構造的原因を扱う。

### 8.3 ADR

設計判断を記録する。

ADRは「何を実装したか」ではなく、「なぜその方向を選んだか」を残す。

```markdown
# ADR-007: 認証責務をserver session中心に整理する

Status: ACCEPTED
Context: BUG-012の原因が単一不具合ではなく、認証責務の分散だった
Decision: 認証状態管理をserver session中心に統一する
Rejected: 既存構造への局所パッチ継続
Consequence: 認証関連の依存方向とディレクトリ構造を整理する
```

ADRを書くべきケースは以下である。

* 認証方式を変える
* DBや永続化方式を変える
* ディレクトリ構造を大きく変える
* 状態管理方式を変える
* 外部サービスを採用または廃止する
* 大規模リファクタリングの方針を決める
* 将来また議論になりそうな選択を固定する

ADRにしないケースは以下である。

* nullチェック追加
* typo修正
* 単純なUI文言変更
* 小さなバグ修正
* import整理

### 8.4 CHG

実際に行った変更を記録する。

```markdown
# CHG-021: 認証関連ディレクトリを再編成

Status: SHIPPED
Reason: ADR-007
Change: 認証責務をserver session中心に集約した
FilesChanged: src/auth/**, src/session/**, src/middleware/**
Verification: 認証フローの主要ケースを確認
```

CHGは「何を変えたか」を記録する。

設計判断の理由はADRに書く。

### 8.5 REV

過去の変更を戻した場合に記録する。

```markdown
# REV-003: 新キャッシュ方式をロールバック

Status: CLOSED
Reverted: CHG-008
Reason: 一部環境でキャッシュ不整合が発生した
Result: 旧キャッシュ方式へ戻した
FollowUp: 再導入する場合はSOL-007の検証条件を満たす
```

REVは、何を戻したか、なぜ戻したか、再導入条件があるかを残す。

### 8.6 SOL

再利用可能な解決パターンを記録する。

```markdown
# SOL-004: JWT期限判定はUTCに統一する

Status: ACTIVE
ProblemPattern: JWT期限判定が環境差で壊れる
Solution: 時刻処理をUTCに統一する
AppliesTo: auth, session, token-expiry
Related: BUG-001, BUG-018
```

同じ原因や似た問題が複数回発生した場合、BUGの情報をSOLに昇格する。

### 8.7 SUP

過去の実装、判断、変更を無効化または置き換えたことを記録する。

```markdown
# SUP-003: 旧認証実装を廃止

Status: ACTIVE
Supersedes: CHG-004, CHG-009, BUG-012, BUG-018
Reason: 同系統の問題が再発し、部分修正では保守性が悪化した
NewBaseline: ADR-007
```

SUPは、大規模リファクタリングや作り直しに至った文脈を残すために重要である。

## 9. 状態管理

各レコードはStatusを持つ。

例：

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

特に、`FIXED` と `VERIFIED` は分ける。

* FIXED: 修正は入ったが、確認が完了していない
* VERIFIED: 修正後の挙動確認が完了している

## 10. FutureRelevance

BUGやSOLには、今後どの程度気にするべきかを明示する。

```text
ignore
watch
reusable
```

意味は以下である。

* ignore: 解決済みで通常は考慮不要
* watch: 再発や副作用に注意
* reusable: 同種問題の解決策として再利用価値がある

例：

```markdown
FutureRelevance: watch
```

これにより、LLMが過去の解決済み問題に過剰に引っ張られることを防ぐ。

## 11. ファイル構成

最小構成：

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

例：

```text
memory/
├─ bugs/
│  └─ BUG-012.md
├─ problems/
│  └─ PROB-001.md
├─ decisions/
│  └─ ADR-007.md
├─ changes/
│  └─ CHG-021.md
├─ reversions/
│  └─ REV-003.md
├─ solutions/
│  └─ SOL-004.md
├─ supersessions/
│  └─ SUP-003.md
└─ generated/
   ├─ active.md
   ├─ by-area.md
   └─ open.md
```

`generated/` 配下は生成物であり、手動編集しない。

## 12. インデックスの扱い

`index.md` や `active.md` を手動で編集しない。

手動管理すると、状態変更時に複数ファイルを修正する必要があり、二重管理になる。

正本は各レコードファイルとする。

生成物：

```text
memory/generated/active.md
memory/generated/by-area.md
memory/generated/by-status.md
memory/generated/by-type.md
```

生成ファイルには以下を明記する。

```markdown
<!-- This file is generated. Do not edit manually. -->
```

## 13. 検索と関連付け

### 13.1 CauseKeyの問題

自由入力のCauseKeyは長期運用で壊れやすい。

例：

```text
auth.jwt.timezone
auth.jwt.timezone-mismatch
auth.jwt.utc
auth.time.utc
jwt-timezone
```

これらは同じ問題を指している可能性があるが、文字列としては別物になる。

したがって、CauseKeyを正本にしない。

### 13.2 検索による候補提示

関連性は、以下を組み合わせて探す。

* ripgrep
* SQLite FTS
* ベクトル検索
* LLMによる候補判定

人間が全件を総当たり確認するのではなく、`MemADR` が候補を数件提示する。

```text
新規BUG作成
↓
既存BUG/SOL/ADRから関連候補を検索
↓
候補を人間またはLLMが確認
↓
必要ならRelatedやResolvedByに記録
```

### 13.3 関係性の種類

関係性は必要になったものだけ使う。

```text
ResolvedBy
Decision
SupersededBy
Reverts
Related
DuplicateOf
SimilarTo
CausedBy
```

双方向リンクは必須にしない。

例：

```markdown
Decision: ADR-007
ResolvedBy: CHG-021
```

逆引きは `memadr` が行う。

## 14. SQLiteの位置付け

SQLiteを使う場合でも、正本にはしない。

推奨構成：

```text
Source of Truth: Markdown
Search / Index: SQLite
```

Markdownは人間、Git、LLMに優しい。

SQLiteは検索、一覧、状態集計、関連候補抽出に向いている。

```text
Markdown
↓
パーサー
↓
SQLite FTS / edges table
↓
検索・一覧・関連候補提示
```

SQLiteは生成物として扱う。

```text
memory/generated/memory.sqlite
```

正本ではないため、壊れてもMarkdownから再生成できる。

## 15. CLI案

ツール名は `MemADR` とする。

コマンド名は `memadr` とする。

想定コマンド：

```bash
memadr init

memadr new bug
memadr new prob
memadr new adr
memadr new chg
memadr new rev
memadr new sol
memadr new sup

memadr list --status unresolved
memadr list --area auth
memadr search "ログイン後に切れる"
memadr related BUG-012

memadr close BUG-012 --resolved-by CHG-021 --verified
memadr split BUG-012 --into adr
memadr check
memadr index
```

ただし、最初からCLIを作り込まない。

まずはGit、Markdown、命名規則、ripgrepで運用し、必要になった箇所だけ薄くツール化する。

## 16. 途中で記録を分ける運用

最初はBUGだけで始める。

```markdown
# BUG-012: 認証状態が画面遷移後に壊れる

Status: INVESTIGATING
Area: auth
Symptom: 画面遷移後にログイン状態が不安定になる
Finding: 原因調査中
```

調査中に「単一バグではなく設計問題」と分かった場合、ADRを作る。

```markdown
# ADR-007: 認証責務をserver session中心に整理する

Status: PROPOSED
Context: BUG-012の原因が単一不具合ではなく、認証責務の分散だった
Decision: 認証状態管理をserver session中心に統一する
Rejected: 既存構造への局所パッチ継続
Consequence: 認証関連の依存方向とディレクトリ構造を整理する
```

実装したらCHGを作る。

```markdown
# CHG-021: 認証関連ディレクトリを再編成

Status: SHIPPED
Reason: ADR-007
Change: 認証責務をserver session中心に集約した
FilesChanged: src/auth/**, src/session/**, src/middleware/**
Verification: 認証フローの主要ケースを確認
```

最後にBUGを閉じる。

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

## 17. 大規模リファクタリング時の扱い

Aというバグを解決するためにディレクトリ構造を大きく変える場合、BUGファイルだけを肥大化させない。

分け方：

```text
BUG: 問題そのもの
ADR: なぜ大きく構造変更する判断をしたか
CHG: 実際に何を変えたか
SUP: 何を無効化したか
```

例：

```text
BUG-012: 認証状態が画面遷移後に壊れる
ADR-007: 認証責務をserver session中心に整理する
CHG-021: 認証関連ディレクトリを再編成
SUP-003: 旧認証実装を廃止
```

## 18. Git運用

Gitで十分なもの：

* 履歴
* 差分
* 分岐
* マージ
* レビュー
* いつ変えたか
* どのコード変更と一緒に記録したか

`MemADR` でやらないもの：

* ブランチ管理
* マージ管理
* コンフリクト解決
* ファイル履歴管理
* コミット履歴管理

コミット単位は、コード変更とmemory更新をまとめる。

例：

```text
fix(auth): reorganize auth responsibility

modified:
  src/auth/...
  src/session/...
  memory/bugs/BUG-012.md
  memory/decisions/ADR-007.md
  memory/changes/CHG-021.md
```

バグファイル1つの更新だけで毎回1コミット作る必要はない。

## 19. LLMへの指示

AGENTS.mdには以下のような方針を書く。

```markdown
## Memory Policy

Use Japanese for memory records by default.

English field names are allowed for structured fields such as Status, Area, Cause, Fix, Verification, FutureRelevance, Context, Decision, Rejected, Consequence, Related, ResolvedBy, and SupersededBy.

Do not write work logs.

The main purpose of memory records is to reduce LLM tokens and context size while preserving reusable development knowledge.

When a bug, behavior change, reversion, architectural decision, reusable solution, or supersession is discovered, create or update a compact memory record.

Each memory record must contain only the essential outcome:
- Status
- Area
- Symptom, Cause, Fix
- Decision, Rejected, Consequence
- Verification
- FutureRelevance
- Related references

Do not include:
- command histories
- step-by-step implementation logs
- long diffs
- speculative reasoning
- full walkthroughs

Each record is the source of truth for its own status.

Do not manually update aggregate files such as active.md, index.md, or by-area.md.

Aggregate views must be generated by tooling such as memadr.

Before creating a new BUG or SOL, search existing records for similar symptoms, causes, and related areas.

If the same problem pattern appears repeatedly, propose creating or updating a SOL record.
```

## 20. 最小導入案

最初からSQLite、ベクトル検索、専用CLIを作り込まない。

第1段階：

```text
Git
Markdown
命名規則
短いテンプレート
ripgrep
```

第2段階：

```text
memadrによる形式チェック
memadrによる一覧生成
memadrによる検索補助
```

第3段階：

```text
SQLite FTS
関連候補提示
LLMによる重複候補判定
```

第4段階：

```text
必要に応じてベクトル検索
```

## 21. 採用する最小ルール

### 21.1 作業ログを書かない

`memory/` 配下には作業ログを書かない。

### 21.2 各レコードは短くする

目安：

```text
BUG: 5〜10行
PROB: 5〜10行
ADR: 5〜10行
CHG: 5〜10行
REV: 5〜10行
SOL: 5〜15行
SUP: 5〜10行
```

### 21.3 状態は各ファイルに書く

active一覧は生成物とする。

### 21.4 関係は最小限にする

双方向リンクを必須にしない。

### 21.5 Gitを正本にする

Markdownを正本とし、SQLiteやindexは再生成可能な補助情報とする。

### 21.6 日本語を基本言語にする

本文、説明、原因、判断理由、修正内容は日本語で書く。

フィールド名のみ英語を許容する。

## 22. 開発言語・提供形態

`MemADR` は、Goによる単体CLIバイナリとして開発する。

利用者側ではビルドやランタイム導入を不要とし、Windows、macOS、Linux向けの実行ファイルを配布する。

開発環境と利用環境は完全に切り分け、利用者は対象OSのバイナリを取得して、リポジトリ内で `memadr` コマンドを実行するだけで利用できる構成とする。

正本はあくまで `memory/` 配下のMarkdownファイルとし、`MemADR` はその上に薄く乗る補助ツールとする。

`MemADR` の責務は、テンプレート生成、形式チェック、一覧生成、検索補助、関連候補提示に限定する。

SQLiteを利用する場合でも正本にはせず、`memory/generated/memory.sqlite` としてMarkdownから再生成可能な補助インデックスとして扱う。

初期段階ではSQLiteやベクトル検索を導入せず、Git、Markdown、命名規則、短いテンプレート、CLIによる形式チェックと一覧生成を中心に実装する。

## 23. 企画の一文要約

`MemADR` は、LLM作業ログを保存するのではなく、LLMに渡すトークン量とコンテキスト量を削減するために、次回以降の開発で再利用できる「問題・判断・変更・解決策・無効化された前提」だけを、Git管理可能な短い日本語Markdownレコードとして蓄積するGo製CLIツールである。

## 24. 暫定結論

現時点では、専用ツールを大きく作るよりも、まず以下で始めるのがよい。

```text
Git + Markdown + 短いテンプレート + 命名規則 + ripgrep
```

そのうえで、実際に運用上の痛みが出た部分だけを `MemADR` として薄くツール化する。

最初から作るべきではないもの：

* Gitの代替
* 完全なIssue管理システム
* 複雑なグラフDB
* 大きなGUI
* 全自動分類システム

最初に作る価値があるもの：

* テンプレート
* AGENTS.mdのルール
* 最小ファイル構成
* 状態とFutureRelevanceの規約
* grepしやすいMarkdown形式
* `memadr check`
* `memadr index`
* `memadr list`
