package main

import (
  "flag"
  "fmt"
  "os"
  "strings"
  "gosearch/invertedindex"
)

func main() {
  limit := flag.Int("limit", 100, "maximum number of results to display, 0 - unlimited")
  indexDir := flag.String("index-dir", "wiki-index", "index directory")
  flag.Parse()

  args := flag.Args()

  if len(args) != 1 {
    usage()
  }

  index, err := invertedindex.CreateDirIndexReader(*indexDir)
  if err != nil {
    fmt.Printf("%v", err)
    os.Exit(1)
  }

  results, serr := index.TermSearch(strings.ToLower(args[0]))
  if serr != nil {
    fmt.Printf("%v", serr)
    os.Exit(1)
  }

  if *limit > 0 {
    results = results[0:*limit]
  }

  for _, r := range results {
    doc, derr := index.GetDoc(r.DocId)
    if derr != nil {
      fmt.Printf("%v", derr)
      os.Exit(1)
    }
    fmt.Printf("%v\t%v\n", r.Score, doc)
  }

  os.Exit(0)
}

func usage() {
  fmt.Println("usage: search [-index-dir dir] <term>")
  flag.PrintDefaults()
  os.Exit(1)
}
