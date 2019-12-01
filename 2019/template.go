package MAIN

import (
	"flag"
	"fmt"
	"os"
)

var (
	partB       = flag.Bool("partB", false, "Perform part B solution?")
	inputFile   = flag.String("inputFile", "inputs/day!DAY!a.txt", "Input File")
	inputString = flag.String("input", "", "Input string")
	debug       = flag.Bool("debug", false, "Debug?")
)

func main() {
	flag.Parse()

	if !*partB {
		// part A
		fmt.Printf("Part A")
	} else {
		// part B
	}

	os.Exit(0)
}
