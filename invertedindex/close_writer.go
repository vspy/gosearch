package invertedindex

import (
  "encoding/gob"
  "os"
)

func (writer *dirIndexWriter) Close() error {
  e := gob.NewEncoder( writer.termsDictWriter )
  err := e.Encode( writer.termsDict )
  if err != nil {
    return err
  }

  dwErr := writer.docWriter.(*os.File).Close()
  if dwErr != nil {
    return dwErr
  }
  diwErr := writer.docIndexWriter.(*os.File).Close()
  if diwErr != nil {
    return diwErr
  }
  iwErr := writer.indexWriter.(*os.File).Close()
  if iwErr != nil {
    return iwErr
  }
  tdErr := writer.termsDictWriter.(*os.File).Close()
  if tdErr != nil {
    return tdErr
  }

  return nil
}
