package invertedindex

import (
  "os"
  "io"
)

type dirIndexWriter struct {
  location string
  termDictionary map[string]uint64
  docId uint64
  docWriter io.Writer
  docWriterPos uint64
  docIndexWriter io.Writer
  indexWriter io.Writer
  termDictWriter io.Writer
}

func CreateDirIndexWriter (location string) (*dirIndexWriter, error) {
  docWriter, docErr := os.Create(docLocation(location))
  if docErr != nil {
    return nil, docErr
  }

  docIndexWriter, docIndexErr := os.Create(docIndexLocation(location))
  if docIndexErr != nil {
    return nil, docIndexErr
  }

  indexWriter, indexErr := os.Create(indexLocation(location))
  if indexErr != nil {
    return nil, indexErr
  }

  termDictWriter, termDictErr := os.Create(termsLocation(location))
  if termDictErr != nil {
    return nil, termDictErr
  }

  return &dirIndexWriter{
    location,
    make(map[string]uint64),
    0,
    docWriter,
    0,
    docIndexWriter,
    indexWriter,
    termDictWriter,
  }, nil
}
