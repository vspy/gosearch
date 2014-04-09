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
  index_dir := flag.String("index-dir", "wiki-index", "output index directory")
  stopwords_file := flag.String("stopwords", "stopwords.txt", "text file with stopwords")
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

  stopwordsIo, serr := os.Open(*stopwords_file)
  if serr != nil {
    log.Fatal(serr)
  }

  scanner := bufio.NewScanner(stopwordsIo)
  stopwords := []string{}
  for scanner.Scan() {
    stopwords = append(stopwords, scanner.Text())
  }

  stopwordsFilter := processing.CreateStopWordsFilter(stopwords)

  analyzer := func(str string) []string {
    return stopwordsFilter(
            processing.LowercaseFilter(
              processing.SimpleTokenizer(str)))
  }

  wikixmlparser.Parse(fi, func(page *wikixmlparser.Page){
    tokens := append( analyzer(page.Title),
                      analyzer(page.Text)...)
  })
}

func usage() {
  fmt.Print("usage: index [-index-dir dir] <source xml file>")
  flag.PrintDefaults()
  os.Exit(1)
}
