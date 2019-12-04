package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	partB       = flag.Bool("partB", false, "Perform part B solution?")
	inputFile   = flag.String("inputFile", "inputs/day04a.txt", "Input File")
	inputString = flag.String("input", "134792-675810", "Input string")
	debug       = flag.Bool("debug", false, "Debug?")
)

// countDigits - how many digits are in the given integer?
func countDigits(number int) int {
	if number == 0 {
		return 1
	}
	return int(math.Floor(math.Log10(math.Abs(float64(number)))) + 1)
}

// splitDigits - Splits +number+ into its digit components.
// The return will be in the same order
// 681901 will be returned [6, 8, 1, 9, 0, 1]
func splitDigits(number int) []uint8 {
	// we have nothing, so add nothing
	if number == 0 {
		return []uint8{0}
	}
	digitCount := countDigits(number)
	ret := make([]uint8, digitCount)

	i := number
	for d := 0; d < digitCount; d++ {
		digit := i % 10
		ret[digitCount-d-1] = uint8(digit)
		i -= digit
		i /= 10
	}
	return ret
}

func validPassword(p int) bool {
	digits := splitDigits(p)
	if len(digits) != 6 {
		return false
	}
	lowestDigit := digits[0]
	lastDigit := digits[0]
	digitFreq := make(map[uint8]int)
	for i, digit := range digits {
		if digit < lastDigit || digit < lowestDigit {
			// we went down
			return false
		}
		if i > 0 && lastDigit == digit {
			digitFreq[digit] += 1
		}
		lastDigit = digit
	}
	return len(digitFreq) != 0
}

func main() {
	flag.Parse()

	bounds := strings.Split(*inputString, "-")

	lower, err := strconv.Atoi(bounds[0])
	if err != nil {
		fmt.Printf("Couldn't convert %s: %s\n", bounds[0], err.Error())
	}
	upper, err := strconv.Atoi(bounds[1])
	if err != nil {
		fmt.Printf("Couldn't convert %s: %s\n", bounds[1], err.Error())
	}
	validPasswords := make([]int, 0)
	for i := lower; i <= upper; i++ {
		if validPassword(i) {
			validPasswords = append(validPasswords, i)
		}
	}
	if !*partB {
		// part A
		fmt.Printf("Out of %d possible, there are %d valid passwords\n", upper-lower+1, len(validPasswords))
	} else {
		// part B
	}

	os.Exit(0)
}
