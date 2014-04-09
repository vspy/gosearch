package invertedindex

import (
  "os"
  "io"
)

type dirIndexWriter struct {
  location string
  termsDict map[string]uint64
  docId uint64
  docWriter io.Writer
  docWriterPos uint64
  docIndexWriter io.Writer

  indexWriter io.Writer
  indexPos uint64

  termsDictWriter io.Writer
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

  termsDictWriter, termsDictErr := os.Create(termsLocation(location))
  if termsDictErr != nil {
    return nil, termsDictErr
  }

  return &dirIndexWriter{
    location,
    make(map[string]uint64),
    0,
    docWriter,
    0,
    docIndexWriter,
    indexWriter,
    0,
    termsDictWriter,
  }, nil
}
