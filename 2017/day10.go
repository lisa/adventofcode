package main

/* Day 10 part A
For a given sequence of lenghts (input) apply the following rules to a circular
list of size 256 ints (numbered 0 to 255):

Starting with the first item in the list of numbers reverse the order of the
first n digits where n is the first length (input).

Once the reversal is complete move forward (skip) that many places in the list.

Repeat this process for each input.

Once the input is processed multiply the first two digits in the list and
provide the answer.

*/

import (
	"container/ring"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var input = flag.String("input", "3,4,1,5", "Input for day 10")
var listLen = flag.Int("listLen", 5, "Length of the list")

func ReverseRingSlice(r *ring.Ring, sliceLen int) *ring.Ring {
	if sliceLen <= 1 {
		//nothing to do
		return r
	}
	returnRing := ring.New(r.Len())
	r = r.Move(sliceLen - 1)
	newRing := ring.New(sliceLen)
	for i := 0; i < sliceLen; i++ {
		newRing.Value = r.Value
		newRing = newRing.Next()
		r = r.Prev()
	}
	// build from newRing until i > newRing.Len(), then use r.
	// Make sure r is ready to be read in the right order, +1 to undo Prev() above
	r = r.Move(sliceLen + 1)
	for i := 0; i < returnRing.Len(); i++ {
		if i < newRing.Len() {
			returnRing.Value = newRing.Value
			newRing = newRing.Next()
		} else {
			returnRing.Value = r.Value
			r = r.Next()
		}
		returnRing = returnRing.Next()
	}
	for i := 0; i < returnRing.Len(); i++ {
		returnRing = returnRing.Next()
	}

	return returnRing
}

func main() {
	flag.Parse()
	ring := ring.New(*listLen)
	for i := 0; i < *listLen; i++ {
		ring.Value = i
		ring = ring.Next()
	}
	skipSize := 0
	totalSkips := 0
	for _, numberString := range strings.Split(*input, ",") {
		number, err := strconv.Atoi(numberString)
		if err != nil {
			fmt.Printf("Couldn't convert %s to a number: %s\n", numberString, err)
			os.Exit(1)
		}
		ring = ReverseRingSlice(ring, number)
		ring = ring.Move(number + skipSize)
		totalSkips += number + skipSize
		skipSize += 1
	} // done with input

	// Go back to the "beginning", ie, undo all of the skipping about done in the
	// above loop.
	ring = ring.Move(-1 * totalSkips)
	fmt.Printf("Product of first two %d * %d = %d\n", ring.Value, ring.Next().Value, (ring.Value.(int))*(ring.Next().Value.(int)))
}
