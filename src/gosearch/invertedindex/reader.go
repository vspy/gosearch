package invertedindex

import (
  "bufio"
  "os"
  "encoding/gob"
)

type dirIndexReader struct {
  location string
  termsDict *map[string]uint64

  docFile *os.File
  docReader *bufio.Reader

  docIndexFile *os.File

  indexFile *os.File
  indexReader *bufio.Reader
}

func CreateDirIndexReader(location string) (*dirIndexReader, error) {

  termsDictFile, tdfErr := os.Open(termsLocation(location))
  if tdfErr != nil {
    return nil, tdfErr
  }

  d := gob.NewDecoder(termsDictFile)
  termsDict := new(map[string]uint64)
  derr := d.Decode(termsDict)
  if derr != nil {
    return nil, derr
  }

  tdfcErr := termsDictFile.Close()
  if tdfcErr != nil {
    return nil, tdfcErr
  }

  indexFile, ifErr := os.Open(indexLocation(location))
  if ifErr != nil {
    return nil, ifErr
  }

  docFile, dfErr := os.Open(docLocation(location))
  if dfErr != nil {
    return nil, dfErr
  }

  docIndexFile, difErr := os.Open(docIndexLocation(location))
  if difErr != nil {
    return nil, difErr
  }

  return &dirIndexReader {
    location,
    termsDict,
    docFile,
    bufio.NewReader(docFile),
    docIndexFile,
    indexFile,
    bufio.NewReader(indexFile),
  }, nil
}

func (reader *dirIndexReader) Close () error {
  dfErr := reader.docFile.Close()
  if dfErr != nil {
    return dfErr
  }
  difErr := reader.docIndexFile.Close()
  if difErr != nil {
    return difErr
  }
  ifErr := reader.indexFile.Close()
  if ifErr != nil {
    return ifErr
  }

  return nil
}
