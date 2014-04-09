package invertedindex

import (
  "encoding/binary"
  "encoding/gob"
  "bytes"
)

func (writer *dirIndexWriter) Write(docs []IndexDoc) error {
  docMap := make(map[uint64]([]string))

  for _, doc := range docs {
    id, err := writeDoc(writer, doc.Document)
    if err != nil {
      return err
    }
    docMap[id] = doc.Tokens
  }

  e := docsToEntryElements(docMap)

  for term, elements := range e {
    werr := writeEntry(writer, term, elements)
    if werr != nil {
      return werr
    }
  }

  return nil
}

func writeDoc(writer *dirIndexWriter, doc interface{}) (uint64, error){
  b := new(bytes.Buffer)
  e := gob.NewEncoder(b)

  err := e.Encode(doc)
  if err != nil {
    return 0, err
  }

  id := writer.docId
  writer.docId = writer.docId + 1

  pos := writer.docWriterPos
  l, werr := b.WriteTo(writer.docWriter)
  if werr != nil {
    return 0, werr
  }

  writer.docWriterPos = pos + uint64(l)

  iwerr := binary.Write( writer.docIndexWriter, binary.LittleEndian, pos )
  if iwerr != nil {
    return 0, iwerr
  }

  return id, nil
}

type indexEntry struct {
  Entries []indexEntryElement
  TailOffest uint64
}

type indexEntryElement struct {
  DocId uint64
  Freq uint64
}

func docsToEntryElements(docs map[uint64]([]string)) map[string]([]indexEntryElement) {
  result := make(map[string]([]indexEntryElement))

  for docId, tokens := range docs {
    t := freqTable(tokens)

    for token, freq := range t {
      result[token] = append(result[token], indexEntryElement{docId, freq})
  } }

  return result
}

func freqTable(tokens []string) map[string]uint64 {
  result := make(map[string]uint64)
  for _, token := range tokens {
    result[token] = result[token] + 1
  }
  return result
}

const NO_TAIL uint64 = 0xffffffffffffffff

func writeEntry(writer *dirIndexWriter, term string, elements []indexEntryElement) error {
  entry := indexEntry { elements, NO_TAIL }

  if val, ok := writer.termsDict[term]; ok {
    entry.TailOffest = val
  }

  writer.termsDict[term] = writer.indexPos
  b := new(bytes.Buffer)
  e := gob.NewEncoder(b)

  err := e.Encode(entry)
  if err != nil {
    return err
  }

  l, werr := b.WriteTo(writer.indexWriter)
  if werr != nil {
    return werr
  }

  writer.indexPos = writer.indexPos + uint64(l)

  return nil
}

