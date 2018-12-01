package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	partA     = flag.Bool("partB", false, "Perform part B solution?")
	inputFile = flag.String("inputFile", "inputs/day01a.txt", "Input")
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
	for lineReader.Scan() {
		line := lineReader.Text()

		number, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Couldn't parse %s: %e\n", line, err)
			os.Exit(1)
		}
		calibration += number
	}

	fmt.Printf("Final calibration: %d\n", calibration)
}
