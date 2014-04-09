package main

import (
  "flag"
  "fmt"
  "os"
  "bufio"
  "log"

  "gosearch/wikixmlparser"
  "gosearch/processing"
  "gosearch/invertedindex"
)

func main() {
  indexDir := flag.String("index-dir", "wiki-index", "output index directory")
  stopwordsFile := flag.String("stopwords", "stopwords.txt", "text file with stopwords")
  batchSize := flag.Int("batch", 1000, "default batch size")

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

  stopwordsIo, serr := os.Open(*stopwordsFile)
  if serr != nil {
    log.Fatal(serr)
  }

  scanner := bufio.NewScanner(stopwordsIo)
  stopwords := []string{}
  for scanner.Scan() {
    stopwords = append(stopwords, scanner.Text())
  }

  index, ierr := invertedindex.CreateDirIndexWriter(*indexDir)
  if ierr != nil {
    log.Fatal(ierr)
  }

  stopwordsFilter := processing.CreateStopWordsFilter(stopwords)

  analyzer := func(str string) []string {
    return stopwordsFilter(
            processing.LowercaseFilter(
              processing.SimpleTokenizer(str)))
  }

  buffer := []invertedindex.IndexDoc{}

  wikixmlparser.Parse(fi, func(page *wikixmlparser.Page){
    tokens := append( analyzer(page.Title),
                      analyzer(page.Text)...)
    // page.Title as the document body
    buffer := append(buffer, invertedindex.IndexDoc{ page.Title, tokens })
    if len(buffer) >= *batchSize {
      werr := index.Write(buffer)
      if werr != nil {
        log.Fatal(werr)
      }
      buffer = []invertedindex.IndexDoc{}
    }
  })

  index.Close()
}

func usage() {
  fmt.Print("usage: index [-index-dir dir] <source xml file>")
  flag.PrintDefaults()
  os.Exit(1)
}
