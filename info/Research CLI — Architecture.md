# Research CLI — Architecture

Nov 19, 2026 · @Sokratis

## High Level Idea

The project is split into two layers that never mix:

- **`cmd/`** — the interface layer. Knows what the user typed, which flags they used, and delegates to the logic layer.
- **`internal/`** — the logic layer. Does the actual work: fetching papers, parsing responses, saving files. Knows nothing about the CLI.

`main.go` is the entry point — it just starts cobra, which routes everything else.

## Folder Structure

| File | Layer | Job |
| --- | --- | --- |
| `main.go` | Entry point | Starts cobra, nothing else |
| `cmd/root.go` | Interface | Sets up the app, `research --help` |
| `cmd/arxiv.go` | Interface | Handles `research arxiv` command and flags |
| `internal/arxiv/model.go` | Logic | Defines what a Paper looks like (title, authors, abstract) |
| `internal/arxiv/fetch.go` | Logic | Calls arxiv API, parses response, fetches bibtex |
| `internal/output/render.go` | Logic | Prints to terminal or saves as markdown file |

## The Flow

What happens when you run `research arxiv 2305.12345 --bibtex`:

1. `main.go` starts, hands control to cobra
2. `cmd/root.go` receives the command, routes to `cmd/arxiv.go`
3. `cmd/arxiv.go` reads the ID (`2305.12345`) and the `--bibtex` flag
4. calls `internal/arxiv/fetch.go` → fetches paper + bibtex from arxiv
5. calls `internal/output/render.go` → prints the result to terminal

`cmd` only asks *what*. `internal` only knows *how*.
