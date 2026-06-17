# 実装計画

## 方針

最初から全コマンドを作らず、Markdown正本と薄いCLIという前提を崩さない最小スライスで着手する。初回は `init` と `new` に限定し、今後の `check` `list` `index` で再利用できる内部構造を先に分ける。

## 実装順

1. Go実行環境をワークスペース内に用意する
2. GoモジュールとCLIエントリポイントを作成する
3. `init` の期待動作をテストで固定する
4. `new` の期待動作をテストで固定する
5. `init` と `new` を実装する
6. テンプレート生成、採番、ファイル配置の責務を内部パッケージに分ける
7. テスト結果と残課題を `walkthrough.md` に反映する

## 初期アーキテクチャ

- `cmd/memadr/main.go`: CLI入口
- `internal/app`: コマンド分岐
- `internal/mem`: 種別定義、採番、ディレクトリ解決
- `internal/template`: Markdownテンプレート生成

依存方向は `cmd` → `internal/app` → 下位モジュールとし、CLI層から直接細部を触らない。

## 今回見送る範囲

- `check`
- `list`
- `search`
- `index`
- `related`
- SQLite
- JSON出力
