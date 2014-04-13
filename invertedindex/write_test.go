package invertedindex

import (
  "fmt"
  "testing"
  "bytes"
  "bufio"
)

func TestFreqTable(t *testing.T) {
  ref := make(map[string]uint64)
  ref["bar"] = 3
  ref["foo"] = 2
  ref["baz"] = 1

  result := freqTable(&[]string{ "bar", "foo", "bar", "foo", "bar", "baz" })
  if fmt.Sprintf("%v", *result) != fmt.Sprintf("%v", ref) {
    t.Errorf("expected %v == %v", *result, ref)
  }
}

func TestDocsToEntryElements(t *testing.T) {
  arg := make(map[uint64](*[]string))
  arg[0] = &[]string{"foo","foo","bar"}
  arg[1] = &[]string{"foo","baz"}
  arg[2] = &[]string{"bar"}

  ref := make(map[string]([]*indexEntryElement))
  ref["foo"] = []*indexEntryElement{
      &indexEntryElement{0, 2},
      &indexEntryElement{1, 1},
    }
  ref["bar"] = []*indexEntryElement{
      &indexEntryElement{0, 1},
      &indexEntryElement{2, 1},
    }
  ref["baz"] = []*indexEntryElement{
      &indexEntryElement{1, 1},
    }

  result := docsToEntryElements(&arg)

  for idx, r := range *result {
    for idx2, e := range r {
      if fmt.Sprintf("%v", e) != fmt.Sprintf("%v", ref[idx][idx2]) {
        t.Errorf("expected %v == %v", e, ref[idx][idx2])
      }
    }
  }
}

func TestWriteEntry(t *testing.T) {
  writer := testWriter()

  err := writeEntry(writer, "foo", []*indexEntryElement{ &indexEntryElement{0, 1} })

  if err != nil {
    t.Errorf("expected writeEntry to finish succesfully, but got %v", err)
  }

  pos, ok := writer.termsDict["foo"]
  if !ok || pos != 0 {
    t.Error("expected termsDict to be updated")
  }

  if writer.indexPos == 0 {
    t.Error("expected indexPos to be updated")
  }

  if writer.indexWriter.Buffered() == 0 {
    t.Error("expected index file to contain written entries")
  }
}

func TestWriteDoc(t *testing.T) {
  writer := testWriter()

  doc := "just a string"
  idx, err := writeDoc(writer, &doc)

  if err != nil || idx != 0 {
    t.Errorf("expected writeDoc to finish succesfully, but got %v %v", idx, err)
  }

  idx1, err1 :=  writeDoc(writer, &doc)
  if err1 != nil || idx1 != 1 {
    t.Errorf("expected writeDoc to finish succesfully, but got %v %v", idx1, err)
  }

  if writer.docWriterPos == 0 {
    t.Error("expected docWriterPos to be updated")
  }

  if writer.docWriter.Buffered() == 0 {
    t.Error("expected documents file to contain written document")
  }

  if writer.docIndexWriter.Buffered() == 0 {
    t.Error("expected documents index to contain document offset")
  }
}

func testWriter() *dirIndexWriter {
  return &dirIndexWriter{
    "test",
    make(map[string]uint64),
    0,
    nil,
    bufio.NewWriter(new(bytes.Buffer)),
    0,
    nil,
    bufio.NewWriter(new(bytes.Buffer)),
    nil,
    bufio.NewWriter(new(bytes.Buffer)),
    0,
    nil,
    bufio.NewWriter(new(bytes.Buffer)),
  }
}

