package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ling0322/wordseg"
)

// printUsage prints the usage of wordseg
func printUsage() {
	fmt.Println("Usage: wordseg -c MODEL_CONFIG [-i INPUT] [-o OUTPUT]")
	flag.PrintDefaults()
}

// checkAndFatal checks err, fatals if err != nil
func checkAndFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var wsConfig, inputFile, outputFile string
	flag.StringVar(&wsConfig, "c", "", "model config file for wordseg")
	flag.StringVar(&inputFile, "i", "-", "input file, '-' for stdin")
	flag.StringVar(&outputFile, "o", "-", "output file, '-' for stdout")
	flag.Parse()

	// Check parameters
	if wsConfig == "" {
		printUsage()
		os.Exit(22)
	}

	// Create segmenter
	s, err := wordseg.NewSegmenter(wsConfig)
	checkAndFatal(err)

	// Open input and output file
	var r io.Reader
	r = os.Stdin
	if inputFile != "-" {
		fd, err := os.Open(inputFile)
		checkAndFatal(err)
		defer fd.Close()
		r = fd
	}

	var w io.Writer
	w = os.Stdout
	if outputFile != "-" {
		fd, err := os.Create(outputFile)
		checkAndFatal(err)
		defer fd.Close()
		w = fd
	}

	// Break line by line
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		words := s.Seg(scanner.Text())
		w.Write([]byte(strings.Join(words, " ") + "\n"))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
