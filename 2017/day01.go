package main

/* Compute the sum of input digits with these rules:

1. The input is circular (wraps around) so i[len(i)+1] == i[0]
2. Do work (add) if i[n] == i[n+1] where n is the candidate digit.

*/

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var input = flag.String("input", "1122", "Digits to use")

func main() {
	flag.Parse()

	digits := make([]int, len(*input))
	for i, d := range strings.Split(*input, "") {
		digit, err := strconv.Atoi(d)
		if err != nil {
			fmt.Printf("Couldn't convert %d to a string: %v", d, err)
			return
		}
		digits[i] = digit
	}
	sum := 0

	for n := 0; n < len(digits); n++ {
		var lookup int
		if n == len(digits)-1 {
			// need to look at 0
			lookup = 0
		} else {
			lookup = n + 1
		}
		if digits[n] == digits[lookup] {
			sum += digits[n]
		}
	}

	fmt.Printf("Sum: %d\n", sum)
}
