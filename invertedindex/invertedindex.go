package invertedindex

import (
  "io"
  "path"
)


type IndexWriter interface {
  func Write(doc interface{}, tokens []string)
  func Close()
}

type IndexReader interface {
  func TokenSearch(token string) interface{}
}

type dirIndexWriter struct {
  location string
  termDictionary map[string]uint64
  documentsFile io.Writer
  documentsIndexFile io.Writer
  postingsFile io.Writer
}

type dirIndexReader struct {
  location string
  termDictionary map[string]uint64
  documentsFile io.Reader
  documentsIndexFile io.Reader
  postingsFile io.Reader
}

func CreateDirIndexWriter (location string) (*dirIndexWriter, err error) {
  return &dirIndexWriter{
    location,
    make(map[string]uint64),
    os.Create(path.Join(location, "doc"))
    os.Create(path.Join(location, "docindex"))
    os.Create(path.Join(location, "postings"))
  }, nil
}

func CreateDirIndexReader(location string) (*dirIndexReader, err error) {

  return &dirIndexReader{
    }
}
