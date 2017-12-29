package main

/*
Day 6, part A
Given a set of memory banks, each having a number of blocks stored in them, go through in turn to rebalance the banks by:

Zero out the largest bank, saving the number of blocks. Starting with the next bank, deposit one block at a time, circling around to banks until all the blocks are gone. Keep track of the final balanced configurations and stop balancing once a repeat configuration is encountered; print how many balancing passes it took to reach.

*/

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "./inputs/day06-example.txt", "Input file")
var partB = flag.Bool("partB", false, "Perform part B solution?")

// This will modify the original slice
func balance(banks []int) ([]int, string) {
	index := 0              // index of the highest element
	highest := banks[index] // we'll assume the first element is the biggest

	for bank := 0; bank < len(banks); bank++ {
		if banks[bank] > highest {
			highest = banks[bank]
			index = bank
		}
	} //end find highest

	banks[index] = 0
	rounds := 0 // we've balanced this many times
	index += 1  // Index into the array, starting with the "next" memory bank after highest; could wrap
	if index >= len(banks) {
		// wrap around if we need to
		index = 0
	}

	// Loop highest times and do manual positioning management of the list
	for rounds < highest {
		banks[index] += 1
		rounds += 1
		index += 1
		if index >= len(banks) {
			// wrap around if we need to
			index = 0
		}
	}
	return banks, arrayToString(banks)
}

func arrayToString(ary []int) string {
	var ret string
	for i := 0; i < len(ary); i++ {
		ret += strconv.Itoa(ary[i])
	}

	return ret
}
func main() {
	flag.Parse()
	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't open %s for read: %v", inputFile, err)
		os.Exit(1)
	}
	defer input.Close()

	lineReader := bufio.NewScanner(input)
	var line string
	banks := 0
	for lineReader.Scan() {
		line = lineReader.Text()
		banks = bytes.Count([]byte(line[:len(line)]), []byte{'\t'})
	}
	banks += 1 // n-1 tabs

	// Create our memory banks
	memoryBanks := make([]int, banks)

	// there is a place in hell for people that use synonyms like this, but,
	// we'll be converting memoryBanks to a string (length `banks`) and
	// using that as the key to store stuff we've already seen.
	observedPatterns := make(map[string]bool)

	// build memory bank
	for i, d := range strings.Split(line, "\t") {
		mag, err := strconv.Atoi(d)
		if err != nil {
			fmt.Printf("Couldn't convert: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Assigning %d to a memory bank\n", mag)
		memoryBanks[i] = mag
	}
	fmt.Printf("Memory banks: %v\n", memoryBanks)

	count := 0
	countSinceFirstObservation := 0
	partBCount := false
	var partBFirstSeen string
balanceLoop:
	for {
		_, strFormatted := balance(memoryBanks)
		count += 1
		if *partB && partBCount {
			countSinceFirstObservation += 1
		}
		if observedPatterns[strFormatted] {
			if *partB {
				if partBFirstSeen == "" {
					// set the pattern to look for
					partBFirstSeen = strFormatted
				} else if partBFirstSeen == strFormatted {
					// We did it! We're done
					break balanceLoop
				}
				partBCount = true
			} else {
				break balanceLoop
			}
		} else {
			observedPatterns[strFormatted] = true
		}
	}
	if *partB {
		fmt.Printf("Encountered the first pattern (%s) after %d loops\n", partBFirstSeen, countSinceFirstObservation)
	} else {
		fmt.Printf("Saw a duplicate pattern after %d iterations\n", count)
	}

}
