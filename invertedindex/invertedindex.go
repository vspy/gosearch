package invertedindex

type SearchResult struct {
  DocId uint64
  Score float64
}

type IndexDoc struct {
  Document string
  Tokens []string
}

type IndexWriter interface {
  Write(docs []*IndexDoc) error
  Close()
}

type IndexReader interface {
  TermSearch(term string) ([]SearchResult, error)
  GetDoc(id uint64) (string, error)
  Close()
}
