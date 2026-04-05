# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

yqp is a TUI playground for exploring yq (YAML processor), forked from jqp. It lets users write yq expressions interactively against YAML data with live results. Built with Bubbletea (Elm architecture TUI framework), using mikefarah/yq as the query engine.

## Commands

```bash
just build        # Build binary
just test         # Run tests with race detector
just lint         # Run golangci-lint (config in .golangci.yaml)
just check        # Run all checks: verify + lint + test
just audit        # go vet + mod verify + govulncheck
```

Single test: `go test -run TestName -count=1 -race ./path/to/package`

## Architecture

The app follows the Bubbletea Elm architecture (Model → Update → View):

- **`main.go`** → **`cmd/`**: CLI entry via cobra. Handles flags (`-f` file, `-t` theme, `--config`), stdin pipe detection, config loading via viper (`~/.yqp.yaml`).
- **`tui/bubbles/yqplayground/`**: Top-level Bubbletea model (`Bubble`) that composes all sub-bubbles. Split across `model.go` (struct/constructor), `init.go`, `update.go` (message handling), `view.go` (rendering), `commands.go` (async tea.Cmds).
- **`tui/bubbles/`**: Each sub-directory is an independent Bubbletea bubble (component): `queryinput`, `inputdata`, `output`, `statusbar`, `help`, `fileselector`, `state`. Each has its own model, styles, and optionally keys.
- **`tui/theme/`**: Theme system backed by chroma styles. Supports chroma syntax highlighting overrides and custom UI color overrides (`primary`, `secondary`, `error`, `inactive`, `success`).
- **`tui/utils/`**: Terminal size helpers and yq YAML evaluation logic.

## Linting

golangci-lint is configured with strict revive rules (cyclomatic complexity ≤5, cognitive complexity ≤7). Import ordering is enforced via gci: standard → third-party → local module. Uses `depguard` to block `github.com/pkg/errors` and `golang.org/x/net/context`.

## Key Dependencies

- `mikefarah/yq/v4` — YAML query engine (Go library, not CLI)
- `charmbracelet/bubbletea` + `bubbles` + `lipgloss` — TUI framework
- `alecthomas/chroma/v2` — Syntax highlighting
- `atotto/clipboard` — System clipboard (requires xclip/xsel on Linux)
