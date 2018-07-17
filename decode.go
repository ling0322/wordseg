package wordseg

import (
	"strings"

	"github.com/ling0322/lexicon"
)

// initLattice initialize the lattice
func initLattice(costs []float32, lattice []int) {
	for i := range costs {
		costs[i] = 1e38
		lattice[i] = -1
	}

	costs[0] = 0
	lattice[0] = 0
}

// buildLattice builds the lattice for segmenter
func (s *Segmenter) buildLattice(toks []string, costs []float32, lattice []int) {
	for i := 0; i < len(toks); i++ {
		state := lexicon.InitialState()
		for j := i; j < len(toks) && state.Valid(); j++ {
			val, ok := s.lexicon.Traverse(toks[j], &state)
			cntCost := float32(20.0)
			if ok && int(val) < len(s.uniCost) {
				cntCost = s.uniCost[val]
			}

			// Always add single token into lattice
			if ok || j == i {
				// tok[i:j+1] is in lexicon
				cost := costs[i] + cntCost
				if cost < costs[j+1] {
					costs[j+1] = cost
					lattice[j+1] = i
				}
			}
		}
	}
}

// getBestResult gets the best result from lattice
func getBestResult(lattice []int, toks []string) []string {
	words := []string{}
	for cntPos := len(lattice) - 1; cntPos > 0; {
		prevPos := lattice[cntPos]
		if prevPos < 0 {
			// Decoding failed
			return []string{}
		}
		word := strings.Join(toks[prevPos:cntPos], "")
		words = append(words, word)
		cntPos = prevPos
	}

	// Current words is in reversed order
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}

	return words
}

// Seg breaks qyery into words
func (s *Segmenter) Seg(query string) []string {
	toks := tokenize(query)
	costs := make([]float32, len(toks)+1)
	lattice := make([]int, len(toks)+1)

	// Initialize the lattice
	initLattice(costs, lattice)

	// Build lattice
	s.buildLattice(toks, costs, lattice)

	// Return the best result
	return getBestResult(lattice, toks)
}
