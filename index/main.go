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
  cpus := runtime.NumCPU()
  runtime.GOMAXPROCS(cpus)

  indexDir := flag.String("index-dir", "wiki-index", "output index directory")
  stopwordsFile := flag.String("stopwords", "stopwords.txt", "text file with stopwords")
  batchCount := flag.Int("batch", 100, "default batch size")
  batchSize := flag.Int("batch-size", 1024*1024*4, "maximum batch size in bytes")
  cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
  maxTokenSize := flag.Int("max-token-size", 64, "maximum token size in characters")

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

  analyzer := createAnalyzer(*stopwordsFile, *maxTokenSize)

  lemmatizeWorker := func( idx int, wg *sync.WaitGroup, in chan []*wikixmlparser.Page, out chan []*invertedindex.IndexDoc ) {
    log.Printf("started lemmatizer worker %v", idx)
    for pages := range in {
      buffer := []*invertedindex.IndexDoc{}
      for _, page := range pages {
        tokens := append( analyzer(page.Title),
                          analyzer(page.Text)...)
        buffer = append(buffer, &invertedindex.IndexDoc{ page.Title, tokens })
      }
      out <- buffer
    }

    log.Printf("stopping lemmatizer worker %v", idx)
    wg.Done()
  }

  index, ierr := invertedindex.CreateDirIndexWriter(*indexDir)
  if ierr != nil {
    log.Fatal(ierr)
  }

  aggregatorWorker := func( wg *sync.WaitGroup, in chan []*invertedindex.IndexDoc ) {
    log.Print("starting aggregator worker")
    aggMessages :=0
    aggCounter := 0
    for docs := range in {
      werr := index.Write(docs)
      if werr != nil {
        log.Fatal(werr)
      }
      aggMessages = aggMessages + 1
      aggCounter = aggCounter + len(docs)
      if aggMessages % 10 == 0 {
        log.Printf("written to disk %v documents (in %v batches)", aggCounter, aggMessages)
        index.LogStats()
      }
    }
    log.Printf("written to disk %v documents (in %v batches)", aggCounter, aggMessages)
    log.Print("shutting down aggregator worker")
    wg.Done()
  }

  var wgLem sync.WaitGroup
  lemChan := make(chan []*wikixmlparser.Page, cpus)
  aggChan := make(chan []*invertedindex.IndexDoc, cpus)

  for i := 0; i < cpus; i++ {
    wgLem.Add(1)
    go lemmatizeWorker( i, &wgLem, lemChan, aggChan )
  }

  var wgAgg sync.WaitGroup
  wgAgg.Add(1)
  go aggregatorWorker( &wgAgg, aggChan )

  cnt := 0

  buffer := []*wikixmlparser.Page{}
  bufferSize := 0

  wikixmlparser.Parse(fi, func(page *wikixmlparser.Page) bool {
    if page.Redirect.Title == "" {
      buffer = append(buffer, page)
      bufferSize = bufferSize + len(page.Title) + len(page.Text)
    }

    if len(buffer) >= *batchCount || bufferSize >= *batchSize {
      log.Printf("+ %v (%v bytes)", len(buffer), bufferSize)
      lemChan <- buffer
      buffer = []*wikixmlparser.Page{}
      bufferSize = 0
    }

    cnt = cnt + 1
    if profile && cnt > 5000 {
      return false
    }

    return true
  })

  if len(buffer) > 0 {
    log.Printf("+ %v (%v bytes)", len(buffer), bufferSize)
    lemChan <- buffer
  }

  close(lemChan)
  wgLem.Wait()
  close(aggChan)
  wgAgg.Wait()
  index.Close()
}

func createAnalyzer(stopwordsFile string, maxTokenSize int) (func(string) []string){
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
              processing.SimpleTokenizer(str, maxTokenSize)))
  }
}

func usage() {
  fmt.Println("usage: index [-index-dir dir] <source xml file>")
  flag.PrintDefaults()
  os.Exit(1)
}
