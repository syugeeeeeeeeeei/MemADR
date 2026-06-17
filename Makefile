SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command

GO := $(if $(wildcard .tools/go/bin/go.exe),.tools/go/bin/go.exe,go)
BIN := memadr.exe
TITLE ?= タイトル未設定
ARGS ?=

.PHONY: help build test init new-bug new-adr exec-init exec-new-bug exec-new-adr run clean

help:
	@Write-Host "targets:" 
	@Write-Host "  build        - memadr.exe をビルド"
	@Write-Host "  test         - Goテストを実行"
	@Write-Host "  init         - go run で memory/ を初期化"
	@Write-Host "  new-bug      - go run で BUG を新規作成 (TITLE=... 指定可)"
	@Write-Host "  new-adr      - go run で ADR を新規作成 (TITLE=... 指定可)"
	@Write-Host "  exec-init    - ビルド済み memadr.exe で init を実行"
	@Write-Host "  exec-new-bug - ビルド済み memadr.exe で BUG を新規作成 (TITLE=... 指定可)"
	@Write-Host "  exec-new-adr - ビルド済み memadr.exe で ADR を新規作成 (TITLE=... 指定可)"
	@Write-Host "  run          - ビルド済み memadr.exe を任意引数で実行 (ARGS=...)"
	@Write-Host "  clean        - memadr.exe を削除"

build:
	& '$(GO)' build -o '$(BIN)' ./cmd/memadr

test:
	& '$(GO)' test ./...

init:
	& '$(GO)' run ./cmd/memadr init

new-bug:
	& '$(GO)' run ./cmd/memadr new bug '$(TITLE)'

new-adr:
	& '$(GO)' run ./cmd/memadr new adr '$(TITLE)'

exec-init: build
	& '.\$(BIN)' init

exec-new-bug: build
	& '.\$(BIN)' new bug '$(TITLE)'

exec-new-adr: build
	& '.\$(BIN)' new adr '$(TITLE)'

run: build
	& '.\$(BIN)' $(ARGS)

clean:
	if (Test-Path '$(BIN)') { Remove-Item '$(BIN)' }
