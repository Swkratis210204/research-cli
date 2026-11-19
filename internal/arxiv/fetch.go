package arxiv

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// these are the XML structs that match the arxiv API response format
type feed struct {
	Entries []entry `xml:"entry"`
}

type entry struct {
	Title     string   `xml:"title"`
	Abstract  string   `xml:"summary"`
	Published string   `xml:"published"`
	Authors   []author `xml:"author"`
	ID        string   `xml:"id"`
}

type author struct {
	Name string `xml:"name"`
}

// FetchPaper calls the arxiv API and returns a Paper struct
func FetchPaper(id string) (*Paper, error) {
	// build the API URL
	url := fmt.Sprintf("https://export.arxiv.org/api/query?id_list=%s", id)

	// make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		// check if it's a network error
		if strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "dial") {
			return nil, fmt.Errorf("no internet connection or arxiv is unreachable")
		}
		return nil, fmt.Errorf("failed to fetch paper: %w", err)
	}
	defer resp.Body.Close()

	// check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("arxiv returned an error (status %d), try again later", resp.StatusCode)
	}

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// parse the XML response
	var f feed
	if err := xml.Unmarshal(body, &f); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// check if we got a result
	if len(f.Entries) == 0 {
		return nil, fmt.Errorf("no paper found for ID: %s", id)
	}

	// build the Paper struct from the first entry
	e := f.Entries[0]

	// collect author names
	authors := make([]string, len(e.Authors))
	for i, a := range e.Authors {
		authors[i] = a.Name
	}

	return &Paper{
		ID:        id,
		Title:     strings.TrimSpace(e.Title),
		Authors:   authors,
		Abstract:  strings.TrimSpace(e.Abstract),
		Published: e.Published[:10],
		URL:       fmt.Sprintf("https://arxiv.org/abs/%s", id),
	}, nil
}

// stripVersion removes the version suffix from an arxiv ID
// e.g. "1706.03762v2" -> "1706.03762"
func StripVersion(id string) string {
	if idx := strings.Index(id, "v"); idx != -1 {
		return id[:idx]
	}
	return id
}
