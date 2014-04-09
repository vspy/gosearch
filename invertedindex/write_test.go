package invertedindex

import (
  "fmt"
  "testing"
)

func TestFreqTable(t *testing.T) {
  ref := make(map[string]uint64)
  ref["bar"] = 3
  ref["foo"] = 2
  ref["baz"] = 1

  result := freqTable([]string{ "bar", "foo", "bar", "foo", "bar", "baz" })
  if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", ref) {
    t.Errorf("expected %v == %v", result, ref)
  }
}

func TestDocsToEntryElements(t *testing.T) {
  arg := make(map[uint64]([]string))
  arg[0] = []string{"foo","foo","bar"}
  arg[1] = []string{"foo","baz"}
  arg[2] = []string{"bar"}

  ref := make(map[string]([]indexEntryElement))
  ref["foo"] = []indexEntryElement{
      indexEntryElement{0, 2},
      indexEntryElement{1, 1},
    }
  ref["bar"] = []indexEntryElement{
      indexEntryElement{0, 1},
      indexEntryElement{2, 1},
    }
  ref["baz"] = []indexEntryElement{
      indexEntryElement{1, 1},
    }

  result := docsToEntryElements(arg)
  if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", ref) {
    t.Errorf("expected %v == %v", result, ref)
  }
}
