package wordseg

import (
	"regexp"
)

// Regexp for tokens
var reTok = regexp.MustCompile(`\w+|.`)

// Tokenize splits string into tokens
func tokenize(query string) []string {
	return reTok.FindAllString(query, -1)
}
