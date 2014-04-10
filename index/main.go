package main

import (
  "flag"
  "fmt"
  "os"
  "bufio"
  "log"
  "runtime"
  "runtime/pprof"
  "sync"

  "gosearch/wikixmlparser"
  "gosearch/processing"
  "gosearch/invertedindex"
)

func main() {
  runtime.GOMAXPROCS(4)

  indexDir := flag.String("index-dir", "wiki-index", "output index directory")
  stopwordsFile := flag.String("stopwords", "stopwords.txt", "text file with stopwords")
  batchSize := flag.Int("batch", 1000, "default batch size")
  cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")

  flag.Parse()

  profile := (*cpuprofile != "")
  if profile {
    f, err := os.Create(*cpuprofile)
    if err != nil {
      log.Fatal(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
  }

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

  lemmatizeWorker := func( wg *sync.WaitGroup, in chan *wikixmlparser.Page, out chan *invertedindex.IndexDoc ) {
    fmt.Print("<")
    for page := range in {
     tokens := append( analyzer(page.Title),
                        analyzer(page.Text)...)
     out <- &invertedindex.IndexDoc{ page.Title, tokens }
    }
    wg.Done()
    fmt.Print(">")
  }

  buffer := []invertedindex.IndexDoc{}

  aggregatorWorker := func( wg *sync.WaitGroup, in chan *invertedindex.IndexDoc ) {
    for doc := range in {
      buffer = append(buffer, *doc)
      // page.Title as the document body
      if len(buffer) >= *batchSize {
        werr := index.Write(buffer)
        if werr != nil {
          log.Fatal(werr)
        }
        buffer = []invertedindex.IndexDoc{}
        fmt.Print(".")
      }
    }

    if len(buffer) > 0 {
      werr := index.Write(buffer)
      if werr != nil {
        log.Fatal(werr)
      }
      fmt.Print(".")
    }
    wg.Done()
  }

  var wgLem sync.WaitGroup
  lemChan := make(chan *wikixmlparser.Page, 4)
  aggChan := make(chan *invertedindex.IndexDoc)

  for i := 0; i < 4; i++ {
    wgLem.Add(1)
    go lemmatizeWorker( &wgLem, lemChan, aggChan )
  }

  var wgAgg sync.WaitGroup
  wgAgg.Add(1)
  go aggregatorWorker( &wgAgg, aggChan )

  cnt := 0
  wikixmlparser.Parse(fi, func(page *wikixmlparser.Page) bool {
    lemChan <- page
    cnt = cnt + 1
    if profile && cnt > 5000 {
      return false
    }
    return true
  })

  close(lemChan)
  wgLem.Wait()
  close(aggChan)
  wgAgg.Wait()
  index.Close()
}

func usage() {
  fmt.Println("usage: index [-index-dir dir] <source xml file>")
  flag.PrintDefaults()
  os.Exit(1)
}
