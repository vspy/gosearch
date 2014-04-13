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
  "strings"

  "gosearch/wikixmlparser"
  "gosearch/processing"
  "gosearch/invertedindex"
)

func main() {
  cpus := runtime.NumCPU()
  runtime.GOMAXPROCS(cpus)

  indexDir := flag.String("index-dir", "wiki-index", "output index directory")
  stopwordsFile := flag.String("stopwords", "stopwords.txt", "text file with stopwords")
  batchSize := flag.Int("batch", 100, "default batch size")
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

  analyzer := createAnalyzer(*stopwordsFile)

  lemmatizeWorker := func( wg *sync.WaitGroup, in chan []*wikixmlparser.Page, out chan []*invertedindex.IndexDoc ) {
    fmt.Print("<")
    for pages := range in {
      buffer := []*invertedindex.IndexDoc{}
      for _, page := range pages {
        tokens := append( analyzer(page.Title),
                          analyzer(page.Text)...)
        buffer = append(buffer, &invertedindex.IndexDoc{ page.Title, tokens })
      }
      out <- buffer
    }
    wg.Done()
    fmt.Print(">")
  }

  index, ierr := invertedindex.CreateDirIndexWriter(*indexDir)
  if ierr != nil {
    log.Fatal(ierr)
  }

  aggregatorWorker := func( wg *sync.WaitGroup, in chan []*invertedindex.IndexDoc ) {
    aggCnt := 0
    for docs := range in {
      werr := index.Write(docs)
      if werr != nil {
        log.Fatal(werr)
      }
      aggCnt = aggCnt + 1
      fmt.Print(".")
      if aggCnt % 10 == 0 {
        fmt.Printf("%v\n", aggCnt)
        index.PrintStats()
      }
    }
    wg.Done()
  }

  var wgLem sync.WaitGroup
  lemChan := make(chan []*wikixmlparser.Page, cpus)
  aggChan := make(chan []*invertedindex.IndexDoc, cpus)

  for i := 0; i < cpus; i++ {
    wgLem.Add(1)
    go lemmatizeWorker( &wgLem, lemChan, aggChan )
  }

  var wgAgg sync.WaitGroup
  wgAgg.Add(1)
  go aggregatorWorker( &wgAgg, aggChan )

  cnt := 0
  buffer := []*wikixmlparser.Page{}
  wikixmlparser.Parse(fi, func(page *wikixmlparser.Page) bool {
    if page.Redirect.Title == "" &&
        !strings.HasPrefix(page.Title, "Talk:") &&
        !strings.HasPrefix(page.Title, "User:") {
      buffer = append(buffer, page)
    }

    if len(buffer) >= *batchSize {
      lemChan <- buffer
      buffer = []*wikixmlparser.Page{}
    }

    cnt = cnt + 1
    if profile && cnt > 5000 {
      return false
    }

    return true
  })

  if len(buffer) > 0 {
    lemChan <- buffer
  }

  close(lemChan)
  wgLem.Wait()
  close(aggChan)
  wgAgg.Wait()
  index.Close()
}

func createAnalyzer(stopwordsFile string) (func(string) []string){
  stopwordsIo, serr := os.Open(stopwordsFile)
  if serr != nil {
    log.Fatal(serr)
  }

  defer func() {
    if err := stopwordsIo.Close(); err != nil {
      panic(err)
    }
  }()

  scanner := bufio.NewScanner(stopwordsIo)
  stopwords := []string{}
  for scanner.Scan() {
    stopwords = append(stopwords, scanner.Text())
  }

  stopwordsFilter := processing.CreateStopWordsFilter(stopwords)

  return func(str string) []string {
    return stopwordsFilter(
            processing.LowercaseFilter(
              processing.SimpleTokenizer(str)))
  }
}

func usage() {
  fmt.Println("usage: index [-index-dir dir] <source xml file>")
  flag.PrintDefaults()
  os.Exit(1)
}
