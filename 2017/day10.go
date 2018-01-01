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

Part B:

Step 1:
Treat input as a string of bytes instead of numbers. Convert characters to bytes
using their ASCII codes. Append 17,31,73,47,23 to the end of each input
sequence. ex 1,2,3 -> 49,44,50,44,51,17,31,73,47,23 (note: `,` gets converted,
too). The ASCII codes (base 10) are the new numbers to use in place of treating
the provided input as a comma-separated list of numbers.

Step 2:
The part A solution is merely one "round" of the overall completion. For part
B, apply 64 rounds in total, using the same length sequence ("input") in each
round. The current position and skip size should be preserved between rounds.

Step 3:
Once the rounds have been completed the remaining (0..255) is called a sparse
hash. Reduce that list to one, of only 16 numbers (called the dense hash). To do
this, use numeric bitwise XOR to combine each consecutive block of 16 numbers
in the sparse hash (there are 16 such blocks in a list of 256 numbers). So, the
first element in the dense hash is the first sixteen elements of the sparse
hash XOR'd together, the second element in the dense hash is the second sixteen
elements of the sparse hash XOR'd together, etc.

Step 4:
Perform this operation on each of the 16 blocks of 16 numbers in the sparse
hash to determine the 16 numbers in the dense hash.

Step 5:
Represent the dense hash as a hex string.  Convert each number to hex with
leading zero if necessary.

example hashes:
input    = hash
""       = a2582a3a0e66e6e86e3812dcb672a272
AoC 2017 = 33efeb34ea91902bb2f59c9920caa6cd.
1,2,3    = 3efbe78a8d82f29979031a4aa0b16a9d.
1,2,4    = 63960835bcdc130f0b66d7ff4f6a5a8e
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
var debug = flag.Bool("debug", false, "Debug output")
var debugHashComputation = flag.Bool("debugHashComputation", false, "Debug hash computations")
var partB = flag.Bool("partB", false, "Perform part B solution?")

// Append this to the input after it has been converted to ASCII for part B.
var partBSuffix = []int{17, 31, 73, 47, 23}

// For part B:
func ComputeDenseHash(ring *ring.Ring) []byte {
	ret := make([]int, 16)
	for chunk := 0; chunk < 16; chunk++ {
		if *debugHashComputation {
			fmt.Printf("ComputeDenseHash chunk %d\n", chunk)
		}
		// "Seed" the chunk bits with the first value of the 16 digits for ^=
		ret[chunk] = ring.Value.(int)
		if *debugHashComputation {
			fmt.Printf(" * (digit=0) = %d\n", ring.Value.(int))
		}

		ring = ring.Next()
		for digit := 1; digit < 16; digit++ {
			// The digit-th digit in the dense hash
			if *debugHashComputation {
				fmt.Printf(" * (digit=%d) %d ^ %d == %d\n", digit, ret[chunk], ring.Value.(int), ret[chunk]^ring.Value.(int))
			}
			ret[chunk] ^= ring.Value.(int)
			ring = ring.Next()
		} // finsihed 16 digits
		if *debugHashComputation {
			fmt.Printf("Chunk %d: %d=%x\n", chunk, ret[chunk], ret[chunk])
		}
	} // done with the chunks

	// coerce to []byte
	byteRet := make([]byte, len(ret))
	for i := 0; i < len(ret); i++ {
		byteRet[i] = byte(ret[i])
	}
	return byteRet
}

// Print the ring, optionally highlighting the `highlight` value
func PrintRing(r *ring.Ring, highlight int) {
	for i := 0; i < r.Len(); i++ {
		if r.Value == highlight {
			fmt.Printf("(%d)", r.Value)
		} else {
			fmt.Printf("%d", r.Value)
		}
		if i < r.Len()-1 {
			fmt.Printf("->")
		} else {
			fmt.Printf("\n")
		}
		r = r.Next()
	}
}

