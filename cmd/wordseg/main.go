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
	var brkField int
	flag.StringVar(&wsConfig, "c", "", "model config file for wordseg")
	flag.StringVar(&inputFile, "i", "-", "input file, '-' for stdin")
	flag.StringVar(&outputFile, "o", "-", "output file, '-' for stdout")
	flag.IntVar(&brkField, "f", -1, "segment specific field in TSV file, -1 to segment all")
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
	buf := make([]byte, 64*1024*1024)
	scanner.Buffer(buf, len(buf))
	for scanner.Scan() {
		if brkField >= 0 {
			line := scanner.Text()
			fields := strings.Split(line, "\t")
			if len(fields) <= brkField {
				if len(line) > 100 {
					line = line[:100] + "..."
				}
				log.Fatal(fmt.Sprintf("unexpected line: %s", line))
			}
			words := wordseg.RemoveSpace(s.Seg(fields[brkField]))
			fields[brkField] = strings.Join(words, " ")
			w.Write([]byte(strings.Join(fields, "\t") + "\n"))
		} else {
			words := wordseg.RemoveSpace(s.Seg(scanner.Text()))
			w.Write([]byte(strings.Join(words, " ") + "\n"))
		}
	}
	checkAndFatal(scanner.Err())
}
