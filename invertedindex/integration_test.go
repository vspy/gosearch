package invertedindex

import (
  "os"
  "testing"
  "io/ioutil"
)

func TestCompleteCycle(t *testing.T) {
  tmpdir, err := ioutil.TempDir("", "index")
  if err != nil {
    t.Error(err)
  }
  defer os.Remove(tmpdir)

  index, cErr := CreateDirIndexWriter(tmpdir)
  if cErr != nil {
    t.Error(cErr)
  }

  refDoc := "Poffertjes are a traditional Dutch batter treat."

  index.Write([]*IndexDoc{
    &IndexDoc{
      refDoc,
      []string{"poffertjes","are","a","traditional","dutch","batter","treat"}},
  })

  closeErr := index.Close()
  if closeErr != nil {
    t.Error(closeErr)
  }

  reader, rErr := CreateDirIndexReader(tmpdir)
  if rErr != nil {
    t.Error(closeErr)
  }

  results1, err1 := reader.TermSearch("poffertjes")
  if err1 != nil {
    t.Error(err1)
  }
  if len(results1) != 1 {
    t.Error("Expected exactly one result")
  }
  results2, err2 := reader.TermSearch("pommes")
  if err2 != nil {
    t.Error(err2)
  }
  if len(results2) != 0 {
    t.Error("Expected exactly zero results")
  }

  doc, err3 := reader.GetDoc(0)
  if err3 != nil {
    t.Error(err3)
  }
  if doc != refDoc {
    t.Error("Expected %v = %v", doc, refDoc)
  }

  rcloseErr := reader.Close()
  if rcloseErr != nil {
    t.Error(rcloseErr)
  }
}

