package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/sokratis/research-cli/internal/arxiv"
	"github.com/sokratis/research-cli/internal/output"
	"github.com/spf13/cobra"
)

// arxiv ID format: YYMM.NNNNN (e.g. 2305.12345)
var arxivIDPattern = regexp.MustCompile(`^\d{4}\.\d{4,5}$`)

// fetchBibtex fetches the bibtex citation for a paper from arxiv
func fetchBibtex(id string) (string, error) {
	url := fmt.Sprintf("https://arxiv.org/bibtex/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// arxivCmd is the "research arxiv" command
var arxivCmd = &cobra.Command{
	Use:   "arxiv [id]",
	Short: "Fetch a paper from arxiv by ID",
	Long: `Fetch metadata and abstract for one or more arxiv papers by ID.

	Examples:
	research arxiv 1706.03762                          fetch and print paper
	research arxiv 1706.03762 --fullab                 print full abstract
	research arxiv 1706.03762 --bibtex                 print bibtex citation
	research arxiv 1706.03762 --save ./papers          save as markdown
	research arxiv 1706.03762 --save ./papers --bibtex save with bibtex included
	research arxiv 1706.03762 2305.12345 --save ./papers  save multiple papers concurrently

	Notes:
	- version suffixes are stripped automatically (e.g. 1706.03762v2 -> 1706.03762)
	- duplicate IDs in the same command are skipped
	- existing files are never overwritten`,
	Args: cobra.MinimumNArgs(1), // at least one paper ID is required
	Run: func(cmd *cobra.Command, args []string) {
		// strip version suffix and validate each ID
		for i, id := range args {
			args[i] = arxiv.StripVersion(id)
			if !arxivIDPattern.MatchString(args[i]) {
				fmt.Printf("invalid arxiv ID: %s (expected format: YYMM.NNNNN)\n", args[i])
				os.Exit(1)
			}
		}

		// remove duplicate IDs
		seen := make(map[string]bool)
		uniqueIDs := []string{}
		for _, id := range args {
			if seen[id] {
				fmt.Printf("skipping duplicate ID: %s\n", id)
				continue
			}
			seen[id] = true
			uniqueIDs = append(uniqueIDs, id)
		}
		args = uniqueIDs

		// read flags
		bibtex, _ := cmd.Flags().GetBool("bibtex")
		savePath, _ := cmd.Flags().GetString("save")
		fullAbstract, _ := cmd.Flags().GetBool("fullab")

		// --bibtex only: print bibtex for each ID, no save
		if bibtex && savePath == "" {
			for _, id := range args {
				bib, err := fetchBibtex(id)
				if err != nil {
					fmt.Printf("error fetching bibtex for %s: %s\n", id, err)
					continue
				}
				fmt.Println(bib)
			}
			return
		}

		// --save: fetch and save each paper concurrently
		if savePath != "" {
			var wg sync.WaitGroup
			for _, id := range args {
				wg.Add(1)
				go func(paperID string) {
					defer wg.Done()
					p, err := arxiv.FetchPaper(paperID)
					if err != nil {
						fmt.Println("error:", err)
						return
					}
					// fetch bibtex too if --bibtex flag was passed
					bib := ""
					if bibtex {
						bib, err = fetchBibtex(paperID)
						if err != nil {
							fmt.Printf("warning: could not fetch bibtex for %s\n", paperID)
						}
					}
					if err := output.SavePaper(p, savePath, bib); err != nil {
						fmt.Println("error:", err)
					}
				}(id)
			}
			wg.Wait()
			return
		}

		// default: fetch and print to terminal
		for _, id := range args {
			paper, err := arxiv.FetchPaper(id)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}

			fmt.Println("Title:    ", paper.Title)
			fmt.Println("Authors:  ", strings.Join(paper.Authors, ", "))
			fmt.Println("Published:", paper.Published)
			fmt.Println("URL:      ", paper.URL)

			if fullAbstract {
				fmt.Println("\nAbstract:\n", paper.Abstract)
			} else {
				fmt.Println("\nAbstract:\n", output.TruncateAbstract(paper.Abstract))
			}
			fmt.Println("---") // separator between papers
		}
	},
}

// init registers arxivCmd as a subcommand of rootCmd
func init() {
	rootCmd.AddCommand(arxivCmd)
	arxivCmd.Flags().Bool("bibtex", false, "fetch bibtex citation instead of metadata")
	arxivCmd.Flags().String("save", "", "save paper as markdown to given path")
	arxivCmd.Flags().Bool("fullab", false, "print the full abstract in terminal")
}
