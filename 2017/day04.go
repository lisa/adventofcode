package main

/*
Part A:

Given a list of candidate passphrases (which consist of a series of words), compute how many are valid given the following rule:

A word may not appear twice in the same passphrase.

aa bb cc dd aa: invalid (aa reappears)
aa bb cc dd aaa: valid (aaa and aa are distinct)

Part B:

Passphrases aren't valid if it contains two words which are anagrams of each other.
NOTE: Part A is a specific case of being an anagram.
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day04b-example.txt", "Input file for Day 4")
var partB = flag.Bool("partB", false, "Perform part B solution?")

// Is needle an anagram of haystack?
func isAnagram(needle, haystack string) bool {
	// if they're a different length it is obviously not an angram
	if len(needle) != len(haystack) {
		return false
	}

	needle_freq := make(map[string]int)
	haystack_freq := make(map[string]int)

	// at this point the needle and haystack are the same length, so this loop is safe
	for i := 0; i < len(needle); i++ {
		needle_freq[needle[i:i+1]] += 1
		haystack_freq[haystack[i:i+1]] += 1
	}

	// go through needle_freq and compare to haystack_freq
	for str, freq := range needle_freq {
		hay_str_freq, ok := haystack_freq[str]
		// if it's not in the haystack, can't be an anagram
		// if the freq of `str` in both aren't identical, not an anagram
		if !ok || hay_str_freq != freq {
			return false
		}
	}

	return true
}

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
		if *partB {
			// take the words and compare against each other if it's anagram
		outer:
			for i := 0; i < len(words); i++ {
				for j := i + 1; j < len(words); j++ {
					valid = !isAnagram(words[i], words[j])
					if !valid {
						break outer
					}
				}
			}
		} else {

			for i := 0; i < len(words); i++ {
				for j := i + 1; j < len(words); j++ {
					if words[i] == words[j] {
						valid = false
						break
					}
				}
			}
		} // end part b check

		if valid {
			validPassphrases += 1
		}

	}
	fmt.Printf("Valid passphrases: %d\n", validPassphrases)
}
