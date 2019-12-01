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

func compute(m int) int {
	return (int(math.Floor(float64(m)/3)) - 2)
}

func Compute(m int) int {
	ret := 0
	if !*partB {
		return compute(m)
	} else {
		if *debug {
			fmt.Printf("Computing part B\n")
		}
		newmass := compute(m)
		if *debug {
			fmt.Printf("Mass %d turns into %d\n", m, newmass)
		}
		for {
			if newmass > 0 {
				ret += newmass
				if *debug {
					fmt.Printf("Computing mass for fuel needed for %d\n", newmass)
				}
				newmass = compute(newmass)
			} else {
				if *debug {
					fmt.Printf("Wishing really hard for mass %d\n", newmass)
				}
				return ret
			}
		}
	}
}

func main() {
	flag.Parse()

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
	fmt.Printf("Sum: %d\n", sum)

}
