package invertedindex

import (
  "encoding/gob"
)

func (writer *dirIndexWriter) Close() error {
  e := gob.NewEncoder( writer.termsDictWriter )
  err := e.Encode( writer.termsDict )
  if err != nil {
    return err
  }

  dwErr := writer.docWriter.Flush()
  if dwErr != nil {
    return dwErr
  }
  dfErr := writer.docFile.Close()
  if dfErr != nil {
    return dfErr
  }

  diwErr := writer.docIndexWriter.Flush()
  if diwErr != nil {
    return diwErr
  }
  difErr := writer.docIndexFile.Close()
  if difErr != nil {
    return difErr
  }

  iwErr := writer.indexWriter.Flush()
  if iwErr != nil {
    return iwErr
  }
  ifErr := writer.indexFile.Close()
  if ifErr != nil {
    return ifErr
  }

  tdErr := writer.termsDictWriter.Flush()
  if tdErr != nil {
    return tdErr
  }
  tdfErr := writer.termsDictFile.Close()
  if tdfErr != nil {
    return tdfErr
  }

  return nil
}
