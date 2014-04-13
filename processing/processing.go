package processing

import (
  "io"
  "strings"
  "unicode"
  "log"
)

type Tokenizer func(str string) []string
type TokenFilter func(tokens []string) []string

const (
  start = 0
  normal = 1
  apostrophe = 2
  skip = 3
)

func SimpleTokenizer(str string, maxTokenSize int) []string {
  return tokenize(normal, maxTokenSize, strings.NewReader(str), 0, "", []string{})
}

func isLetter(ch rune) bool {
  return unicode.IsLetter(ch) && !isCJK(ch)
}

func isCJK(ch rune) bool {
// http://en.wikipedia.org/wiki/CJK_Unified_Ideographs
// 4E00-62FF, 6300-77FF, 7800-8CFF, 8D00-9FFF.
  return ch >= 0x4e00 && ch <= 0x9ffff
}

func tokenize(mode int, maxTokenSize int, str *strings.Reader, l int, token string, buffer []string) []string {
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
          return tokenize(normal, maxTokenSize, str, l+2, token+"'"+string(ch), buffer)
        } else {
          buffer = append(buffer, token)
          return tokenize(start, maxTokenSize, str, 0, "", buffer)
        }
      case start:
        if isLetter(ch) {
          return tokenize(normal, maxTokenSize, str, l+1, token+string(ch), buffer)
        } else if isCJK(ch) {
          buffer = append(buffer, string(ch))
          return tokenize(start, maxTokenSize, str, 0, "", buffer)
        } else {
          return tokenize(start, maxTokenSize, str, 0, "", buffer)
        }
      case skip:
        if isLetter(ch) || ch == '\'' {
          return tokenize(skip, maxTokenSize, str, l+1, token, buffer)
        } else {
          return tokenize(start, maxTokenSize, str, 0, "", buffer)
        }
      default:
        if isLetter(ch) {
          if l+1 > maxTokenSize {
            log.Printf("Token is too long: “%v”, skipping it", token)
            return tokenize(skip, maxTokenSize, str, l+1, token, buffer)
          }
          return tokenize(normal, maxTokenSize, str, l+1, token+string(ch), buffer)
        } else if ch == '\'' && l > 0 {
          return tokenize(apostrophe, maxTokenSize, str, l, token, buffer)
        } else {
          if l > 0 {
            buffer = append(buffer, token)
          }
          return tokenize(start, maxTokenSize, str, 0, "", buffer)
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
