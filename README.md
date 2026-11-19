# research-cli

> Version 1.0.3

A personal research assistant CLI built in Go. Fetch, cite, and save academic papers from arxiv without leaving the terminal.

---

## Installation

```bash
go install github.com/Swkratis210204/research-cli@latest
```

Requires Go 1.27+

---

## Commands

### `research arxiv <id>`

Fetch metadata and abstract for one or more arxiv papers by ID.

```bash
# fetch and print a single paper
research arxiv 1706.03762

# print full abstract
research arxiv 1706.03762 --fullab

# print bibtex citation
research arxiv 1706.03762 --bibtex

# fetch and print multiple papers
research arxiv 1706.03762 2305.12345

# print bibtex for multiple papers
research arxiv 1706.03762 2305.12345 --bibtex

# save single paper as markdown
research arxiv 1706.03762 --save ./papers

# save multiple papers concurrently
research arxiv 1706.03762 2305.12345 --save ./papers

# save multiple papers with bibtex included
research arxiv 1706.03762 2305.12345 --save ./papers --bibtex
```

#### Flags

| Flag | Description |
| --- | --- |
| `--bibtex` | Print or include bibtex citation |
| `--save <path>` | Save paper(s) as markdown to given folder |
| `--fullab` | Print the full abstract in terminal |

#### Saved file format

Each paper is saved as a separate markdown file named:
```
<id>_<slugified_title>.md
e.g. 1706.03762_attention_is_all_you_need.md
```

The file contains title, authors, published date, URL, full abstract, and optionally the bibtex citation.

---

## Notes

- **All commands support multiple IDs** separated by spaces
- **Version suffixes** are stripped automatically — `1706.03762v2` is treated as `1706.03762`
- **Duplicate IDs** in the same command are skipped automatically
- **Existing files** are never overwritten — if a file already exists it is skipped
- **Save path** is created automatically if it does not exist
- **Abstracts** are truncated in the terminal by default — use `--fullab` to print in full or `--save` to get the full abstract in the file
- **Multiple papers** are fetched and saved concurrently

---

## Roadmap (v2)

- `research search <query>` — search papers via Semantic Scholar
- `research explain <concept>` — LLM explanation via Ollama
- `research save <note>` — save a quick research note
- `research compare <id> --mine <path>` — compare a paper to your own work
- Ollama summarization of papers

---

## Built with

- [Go](https://go.dev/)
- [Cobra](https://github.com/spf13/cobra)
- [arxiv API](https://arxiv.org/help/api)

---

## Author

Sokratis-Vissarion Giannoutsos — [github.com/Swkratis210204](https://github.com/Swkratis210204)