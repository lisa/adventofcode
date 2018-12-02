package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	partB     = flag.Bool("partB", false, "Perform part B solution?")
	inputFile = flag.String("input", "inputs/day02.txt", "Input")
	debug     = flag.Bool("debug", false, "Debug?")
)

func main() {
	flag.Parse()

	fmt.Printf("Day 2\n")

	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't open %s: %v\n", *inputFile, err)
		os.Exit(1)
	}
	defer input.Close()

	check2 := 0
	check3 := 0

	lineReader := bufio.NewScanner(input)
	inputRows := make([]string, 0)
	for lineReader.Scan() {
		line := lineReader.Text()
		inputRows = append(inputRows, strings.ToLower(line))
	}
	if *debug {
		fmt.Printf("Number of boxes: %d\n", len(inputRows))
	}

	if !*partB {
		for _, line := range inputRows {
			letters := make(map[string]int)
			had2 := false
			had3 := false

			for _, r := range line {
				if *debug {
					fmt.Printf("Line: %s - Letter %s\n", line, string(r))
				}
				letters[string(r)] += 1
			}
			for _, c := range letters {
				if c == 2 && !had2 {
					check2 += 1
					had2 = true
				}
				if c == 3 && !had3 {
					check3 += 1
					had3 = true
				}
			}
			if *debug {
				fmt.Printf(" Line: %s, Letters: %+v\n", line, letters)
				fmt.Printf(" Check2: %d, Check3: %d\n", check2, check3)
			}
		}
		fmt.Printf("Checksum: %d\n", check2*check3)
	} else {

		// Loop through the boxes and compare each to the other in search of a box that
		// has a single id letter difference. To be a candidate every letter in the
		// right side must match the left side...except for one. This means that we've
		// got to remember when we see a difference and continue onwards and bail out if
		// we see /another/ difference; if we didn't bail we have a pair of good boxes!
		var differentLetter int
	boxCompare:
		for left := 0; left < len(inputRows)-1; left++ {
			if *debug {
				fmt.Printf("Comparing %s\n", inputRows[left])
			}
		rightBox:
			for right := left + 1; right < len(inputRows); right++ {
				if *debug {
					fmt.Printf(" Against %s\n", inputRows[right])
				}
				match := false
				for i := 1; i <= len(inputRows[left]); i++ {
					if inputRows[left][i-1:i] != inputRows[right][i-1:i] {
						if match {
							// uh oh, we have a second difference. try the next box
							if *debug {
								fmt.Printf("  Second difference :( %s vs %s\n", inputRows[left][i-1:i], inputRows[right][i-1:i])
							}
							continue rightBox
						} else {
							if *debug {
								fmt.Printf("  We have our first difference: %s vs %s, that means the differentLetter=%d\n", inputRows[left][i-1:i], inputRows[right][i-1:i], i-1)
							}
							differentLetter = i - 1
							match = true
						}
					}
				}
				if *debug {
					fmt.Printf("I claim that %s and %s are the boxes! The only different letter is at %d\n", inputRows[left], inputRows[right], differentLetter)
				}
				// Now strip out the differences
				for i, r := range inputRows[left] {
					if i == differentLetter {
						if *debug {
							fmt.Printf("Skipping letter %s because i=%d and differentLetter=%d\n", string(r), i, differentLetter)
						}
						continue
					} else {
						fmt.Printf("%s", string(r))
					}
				}
				fmt.Printf("\n")
				break boxCompare
			}
		}

	}

}
