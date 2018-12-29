package main

import (
	"flag"
	"fmt"
	"math"
)

var (
	input = flag.Int("input", 681901, "Input Data")
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

// SplitDigits - Splits +number+ into its digit components.
// The return will be in the same order
// 681901 will be returned [6, 8, 1, 9, 0, 1]
func SplitDigits(number int) []int {
	// we have nothing, so add nothing
	if number == 0 {
		return []int{0}
	}
	digitCount := int(math.Floor(math.Log10(float64(number))) + 1)
	if *debug {
		fmt.Printf("Digit count vs number: %d : %d\n", digitCount, number)
	}
	ret := make([]int, digitCount)

	i := number
	for d := 0; d < digitCount; d++ {
		digit := i % 10
		ret[digitCount-d-1] = digit
		i -= digit
		i /= 10
	}
	return ret
}

// function to split the sum into digits (for iteration/addition)
// need Ruby Mod function again

func main() {
	flag.Parse()

	scores := []int{3, 7}

	elf1Index := 0
	elf2Index := 1

	// less than, since we've already added 1 ([3,7])
	for len(scores) < *input+10 {
		if *debug {
			fmt.Printf("len=%d elf 1 index=%d value=%d, elf2 index=%d value=%d\n",
				len(scores), elf1Index, scores[elf1Index], elf2Index, scores[elf2Index])
		}
		digits := SplitDigits(scores[elf1Index] + scores[elf2Index])
		// append to scores
		scores = append(scores, digits...)
		// move forward
		newElf1Idx := 1 + scores[elf1Index] + elf1Index
		newElf2Idx := 1 + scores[elf2Index] + elf2Index
		elf1Index = rubyMod(newElf1Idx, len(scores))
		elf2Index = rubyMod(newElf2Idx, len(scores))
		if *debug {
			fmt.Printf("scores: %d\n", scores)
		}
	}
	if *debug {
		fmt.Printf("All done - Elf 1 index: %d Elf 2 index: %d\n", elf1Index, elf2Index)

		fmt.Printf("Scores: %d\n", scores)

		fmt.Printf("Slice off ten: %d\n", scores[*input:*input+10])
	}
	ten := scores[*input : *input+10]

	sum := ten[0] * 1000000000
	sum += ten[1] * 100000000
	sum += ten[2] * 10000000
	sum += ten[3] * 1000000
	sum += ten[4] * 100000
	sum += ten[5] * 10000
	sum += ten[6] * 1000
	sum += ten[7] * 100
	sum += ten[8] * 10
	sum += ten[9] * 1

	fmt.Printf("Score of the ten recipes: %010d\n", sum)

}
