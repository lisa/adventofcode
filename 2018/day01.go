package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	partB     = flag.Bool("partB", false, "Perform part B solution?")
	inputFile = flag.String("inputFile", "inputs/day01a.txt", "Input")
	debug     = flag.Bool("debug", false, "Debug?")
)

func main() {
	flag.Parse()

	fmt.Printf("Day 1\n")

	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't open %s: %v\n", *inputFile, err)
		os.Exit(1)
	}
	calibration := 0
	lineReader := bufio.NewScanner(input)
	freqs := make([]int, 0)
	for lineReader.Scan() {
		line := lineReader.Text()

		number, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Couldn't parse %s: %e\n", line, err)
			os.Exit(1)
		}
		freqs = append(freqs, number)

	}
	if !*partB {
		for _, i := range freqs {
			calibration += i
		}
	} else {
		// seen this freq yet?
		seen := make(map[int]bool)
		done := false
		for !done {
			// go through freqs til the end and we need to wrap, or, just bail out if we have a dupe
			for i := 0; !done && i < len(freqs); i++ {
				calibration += freqs[i]
				if seen[calibration] {
					done = true
				}
				seen[calibration] = true
			}
		}

	}

	fmt.Printf("Final calibration: %d\n", calibration)
}
