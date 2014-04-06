package main

import (
  "flag"
  "fmt"
  "os"
  "log"

  "gosearch/wikixmlparser"
)

func main() {
//  index_dir := flag.String("index-dir", "wiki-index", "output index directory")
  flag.Parse()
  args := flag.Args()

  if len(args)<1 {
    usage()
  }

  fi, err := os.Open(args[0])
  defer func() {
    if err := fi.Close(); err != nil {
      panic(err)
    }
  }()

  if err != nil {
    log.Fatal(err)
  }

  wikixmlparser.Parse(fi, func(page *wikixmlparser.Page){})
}

func usage() {
  fmt.Print("usage: index [-index-dir dir] <source xml file>")
  flag.PrintDefaults()
  os.Exit(1)
}
