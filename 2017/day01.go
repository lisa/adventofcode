package main

/* Part A: Compute the sum of input digits with these rules:

1. The input is circular (wraps around) so i[len(i)+1] == i[0]
2. Do work (add) if i[n] == i[n+1] where n is the candidate digit.

Part B:

1. The input is circular (wraps around) so i[len(i)+1] == i[0]
2. Do work (add) iff i[n] == i[len(i)/2]

*/

import (
	"container/ring"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var input = flag.String("input", "1122", "Digits to use")
var partB = flag.Bool("partB", false, "Perform part B solution?")

func main() {
	flag.Parse()

	digits := ring.New(len(*input))
	for _, d := range strings.Split(*input, "") {
		digit, err := strconv.Atoi(d)
		if err != nil {
			fmt.Printf("Couldn't convert %d to a string: %v", d, err)
			return
		}
		digits.Value = digit
		digits = digits.Next()
	}
	// Zip back to the "start"
	digits = digits.Move(-1 * digits.Len())
	sum := 0

	for n := 0; n < digits.Len(); n++ {
		if *partB {
			// Part B logic, compare to position at len/2
			if digits.Value == digits.Move(digits.Len()/2).Value.(int) {
				sum += digits.Value.(int)
			}
		} else {
			// Part A logic
			if digits.Value == digits.Next().Value.(int) {
				sum += digits.Value.(int)
			}
		}
		digits = digits.Next()
	}

	fmt.Printf("Sum: %d\n", sum)
}
