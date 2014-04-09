package invertedindex

import (
  "io"
)

type IndexDoc struct {
  Document interface{}
  Tokens []string
}

type IndexWriter interface {
  Write(docs []IndexDoc) error
  Close()
}

type IndexReader interface {
  TokenSearch(token string) interface{}
}

type dirIndexReader struct {
  location string
  termDictionary map[string]uint64
  docReader io.Reader
  docIndexReader io.ReadSeeker
  indexReader io.ReadSeeker
}

