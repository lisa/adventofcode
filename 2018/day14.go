package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
)

var (
	input = flag.String("input", "681901", "Input Data")
	partB = flag.Bool("partB", false, "do part b?")
	debug = flag.Bool("debug", false, "debug")
)

func rubyMod(d, m int) int {
	r := d % m
	if (r < 0 && m > 0) || (r > 0 && m < 0) {
		r += m
	}
	return r
}

// ScanForSequence - for the needle in the haystack. If it exists in the
// haystack, return true.
func ScanForSequence(needle, haystack *[]uint8, offset int) bool {
	if len(*haystack) < len(*needle) {
		return false
	}

	ret := true
	j := 0
	for i := len(*haystack) - len(*needle) - offset; i < len(*haystack) && j < len(*needle); i++ {
		if i < 0 {
			continue
		}
		ret = ret && (*haystack)[i] == (*needle)[j]
		j++
		if !ret {
			break
		}
	}
	return ret

}

// CountDigits - how many digits are in the given integer?
func CountDigits(number int) int {
	if number == 0 {
		return 1
	}
	return int(math.Floor(math.Log10(math.Abs(float64(number)))) + 1)
}

// SplitDigits - Splits +number+ into its digit components.
// The return will be in the same order
// 681901 will be returned [6, 8, 1, 9, 0, 1]
func SplitDigits(number int) []uint8 {
	// we have nothing, so add nothing
	if number == 0 {
		return []uint8{0}
	}
	digitCount := CountDigits(number)
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

// function to split the sum into digits (for iteration/addition)
// need Ruby Mod function again
func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	scores := []uint8{3, 7}

	elf1Index := 0
	elf2Index := 1

	if !*partB {
		inputNumber, err := strconv.Atoi(*input)
		errorIf("couldnt parse input", err)
		// less than, since we've already added 1 ([3,7])
		for len(scores) < inputNumber+10 {
			digits := SplitDigits(int(scores[elf1Index] + scores[elf2Index]))
			// append to scores
			scores = append(scores, digits...)
			// move forward
			newElf1Idx := 1 + int(scores[elf1Index]) + elf1Index
			newElf2Idx := 1 + int(scores[elf2Index]) + elf2Index
			elf1Index = rubyMod(newElf1Idx, len(scores))
			elf2Index = rubyMod(newElf2Idx, len(scores))
		}
		ten := scores[inputNumber : inputNumber+10]
		sum := 0
		sum += int(ten[0]) * 1000000000
		sum += int(ten[1]) * 100000000
		sum += int(ten[2]) * 10000000
		sum += int(ten[3]) * 1000000
		sum += int(ten[4]) * 100000
		sum += int(ten[5]) * 10000
		sum += int(ten[6]) * 1000
		sum += int(ten[7]) * 100
		sum += int(ten[8]) * 10
		sum += int(ten[9]) * 1

		fmt.Printf("Score of the ten recipes: %010d\n", sum)
	} else {
		needle := make([]uint8, len(*input))
		for i, r := range *input {
			n, err := strconv.ParseInt(string(r), 10, 0)
			errorIf("couldnt parse input", err)
			needle[i] = uint8(n)
		}
		for {

			digits := SplitDigits(int(scores[elf1Index] + scores[elf2Index]))
			// append to scores
			scores = append(scores, digits...)
			// move forward
			newElf1Idx := 1 + int(scores[elf1Index]) + elf1Index
			newElf2Idx := 1 + int(scores[elf2Index]) + elf2Index
			elf1Index = rubyMod(newElf1Idx, len(scores))
			elf2Index = rubyMod(newElf2Idx, len(scores))
			// only need to look at the last 5 digits of the score.
			if len(scores) >= len(needle) {
				if ScanForSequence(&needle, &scores, len(digits)-1) {
					fmt.Printf("%s occurs after %d\n",
						*input, len(scores)-len(*input)-len(digits)+1)
					break
				}
			}
		}
	}
}
