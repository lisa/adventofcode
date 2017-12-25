package main

/*
Part A:

Given a list of candidate passphrases (which consist of a series of words), compute how many are valid given the following rule:

A word may not appear twice in the same passphrase.

aa bb cc dd aa: invalid (aa reappears)
aa bb cc dd aaa: valid (aaa and aa are distinct)
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day04-example.txt", "Input file for Day 4")

func main() {
	flag.Parse()
	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't open %s for read: %v", inputFile, err)
		os.Exit(1)
	}
	lineReader := bufio.NewScanner(input)
	validPassphrases := 0
	for lineReader.Scan() {
		var words []string
		for _, word := range strings.Split(lineReader.Text(), " ") {
			words = append(words, word)
		}
		// We have the candidate passphrase in `words`, let's analyze it.
		var valid bool
		valid = true
		for i := 0; i < len(words); i++ {
			for j := i + 1; j < len(words); j++ {
				if words[i] == words[j] {
					valid = false
					break
				}
			}
		}
		if valid {
			validPassphrases += 1
		}
	}
	fmt.Printf("Valid passphrases: %d\n", validPassphrases)
}
