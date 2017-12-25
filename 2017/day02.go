package main

/* Day 2:
Part A:
 Given a space-separated list of space-separated numbers in a multi-line file compute the line-wise sum of checksums (highest-lowest number). ex:

5 1 9 5  = 8
7 5 3    = 4
2 4 6 8  = 6
checksum = 8 + 4 + 6 = 18

Part B:

For each row, find the two numbers that cleanly divide into each other and add up the result for each row.
5	9	2	8 = 8 / 2 = 4
9	4	7	3 = 9 / 3 = 3
3	8	6	5 = 6 / 3 = 2
checksum = 4 + 3 + 2 = 9

Need to store each parsed number in a list so that they can be iterated over. Need to compare each with the other, numerator and denominator, ex:
i = 0
 j = 0
  next j if i == j
  if remainder of d[i] / d[j] == 0 # got it
 end
end
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
var partB = flag.Bool("partB", false, "Perform part B solution?")

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
		var lineNumbers []int

		for _, d := range strings.Split(lineReader.Text(), "\t") {
			// Parse numbers on the line
			number, err := strconv.Atoi(d)
			if err != nil {
				fmt.Printf("Couldn't convert >%s< to a number: %v", d, err)
				os.Exit(1)
			}
			if *partB {
				// part B
				lineNumbers = append(lineNumbers, number)
			} else {
				// part A
				if number < lowest {
					lowest = number
				}
				if number > highest {
					highest = number
				}
			}
		} // done processing numbers for this line, make checksum
		if *partB {
			for i := 0; i < len(lineNumbers); i++ {
			inner:
				for j := 0; j < len(lineNumbers); j++ {
					if j == i {
						continue inner
					}
					if math.Remainder(float64(lineNumbers[i]), float64(lineNumbers[j])) == 0 {
						sum += lineNumbers[i] / lineNumbers[j]
					}
				}
			}
		} else {
			sum += highest - lowest
		}
	}
	fmt.Printf("Table checksum: %d\n", sum)
}
