package processing

import (
  "regexp"
  "strings"
)

type Tokenizer func(str string) []string
type TokenFilter func(tokens []string) []string

var simpleTokenRegex *regexp.Regexp = regexp.MustCompile(`\w('\w+|\w)*`)

func SimpleTokenizer(str string) []string {
  return simpleTokenRegex.FindAllString(str, -1)
}

func LowercaseFilter(tokens []string) []string {
  for idx, token := range tokens {
    tokens[idx] = strings.ToLower(token)
  }

  return tokens
}

func CreateStopWordsFilter(stopWords []string) TokenFilter {
  m := make(map[string]bool)
  for _, value := range stopWords {
    m[value] = true
  }
  return func (tokens []string) []string {
    result := []string{}
    for idx, value := range tokens {
      if !m[value] {
        result = append(result[:], tokens[idx])
      }
    }
    return result
  }
}
