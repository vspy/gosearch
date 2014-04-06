package wikixmlparser

import (
  "io"
  "log"
  "encoding/xml"
)

type Page struct {
  Title string `xml:"title"`
  Text string `xml:"revision>text"`
}
type PageConsumer func(page *Page)

func Parse(reader io.Reader, fn PageConsumer) {
  decoder := xml.NewDecoder(reader)

  for {
    t, err := decoder.Token()
    if err != nil {
      if err == io.EOF {
        break
      } else {
        log.Fatal(err)
      }
    }

    switch se := t.(type) {
      case xml.StartElement:
        if se.Name.Local == "page" {
          var p Page
          decoder.DecodeElement(&p, &se)
          fn(&p)
        }
    }

  }
}
