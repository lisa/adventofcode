package MAIN

import (
	"flag"
	"fmt"
	"os"
)

var (
	partB = flag.Bool("partB", false, "Perform part B solution?")
	inputFile = flat.String("inputFile", "inputs/day!DAY!a.txt", "Input File")
	inputString = flag.String("input", "", "Input string")
	debug = flag.Bool("debug", false, "Debug?")
)

func main() {
	flag.Parse()
}
