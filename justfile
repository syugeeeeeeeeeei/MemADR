set windows-shell := ["powershell.exe", "-NoProfile", "-Command"]

go_bin := "go"
bin := if os_family() == "windows" { "memadr.exe" } else { "memadr" }
bin_path := if os_family() == "windows" { ".\\" + bin } else { "./" + bin }
clean_cmd := if os_family() == "windows" { "if (Test-Path '" + bin + "') { Remove-Item '" + bin + "' }" } else { "rm -f '" + bin + "'" }

default: _help

# 利用可能なレシピを表示
_help:
    @just --list --unsorted

# memadrバイナリをビルド
build:
    {{ go_bin }} build -o {{ bin }} ./cmd/memadr

# Goテストを実行
test:
    {{ go_bin }} test ./...

# go runでmemoryを初期化
init:
    {{ go_bin }} run ./cmd/memadr init

# go runでBUGを新規作成
new-bug title='タイトル未設定':
    {{ go_bin }} run ./cmd/memadr new bug "{{ title }}"

# go runでADRを新規作成
new-adr title='タイトル未設定':
    {{ go_bin }} run ./cmd/memadr new adr "{{ title }}"

# ビルド済みバイナリでinitを実行
exec-init: build
    {{ bin_path }} init

# ビルド済みバイナリでBUGを新規作成
exec-new-bug title='タイトル未設定': build
    {{ bin_path }} new bug "{{ title }}"

# ビルド済みバイナリでADRを新規作成
exec-new-adr title='タイトル未設定': build
    {{ bin_path }} new adr "{{ title }}"

# ビルド済みバイナリを任意引数で実行
run *args: build
    {{ bin_path }} {{ args }}

# 生成バイナリを削除
clean:
    {{ clean_cmd }}
