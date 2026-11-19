package output

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sokratis/research-cli/internal/arxiv"
)

// threshold for abstract truncation in terminal
const abstractThreshold = 500

// truncateAbstract truncates the abstract if it exceeds the threshold
func TruncateAbstract(abstract string) string {
	if len(abstract) <= abstractThreshold {
		return abstract
	}
	return abstract[:abstractThreshold] + "...\n(use --fullab to read the full abstract)"
}

// slugify converts a title to a filename-safe string
// e.g. "Attention Is All You Need" -> "attention_is_all_you_need"
func slugify(title string) string {
	// lowercase
	s := strings.ToLower(title)
	// replace spaces with underscores
	s = strings.ReplaceAll(s, " ", "_")
	// remove any character that is not a letter, number or underscore
	reg := regexp.MustCompile(`[^\w]`)
	s = reg.ReplaceAllString(s, "")
	return s
}

// SavePaper saves a paper as a markdown file to the given path
// if bibtex is provided it is included in the file
func SavePaper(paper *arxiv.Paper, dir string, bibtex string) error {
	// create the directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// build the filename: <id>_<slugified_title>.md
	filename := fmt.Sprintf("%s_%s.md", paper.ID, slugify(paper.Title))
	fullPath := filepath.Join(dir, filename)

	// check if file already exists
	if _, err := os.Stat(fullPath); err == nil {
		fmt.Printf("already exists, skipping: %s\n", fullPath)
		return nil
	}

	// build the markdown content
	content := fmt.Sprintf(`# %s

**Authors:** %s
**Published:** %s
**URL:** %s

## Abstract

%s
`,
		paper.Title,
		strings.Join(paper.Authors, ", "),
		paper.Published,
		paper.URL,
		paper.Abstract,
	)

	// append bibtex if provided
	if bibtex != "" {
		content += fmt.Sprintf("\n## Citation\n\n```bibtex\n%s\n```\n", bibtex)
	}

	// write the file
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to save paper: %w", err)
	}

	fmt.Printf("saved: %s\n", fullPath)
	return nil
}
