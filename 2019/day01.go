package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
)

var (
	partB       = flag.Bool("partB", false, "Perform part B solution?")
	inputFile   = flag.String("inputFile", "inputs/day01a.txt", "Input File")
	inputString = flag.String("input", "", "Input string")
	debug       = flag.Bool("debug", false, "Debug?")
)

func Compute(m int) int {
	return (int(math.Floor(float64(m)/3)) - 2)
}

func main() {
	flag.Parse()

	if !*partB {
		// Part A
		input, err := os.Open(*inputFile)
		if err != nil {
			fmt.Printf("Couldn't open %s: %v\n", *inputFile, err)
			os.Exit(1)
		}
	
		sum := 0
		lineReader := bufio.NewScanner(input)
		for lineReader.Scan() {
			line := lineReader.Text()
			number, err := strconv.Atoi(line)
			if err != nil {
				fmt.Printf("Couldn't parse %s: %e\n", line, err)
				os.Exit(1)
			}
	
			sum += Compute(number)
		}
		fmt.Printf("Sum: %d\n",sum)

	} else {
		fmt.Printf("Part B\n")
	}
}
