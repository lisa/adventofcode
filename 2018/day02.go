package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	//	"strconv"
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
	for lineReader.Scan() {
		line := lineReader.Text()
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

}
