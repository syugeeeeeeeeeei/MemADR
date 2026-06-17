<!-- This file is generated. Do not edit manually. -->

# By Type

## ADR
- ADR-001 | ACCEPTED | - | 初期実装は標準ライブラリ中心の薄いCLIで開始する
- ADR-002 | ACCEPTED | - | 開発タスクランナーをPowerShell依存Makefileからjustfileへ移行する
- ADR-003 | ACCEPTED | - | レコード処理を共有パーサー基盤に集約する
- ADR-004 | ACCEPTED | - | ヘルプの引数とオプションを定義ベースで描画する
- ADR-005 | ACCEPTED | - | バージョン正本をGit tagとリリース導線へ集約する
- ADR-006 | ACCEPTED | - | GitHub Release生成をタグpush起点のActionsへ委譲する
- ADR-007 | ACCEPTED | - | release再実行時は既存タグをHEADへ再配置する
- ADR-008 | ACCEPTED | - | initでワークフロー案内を同梱する

## CHG
- CHG-001 | SHIPPED | - | `init` と `new` を開始実装した
- CHG-002 | SHIPPED | - | 引数なし実行時のヘルプとヘルプ定義の集約を追加した
- CHG-003 | SHIPPED | - | 開発タスクをjustfileへ移行してマルチプラットフォーム化した
- CHG-004 | SHIPPED | - | 初期版の未実装CLIコマンドを追加した
- CHG-005 | SHIPPED | - | helpにオプション詳細と値候補を追加した
- CHG-006 | SHIPPED | - | versionコマンドとrelease導線を追加した
- CHG-007 | SHIPPED | - | GitHub Actions経由の配布リリース導線を追加した
- CHG-008 | SHIPPED | - | releaseの既存タグ再配置動作を追加した
- CHG-009 | SHIPPED | - | initでワークフロー説明書を生成するようにした
- CHG-010 | SHIPPED | - | READMEを追加した
- CHG-011 | SHIPPED | - | READMEとinit出力へAGENTS.md追記用の指示を追加した
