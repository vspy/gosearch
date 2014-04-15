package invertedindex

const NO_TAIL uint64 = 0x0

type indexEntry struct {
  Elements []indexEntryElement
  TailOffset uint64
}

type indexEntryElement struct {
  DocId uint64
  Freq float64
}

