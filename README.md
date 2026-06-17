# MemADR

MemADR は、LLM と人間のために開発知識レコードを短い Markdown 正本として管理する Go 製 CLI です。

作業ログそのものではなく、次回の開発判断に必要な問題、設計判断、変更、解決策、無効化された前提だけを `memory/` 配下へ残します。

## 目的

- LLM に渡すコンテキスト量を減らす
- 開発判断に必要な知識を短く再利用可能な形で残す
- Git 管理しやすい Markdown 正本を維持する
- 一覧、検索、関連候補提示、集約生成を CLI で補助する

## 特徴

- 正本は `memory/` 配下の Markdown
- 種別ごとの短いテンプレートを自動生成
- `check` で形式と参照整合性を検査
- `list` `search` `related` で既存知識を探索
- `index` で `memory/generated/` を再生成
- `version` と `release` で配布版管理を補助

## レコード種別

- `BUG`: バグ、不具合、失敗
- `PROB`: 根本問題、構造的課題
- `ADR`: 設計判断
- `CHG`: 実際に行った変更
- `REV`: 巻き戻し、撤回
- `SOL`: 再利用可能な解決策
- `SUP`: 過去判断や実装の無効化

## クイックスタート

```bash
memadr init
memadr new bug "認証状態が壊れる"
memadr new adr "認証責務を整理する"
memadr check
memadr index
```

`memadr init` を実行すると、標準ディレクトリに加えて `MEMADR_WORKFLOW.md` を生成します。人間向け導線、LLM 向け導線、代表ユースケースはこのファイルを参照してください。

## AGENTS.md に追記する指示

LLMがMemADRを使用するために、以下を `AGENTS.md` に貼り付けてください。

```md
## MemADR運用ポリシー

このリポジトリでは、MemADRを開発判断メモリとして使用する。

作業するエージェントは、必ず `MEMADR_WORKFLOW.md` を読み、その内容に従うこと。

特に、次を守ること。

- `memory/` に作業ログを書かない
- 現在も価値がある情報と、過去情報を区別する
- 古い判断、削除済み機能、無効化済み情報を現在有効な前提として扱わない
- MemADRレコードを追加または更新した場合は、完了前に `memadr check` と `memadr index` を実行する
```

## 基本ワークフロー

1. `memadr init` で `memory/` と `MEMADR_WORKFLOW.md` を作る
2. 問題が出たら `memadr new bug "..."` または `memadr new prob "..."`
3. 設計判断が必要なら `memadr new adr "..."`
4. 実装後に `memadr new chg "..."`
5. `memadr check` で形式と参照を確認
6. `memadr close BUG-001 --resolved-by CHG-001 --verified` で状態更新
7. `memadr index` で集約ファイルを再生成

## 主なコマンド

- `memadr help`
- `memadr init`
- `memadr new <type> [title]`
- `memadr check`
- `memadr list`
- `memadr search <query>`
- `memadr related <record-id>`
- `memadr close <record-id> [--resolved-by CHG-001] [--verified]`
- `memadr index`
- `memadr version [--verbose]`

詳細なオプションは `memadr help <command>` を参照してください。

## ディレクトリ構成

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

`generated/` 配下は生成物です。手動編集せず、`memadr index` で再生成します。

## バージョンとリリース

- `memadr version` は現在のバージョン文字列を表示します
- `memadr version --verbose` は commit、build 日時、作業ツリー状態も表示します
- 正式版バージョンは Git tag を正本にします
- `just release v0.1.0` で release 導線を実行します
- GitHub Actions は `v*` タグ push を契機にマルチプラットフォームバイナリを作成し、GitHub Release へ公開します

配布対象は次の 4 種です。

- `memadr_windows_amd64.exe`
- `memadr_darwin_amd64`
- `memadr_darwin_arm64`
- `memadr_linux_amd64`

## 開発

主要レシピ:

```bash
just build
just test
just release v0.1.0
just dist v0.1.0
```

ローカル変更がある状態では `just release` は止まります。

## 関連ドキュメント

- `MEMADR_WORKFLOW.md`
- `docs/MemADR_企画書.md`
- `docs/MemADR_仕様書.md`
