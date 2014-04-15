package invertedindex

import (
  "os"
  "sort"
  "errors"
  "encoding/gob"
  "encoding/binary"
)

var errCorrupt = errors.New("Index corrupt.")

func (reader *dirIndexReader) TermSearch (term string) ([]SearchResult, error) {
  firstOffset := (*reader.termsDict)[term]
  if firstOffset > 0 {
    results, err := fetchByTerm(reader, firstOffset, []SearchResult{})
    if err != nil {
      return nil, err
    }

    sorted := byScore(results)
    sort.Sort(sorted)

    return sorted,nil
  } else {
    return []SearchResult{}, nil
  }
}

type byScore []SearchResult
func (s byScore) Len() int {
  return len(s)
}
func (s byScore) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}
func (s byScore) Less(i, j int) bool {
  return s[i].Score > s[j].Score
}

func fetchByTerm(reader *dirIndexReader, offset uint64, buffer []SearchResult) ([]SearchResult, error) {
  entry, err := readEntry(reader, offset)
  if err != nil {
    return buffer, err
  }

  for _, e := range entry.Elements {
    buffer = append(buffer, SearchResult{e.DocId, e.Freq})
  }

  if entry.TailOffset == NO_TAIL {
    return buffer, nil
  } else {
    return fetchByTerm(reader, entry.TailOffset, buffer)
  }
}

func (reader *dirIndexReader) GetDoc(id uint64) (string, error) {
  seekErr := seekFile(reader.docIndexFile, id*8)
  if seekErr!=nil {
    return "", seekErr
  }

  var docOffset uint64
  err := binary.Read(reader.docIndexFile, binary.LittleEndian, &docOffset)
  if err != nil {
    return "", err
  }

  reader.docFile.Seek(int64(docOffset), 0)
  reader.docReader.Reset(reader.docFile)
  d := gob.NewDecoder(reader.docReader)
  var doc string
  derr := d.Decode(&doc)
  if derr != nil {
    return "", derr
  }

  return doc, nil
}

func readEntry(reader *dirIndexReader, offset uint64) (*indexEntry, error) {
  seekErr := seekIndex(reader, offset)
  if seekErr != nil {
    return nil, seekErr
  }

  n, err := binary.ReadUvarint(reader.indexReader)
  if err != nil {
    return nil, err
  }

  entries := []indexEntryElement{}
  for i:=uint64(0); i<n; i++ {
    docId, derr := binary.ReadUvarint(reader.indexReader)
    if derr != nil {
      return nil, derr
    }
    var freq float64
    ferr := binary.Read(reader.indexReader, binary.LittleEndian, &freq)
    if ferr != nil {
      return nil, ferr
    }
    entries = append(entries, indexEntryElement{docId, freq})
  }

  tail, terr := binary.ReadUvarint(reader.indexReader)
  if terr != nil {
    return nil, terr
  }

  return &indexEntry{entries, tail}, nil
}

func seekIndex(reader *dirIndexReader, offset uint64) error {
  err := seekFile(reader.indexFile, offset)
  if err != nil {
    return err
  }
  reader.indexReader.Reset(reader.indexFile)
  return nil
}

func seekFile(file *os.File, offset uint64) error {
  o := int64(offset)

  newOffset, err := file.Seek(o, 0)
  if err != nil {
    return err
  }
  if newOffset != o {
    return errCorrupt
  }
  return nil
}
