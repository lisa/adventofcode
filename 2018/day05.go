package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	partB     = flag.Bool("partB", false, "Perform part B solution?")
	inputFile = flag.String("input", "inputs/day05.txt", "Input")
	debug     = flag.Bool("debug", false, "Debug?")
	debug2    = flag.Bool("debug2", false, "Print new units after match?")

	match = regexp.MustCompile(`aA|Aa|bB|Bb|cC|Cc|dD|Dd|eE|Ee|fF|Ff|gG|Gg|hH|Hh|iI|Ii|jJ|Jj|kK|Kk|lL|Ll|mM|Mm|nN|Nn|oO|Oo|pP|Pp|qQ|Qq|rR|Rr|sS|Ss|tT|Tt|uU|Uu|vV|Vv|wW|Ww|xX|Xx|yY|Yy|zZ|Zz`)
)

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't open %s: %v\n", *inputFile, err)
		os.Exit(1)
	}
	defer input.Close()
	lineReader := bufio.NewScanner(input)
	lineReader.Scan()
	polymer := lineReader.Text()

	if *debug {
		fmt.Printf("Starting units: %d\n", len(polymer))
	}

	for {

		match := match.FindAllStringIndex(polymer, 1)
		if len(match) == 0 {
			break
		}
		if *debug {
			fmt.Printf("Matched %s with left index %d and right index %d. Cutting it out.\n", polymer[match[0][0]:match[0][1]], match[0][0], match[0][1])
		}
		polymer = strings.Replace(polymer, polymer[match[0][0]:match[0][1]], "", 1)
		if *debug2 {
			fmt.Printf("New string is %s\n", polymer)
		}

	}

	fmt.Printf("Remaining units: %d\n", len(polymer))

}
