# 実施内容

## 開始時点

- 企画書と仕様書を読み、初回スコープを `init` と `new` に限定した
- 開発環境には `go` が入っていなかったため、ワークスペース内にローカルGoを展開した
- 開発自体にもADRを使う方針に従い、初回ADRを追加する準備を進めた

## 実装結果

- `go.mod` を作成し、Go CLIプロジェクトとして初期化した
- `cmd/memadr/main.go` から `internal/app` を呼ぶ薄い入口を追加した
- `internal/mem` にディレクトリ初期化、種別定義、採番、ファイル保存を分離した
- `internal/template` に種別ごとのMarkdownテンプレート生成を分離した
- `internal/app/app_test.go` で `init` と `new` の期待動作を先に固定し、その後に実装した
- `Makefile` を追加し、ビルド、テスト、`init`、`new bug`、`new adr` をタスクとして実行できるようにした
- `memory/decisions/ADR-001.md` と `memory/changes/CHG-001.md` を追加し、この開発自体の判断と変更を記録した

## 確認結果

- ローカル展開したGoで `go test ./...` を実行し、全テストが通過した
- `Makefile` 経由でも同等のビルドと実行ができる構成にした

## 今回未着手

- `check`
- `list`
- `search`
- `index`
- `related`
- 設定ファイル対応

## 更新ルール

このファイルには、今回の作業で実際に確定した変更と未着手範囲だけを追記する。試行ログの羅列は書かない。
