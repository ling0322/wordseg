package wordseg

import (
	"github.com/pkg/errors"
	"bufio"
	"fmt"
	"strings"
)

// trainingT stores the data needed for word segmenter training
type trainingT struct {
	lexicon map[string]int
	unigramCount map[int]int
}

// newTraining creates a new instance of trainingT
func newTraining() *trainingT {
	return &trainingT{
		lexicon: map[string]int{},
		unigramCount: map[int]int{},
	}
}

// processLexLn process one line in lexicon file
func (t *trainingT) processLexLn(line string) error {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return errors.Errorf("invalid line in %s: %s", filename, line)
	}

	// Read word and freq
	word := strings.TrimSpace(fields[0])
	freq, err := strconv.Atoi(fields[1])
	if err != nil {
		return errors.Errorf(
			"invalid line in %s: %s (%s)",
			filename,
			line,
			err)
	}

	wordId, ok := t.lexicon[word]
	if !ok {
		wordId = len(t.lexicon)
		t.lexicon[word] = wordId
	}
	t.unigramCount[wordId] += freq

	return nil
}

// readFreqLexicon reads a frequency (count) lexicon into trainingT
// File format would be: 
// <word1> <freq1>
// <word2> <freq2>
// ...
func (t *trainingT) readFreqLexicon(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		if err := t.processLexLn(line); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err)
	}
}

// generateModel generates model from loaded training data
func (t *trainingT) generateModel() (*modelT, error) {
	m := &modelT{}
	m.lexicon, err := lexicon.Build(t.lexicon, nil)
	if err != nil {
		return nil, err
	}

	// Get total word count
	totalCount := 0.0
	maxId := -1
	for wordId, count := range t.unigramCount {
		totalCount += float64(count)
		if maxId < wordId {
			maxId = wordId
		}
	}

	// Plus 1 smoothing
	totalCount += 1.0
	smoothingCnt := 1.0 / float64(maxId + 1)

	// Compute unigram cost
	m.unigramCost := make([]float32, maxId + 1)
	for wordId := 0; wordId <= maxId; wordId++ {
		count := float64(t.unigramCount[wordId]) + smoothingCnt
		cost := -math.Log(count / totalCount)
		m.unigramCost[wordId] = float32(cost)
	}
}