package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	partB     = flag.Bool("partB", false, "Perform part B solution?")
	inputFile = flag.String("input", "inputs/day05.txt", "Input")
	debug     = flag.Bool("debug", false, "Debug?")
	debug2    = flag.Bool("debug2", false, "Print new units after match?")

	parta     = regexp.MustCompile(`aA|Aa|bB|Bb|cC|Cc|dD|Dd|eE|Ee|fF|Ff|gG|Gg|hH|Hh|iI|Ii|jJ|Jj|kK|Kk|lL|Ll|mM|Mm|nN|Nn|oO|Oo|pP|Pp|qQ|Qq|rR|Rr|sS|Ss|tT|Tt|uU|Uu|vV|Vv|wW|Ww|xX|Xx|yY|Yy|zZ|Zz`)
	partBList = []string{`a`, `b`, `c`, `d`, `e`, `f`, `g`, `h`, `i`, `j`, `k`, `l`, `m`, `n`, `o`, `p`, `q`, `r`, `s`, `t`, `u`, `v`, `w`, `x`, `y`, `z`}
)

// React - Do the reaction.
func React(inputPolymer *string) {
	for {

		match := parta.FindAllStringIndex(*inputPolymer, 1)
		if len(match) == 0 {
			break
		}
		if *debug2 {
			fmt.Printf("Removed %s at [%d:%d] -> ", (*inputPolymer)[match[0][0]:match[0][1]], match[0][0], match[0][1])
		}
		*inputPolymer = strings.Replace(*inputPolymer, (*inputPolymer)[match[0][0]:match[0][1]], "", 1)
		if *debug2 {
			fmt.Printf("New string is %s\n", *inputPolymer)
		}
	}
	if *debug2 {
		fmt.Printf("Final string is %s\n", *inputPolymer)
	}
}

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
	input.Close()
	if *debug2 {
		fmt.Printf("Starting string: %s\n", polymer)
	}

	if *debug {
		fmt.Printf("Starting units: %d\n", len(polymer))
	}

	React(&polymer)

	if !*partB {
		fmt.Printf("Remaining units: %d\n", len(polymer))
	} else {
		shortestReduction := len(polymer)
		originalPolymer := polymer

		for _, letterToCut := range partBList {
			polymer = originalPolymer
			polymer = strings.Replace(polymer, letterToCut, "", -1)
			polymer = strings.Replace(polymer, strings.ToUpper(letterToCut), "", -1)
			React(&polymer)
			if len(polymer) < shortestReduction {
				shortestReduction = len(polymer)
			}
		}
		fmt.Printf("Shortest possible reaction is %d\n", shortestReduction)
	}
}
