package invertedindex

import (
  "os"
  "bufio"
  "fmt"
)

type dirIndexWriter struct {
  location string
  termsDict map[string]uint64
  docId uint64

  docFile *os.File
  docWriter *bufio.Writer
  docWriterPos uint64

  docIndexFile *os.File
  docIndexWriter *bufio.Writer

  indexFile *os.File
  indexWriter *bufio.Writer
  indexPos uint64

  termsDictFile *os.File
  termsDictWriter *bufio.Writer
}

func (writer * dirIndexWriter) PrintStats() {
  fmt.Printf("terms:%v\n", len(writer.termsDict))
}

func CreateDirIndexWriter (location string) (*dirIndexWriter, error) {
  dirErr := os.MkdirAll(location, 0755)
  if dirErr != nil {
    return nil, dirErr
  }

  docFile, docErr := os.Create(docLocation(location))
  if docErr != nil {
    return nil, docErr
  }

  docIndexFile, docIndexErr := os.Create(docIndexLocation(location))
  if docIndexErr != nil {
    return nil, docIndexErr
  }

  indexFile, indexErr := os.Create(indexLocation(location))
  if indexErr != nil {
    return nil, indexErr
  }

  termsDictFile, termsDictErr := os.Create(termsLocation(location))
  if termsDictErr != nil {
    return nil, termsDictErr
  }

  indexWriter := bufio.NewWriter(indexFile)
  indexWriter.Write( []byte{ 0xde, 0xad, 0xbe, 0xef } )

  return &dirIndexWriter{
    location,
    make(map[string]uint64),
    0,
    docFile, bufio.NewWriter(docFile),
    0,
    docIndexFile, bufio.NewWriter(docIndexFile),
    indexFile, indexWriter,
    4,
    termsDictFile, bufio.NewWriter(termsDictFile),
  }, nil
}
