package wikixmlparser

import (
  "testing"
  "strings"
)

func TestParse(t *testing.T) {
  reader := strings.NewReader(`
    <mediawiki xmlns="http://www.mediawiki.org/xml/export-0.8/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.mediawiki.org/xml/export-0.8/ http://www.mediawiki.org/xml/export-0.8.xsd" version="0.8" xml:lang="en">
    <page>
      <title>Page 1 title</title>
      <revision>
        <id>123</id>
        <text>Page 1 text</text>
      </revision>
    </page>
    <page>
      <title>Page 2 title</title>
      <revision>
        <id>234</id>
        <text>Page 2 text</text>
      </revision>      
    </page>
    </mediawiki>
    `)

  expected := []Page{
    Page{Title: "Page 1 title", Text: "Page 1 text"},
    Page{Title: "Page 2 title", Text: "Page 2 text"},
  }

  Parse(reader, func(page *Page) {
    var e Page
    e, expected = expected[0], expected[1:len(expected)]
    if e != *page {
      t.Errorf("expected %v == %v", page, e)
    }
  })
}
