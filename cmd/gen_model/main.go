package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ling0322/lexicon"
	"github.com/ling0322/wordseg"
	"github.com/pkg/errors"
)

const (
	// UniCostFilename is the Filename for unigram cost
	UniCostFilename = "cost.uni"

	// LexiconFilename is the filename for lexicon
	LexiconFilename = "lexicon"

	// ConfigFilename is the filename for config file
	ConfigFilename = "wordseg.conf"
)

// trainingT stores the data needed for word segmenter training
type trainingT struct {
	lexicon  map[string]int32
	uniCount map[int]int
	uniCost  []float32
}

// newTraining creates a new instance of trainingT
func newTraining() *trainingT {
	return &trainingT{
		lexicon:  map[string]int32{},
		uniCount: map[int]int{},
	}
}

// processLexLn process one line in lexicon file
func (t *trainingT) processLexLn(line string) error {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return errors.Errorf("invalid line: %s", line)
	}

	// Read word and freq
	word := strings.TrimSpace(fields[0])
	freq, err := strconv.Atoi(fields[1])
	if err != nil {
		return errors.Errorf("invalid line: %s (%s)", line, err)
	}

	wordID, ok := t.lexicon[word]
	if !ok {
		wordID = int32(len(t.lexicon))
		t.lexicon[word] = wordID
	}
	t.uniCount[int(wordID)] += freq

	return nil
}

// readFreqLexicon reads a frequency (count) lexicon into trainingT
// File format would be:
// <word1> <freq1>
// <word2> <freq2>
// ...
func (t *trainingT) read(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return errors.WithStack(err)
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
		return errors.WithStack(err)
	}

	return nil
}

// calcCost compute the cost for each word
func (t *trainingT) calcCost() {
	// Get total word count
	totalCount := 0.0
	maxID := -1
	for wordID, count := range t.uniCount {
		totalCount += float64(count)
		if wordID > maxID {
			maxID = wordID
		}
	}

	// Compute unigram cost
	t.uniCost = make([]float32, maxID+1)
	for wordID := 0; wordID <= maxID; wordID++ {
		count := float64(t.uniCount[wordID])
		cost := -math.Log(count / totalCount)
		t.uniCost[wordID] = float32(cost)
	}
}

func binaryWrite(fd io.Writer, data interface{}, err error) error {
	if err != nil {
		return err
	}

	err = binary.Write(fd, binary.LittleEndian, data)
	return err
}

// saveUniCost saves the uni_cost data into file
func (t *trainingT) saveUniCost(filename string) error {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	err = binaryWrite(fd, []byte(wordseg.UniCostHeader), err)
	err = binaryWrite(fd, int32(len(t.uniCost)), err)
	err = binaryWrite(fd, t.uniCost, err)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// printUsage prints the usage of this tool
func printUsage() {
	fmt.Println("Usage: gen_model -in INPUT_FILE -out OUTPUT_DIR")
	flag.PrintDefaults()
}

// checkAndFatal checks err, fatals if err != nil
func checkAndFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var inputFile, outputDir string
	flag.StringVar(&inputFile, "in", "", "the input word-count file")
	flag.StringVar(&outputDir, "out", "", "the output model directory")
	flag.Parse()

	// Check parameters
	if inputFile == "" || outputDir == "" {
		printUsage()
		os.Exit(22)
	}

	t := newTraining()
	checkAndFatal(t.read(inputFile))
	t.calcCost()

	// Save uniCost file
	checkAndFatal(t.saveUniCost(filepath.Join(outputDir, UniCostFilename)))

	// Build and save lexicon
	lexicon, err := lexicon.Build(t.lexicon, nil)
	checkAndFatal(err)
	checkAndFatal(lexicon.Save(filepath.Join(outputDir, LexiconFilename)))

	config := wordseg.Config{
		LexiconPath: LexiconFilename,
		UniCostPath: UniCostFilename,
	}
	b, err := json.MarshalIndent(config, "", "  ")
	checkAndFatal(err)
	checkAndFatal(ioutil.WriteFile(filepath.Join(outputDir, ConfigFilename), b, 0664))

	fmt.Println("ok")
}
