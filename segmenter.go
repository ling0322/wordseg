package wordseg

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ling0322/lexicon"
	"github.com/pkg/errors"
)

// UniCostHeader is header for uni_cost file
const UniCostHeader = "wordseg.UniCost.v1"

// Config is the struct to marshal config for wordseg
type Config struct {
	LexiconPath string `json:"lexicon_path"`
	UniCostPath string `json:"unicost_path"`
}

// Segmenter provides the API for segment query into words
type Segmenter struct {
	lexicon *lexicon.Lexicon
	uniCost []float32
}

// readConfig reads config from file and fix the path of each files
func readConfig(configFile string) (*Config, error) {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	config := Config{}
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, errors.WithStack(err)
	}

	// Fix path
	configDir := filepath.Dir(configFile)
	if !filepath.IsAbs(config.LexiconPath) {
		config.LexiconPath = filepath.Join(configDir, config.LexiconPath)
	}
	if !filepath.IsAbs(config.UniCostPath) {
		config.UniCostPath = filepath.Join(configDir, config.UniCostPath)
	}

	return &config, nil
}

// readUniCost reads the uniCost file
func readUniCost(uniCostFile string) ([]float32, error) {
	fd, err := os.Open(uniCostFile)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// Read header
	header := make([]byte, len(UniCostHeader))
	if err := binary.Read(fd, binary.LittleEndian, &header); err != nil {
		return nil, errors.WithStack(err)
	}
	if string(header) != UniCostHeader {
		return nil, errors.New("invalid uni-cost file")
	}

	// Read unigram cost
	var numCosts int32
	if err := binary.Read(fd, binary.LittleEndian, &numCosts); err != nil {
		return nil, errors.WithStack(err)
	}
	uniCost := make([]float32, numCosts)
	if err := binary.Read(fd, binary.LittleEndian, &uniCost); err != nil {
		return nil, errors.WithStack(err)
	}

	return uniCost, nil
}

// NewSegmenter creates a new instance of segmenter
func NewSegmenter(configFile string) (*Segmenter, error) {
	// Read config file
	config, err := readConfig(configFile)
	if err != nil {
		return nil, err
	}

	// Read lexicon
	lexicon, err := lexicon.Read(config.LexiconPath)
	if err != nil {
		return nil, err
	}

	// Read unigram-cost
	uniCost, err := readUniCost(config.UniCostPath)
	if err != nil {
		return nil, err
	}

	return &Segmenter{
		uniCost: uniCost,
		lexicon: lexicon,
	}, err
}
