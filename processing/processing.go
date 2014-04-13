package processing

import (
  "io"
  "strings"
  "unicode"
)

type Tokenizer func(str string) []string
type TokenFilter func(tokens []string) []string

const (
  normal = 0
  apostrophe = 1
)

func SimpleTokenizer(str string) []string {
  return tokenize(normal, strings.NewReader(str), 0, "", []string{})
}

func isLetter(ch rune) bool {
  return unicode.IsLetter(ch) && !isCJK(ch)
}

func isCJK(ch rune) bool {
// http://en.wikipedia.org/wiki/CJK_Unified_Ideographs
// 4E00-62FF, 6300-77FF, 7800-8CFF, 8D00-9FFF.
  return ch >= 0x4e00 && ch <= 0x9ffff
}

func tokenize(mode int, str *strings.Reader, l int, token string, buffer []string) []string {
  ch, _, err := str.ReadRune()
  if err == io.EOF {
    if l > 0 {
      buffer = append(buffer, token)
    }
    return buffer
  } else {
    switch mode {
      case apostrophe:
        if isLetter(ch) {
          return tokenize(normal, str, l+2, token+"'"+string(ch), buffer)
        } else {
          buffer = append(buffer, token)
          return tokenize(normal, str, 0, "", buffer)
        }
      default:
        if isLetter(ch) {
          return tokenize(normal, str, l+1, token+string(ch), buffer)
        } else if ch == '\'' && l > 0 {
          return tokenize(apostrophe, str, l, token, buffer)
        } else {
          if l > 0 {
            buffer = append(buffer, token)
          }
          return tokenize(normal, str, 0, "", buffer)
        }
    }
  }
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