func PrintRingFrom(r *ring.Ring, startAt int) {
	// Find startAt in the ring and then print based off of it. Make a new Ring
	// to keep caller safe.
	if startAt < 0 {
		PrintRing(r, -1)
		return
	}
	tempRing := ring.New(r.Len())
	for i := 0; i < r.Len(); i++ {
		tempRing.Value = r.Value
		tempRing = tempRing.Next()
		r = r.Next()
	}
	// Both rings will do a full run through their lengths and end up back at
	// "start"
	i := 0
	for tempRing.Value != startAt || i == tempRing.Len() {
		tempRing = tempRing.Next()
		i += 1 // just in case we loop around
	} // Found the start
	PrintRing(tempRing, -1)
	return
}

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
	if *debug {
		fmt.Printf("Set up temporary ring (len=%d): ", newRing.Len())
		PrintRing(newRing, -1)
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

/*
return a slice of input numbers.
For part A this will treat `input` as a comma-separated list.
For part B this will treat each character as something to conver to its ASCII
representation.
*/
func userInputToLengths(input *string) []int {
	ret := make([]int, 0)
	if *partB {
		for _, char := range strings.Split(*input, "") {
			ret = append(ret, int(char[0]))
		}
		for i := 0; i < len(partBSuffix); i++ {
			ret = append(ret, partBSuffix[i])
		}
	} else {
		for _, numberString := range strings.Split(*input, ",") {
			number, err := strconv.Atoi(numberString)
			if err != nil {
				fmt.Printf("Couldn't convert %s to a number: %s\n", numberString, err)
				os.Exit(1)
			}
			ret = append(ret, number)
		}
	}
	return ret
}

/* Performs a single round */
func doRound(inputLengths []int, skipSize, totalSkips *int, ring *ring.Ring) *ring.Ring {
	if *debug {
		fmt.Printf("doRound inputLengths: %v, skipSize=%d, totalSkips=%d\n", inputLengths, *skipSize, *totalSkips)
		fmt.Printf("Input ring ")
		PrintRingFrom(ring, -1)
	}
	for number := 0; number < len(inputLengths); number++ {
		if *debug {
			fmt.Printf(" Reversing %d digits\n", inputLengths[number])
		}
		ring = ReverseRingSlice(ring, inputLengths[number])
		if *debug {
			fmt.Printf(" Reversed portion : ")
			PrintRingFrom(ring, -1)
		}
		// Current position moves forward by length + skipSize
		if *debug {
			fmt.Printf("Skipping main ring by length=%d+skipsize=%d total=%d\n", inputLengths[number], *skipSize, inputLengths[number]+*skipSize)
		}
		ring = ring.Move(inputLengths[number] + *skipSize)
		if *debug {
			fmt.Printf("Input (%d/%d) Skip ring: ", number+1, len(inputLengths))
			PrintRingFrom(ring, -1)
		}
		// Save the total number of skips for later rewinding
		*totalSkips += inputLengths[number] + *skipSize
		// Then increase skipSize
		*skipSize += 1
	} // done with input
	return ring
}

func main() {
	flag.Parse()
	skipSize := 0
	totalSkips := 0
	var rounds int

	ring := ring.New(*listLen)
	for i := 0; i < *listLen; i++ {
		ring.Value = i
		ring = ring.Next()
	}
	// Step 1
	inputLengths := *input
	if *partB {
		rounds = 64
		//		inputLengths += partBSuffix
	} else {
		rounds = 1
	}
	if *debug {
		fmt.Printf("rounds: %d\n", rounds)
	}
	numbers := userInputToLengths(&inputLengths)
	if *debug {
		fmt.Printf("Input lengths: %d, literal=%s\n", numbers, inputLengths)
	}
	//Step 1 complete

	// Step 2 - Perform rounds
	for i := 0; i < rounds; i++ {
		if *debug {
			fmt.Printf("Running round %d. skipSize=%d, totalSkips=%d, numbers=%v\n",
				i+1, skipSize, totalSkips, numbers)
		}
		ring = doRound(numbers, &skipSize, &totalSkips, ring)

		if *debug {
			fmt.Printf("Round %d is over. skipSize=%d, totalSkips=%d\n", i, skipSize, totalSkips)
			fmt.Printf("Ring: ")
			PrintRing(ring, -1)
			fmt.Println()
		}
	}
	//part A only?
	if true || !*partB {
		// Go back to the "beginning," undoing skipping from rounds.
		if *debug {
			fmt.Printf("Moving back %d spots to get back to the 'start' of the ring\n", -1*totalSkips)
		}

		ring = ring.Move(-1 * totalSkips)
	}

	if *debug {
		fmt.Printf("Ring after rounds: ")
		PrintRing(ring, -1)
	}

	// Steps 3-4
	dense := ComputeDenseHash(ring)

	// Step 5

	if *partB {
		fmt.Printf("Dense hash: %02x\n", dense)
	} else {
		fmt.Printf("Product of first two %d * %d = %d\n", ring.Value, ring.Next().Value, (ring.Value.(int))*(ring.Next().Value.(int)))
	}
}
