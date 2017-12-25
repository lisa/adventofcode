package main

/* Day 2:
Part A:
 Given a space-separated list of space-separated digits in a multi-line file compute the line-wise sum of checksums (highest-lowest digit). ex:

5 1 9 5  = 8
7 5 3    = 4
2 4 6 8  = 6
checksum = 8 + 4 + 6 = 18
*/

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day02-example.txt", "Input file for Day 2")

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't open %s for read: %v", inputFile, err)
		os.Exit(1)
	}

	sum := 0

	lineReader := bufio.NewScanner(input)
	for lineReader.Scan() {
		var lowest = math.MaxInt32
		var highest = -1

		for _, d := range strings.Split(lineReader.Text(), "\t") {

			digit, err := strconv.Atoi(d)
			if err != nil {
				fmt.Printf("Couldn't convert >%s< to a digit: %v", d, err)
				os.Exit(1)
			}
			if digit < lowest {
				lowest = digit
			}
			if digit > highest {
				highest = digit
			}
		} // done processing digits for this line, make checksum

		sum += highest - lowest
	}
	fmt.Printf("Table checksum: %d\n", sum)
}
