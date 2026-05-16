# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`egc` (emoji-generator-cli) is a single-binary Go CLI that renders short Japanese text into a 500×500 PNG suitable for use as a Slack custom emoji. It is distributed via `go install` / `mise use -g go:github.com/FukeKazki/egc@latest`, so the module path `github.com/FukeKazki/egc` must remain importable as a `main` package.

## Commands

```bash
# Run locally without installing
go run . 完全に理解した
go run . 完全に理解した -c yellow   # color flag: pink|yellow|black|red|green|blue
go run . 完全に理解した -f mplus    # font flag: mono|noto|mplus (default: mono)

# Build a local binary
go build -o egc .

# Install to $GOBIN
go install .

# Format / vet (no test suite exists)
go fmt ./...
go vet ./...
```

Go toolchain is pinned to 1.26.3 via `mise.toml`, while `go.mod` declares `go 1.22` for module compatibility. Don't bump the `go.mod` directive without intent — it controls the minimum version users need when running `go install`.

## Architecture

Everything lives in `main.go`. The flow is linear and worth understanding before editing:

1. **Fonts are embedded at compile time** via separate `//go:embed` directives in `fonts/` — `NotoSansMonoCJKjp-Bold.otf` (default, `mono`), `NotoSansJP-Bold.ttf` (`noto`), `MPLUS1-Black.ttf` (`mplus`). The binary has no runtime font dependency, so every embedded path must stay in the repo. `pickFont` is the sole switch that maps the `-f`/`--font` flag value to the chosen byte slice; adding a font requires a new embed, a new case in `pickFont`, and a README entry.
2. **`fitSize(ttf, text, maxW, maxH)`** does a linear shrink from 500pt downward in 5pt steps, measuring the rendered bounding box at each step until the glyph run fits inside 500×500. This is the part to touch if text overflows or renders too small — not the drawing code below it.
3. **Drawing** uses `golang.org/x/image/font` with `opentype.NewFace` + `font.Drawer`. Glyph centering is computed from `font.BoundString` (the returned `bounds.Min` is typically negative for ascenders, hence the `-minX`/`-minY` offset when computing `originX`/`originY`).
4. **`pickColor`** is the sole source of supported colors. Adding a color requires editing both this switch and the README's color list / flag usage string.
5. **Output filename is `<text>.png` in the current working directory.** No flag overrides this. CJK/space/symbol characters in the argument become part of the filename — keep that in mind when changing output behavior.

## Issue tracking convention

Issues are tracked as YAML files in `.issues/` (e.g. `.issues/1.yaml`) with fields `id`, `title`, `status`, `description`, `references`, `scope`, `created_at`, `updated_at`. This is a local convention, not GitHub Issues. Preserve the schema when adding new ones.
