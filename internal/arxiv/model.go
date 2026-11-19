package arxiv

// Paper holds the metadata for an arxiv paper
type Paper struct {
	ID        string   // arxiv ID e.g. 2305.12345
	Title     string   // paper title
	Authors   []string // list of authors
	Abstract  string   // paper abstract
	Published string   // publication date
	URL       string   // link to the paper
}
