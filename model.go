package wordseg

import (
	"github.com/ling0322/lexicon"
)

// modelT stores data of a segmenter model 
type modelT struct {
	lexicon *lexicon.Lexicon
	unigramCost []float32
}

func (m *modelT) save() error