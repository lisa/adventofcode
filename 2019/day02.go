package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	inputText string = `1,0,0,3,1,1,2,3,1,3,4,3,1,5,0,3,2,9,1,19,1,19,5,23,1,9,23,27,2,27,6,31,1,5,31,35,2,9,35,39,2,6,39,43,2,43,13,47,2,13,47,51,1,10,51,55,1,9,55,59,1,6,59,63,2,63,9,67,1,67,6,71,1,71,13,75,1,6,75,79,1,9,79,83,2,9,83,87,1,87,6,91,1,91,13,95,2,6,95,99,1,10,99,103,2,103,9,107,1,6,107,111,1,10,111,115,2,6,115,119,1,5,119,123,1,123,13,127,1,127,5,131,1,6,131,135,2,135,13,139,1,139,2,143,1,143,10,0,99,2,0,14,0`
)

var (
	partB       = flag.Bool("partB", false, "Perform part B solution?")
	inputFile   = flag.String("inputFile", "inputs/day02a.txt", "Input File")
	inputString = flag.String("input", inputText, "Input string")
	debug       = flag.Bool("debug", false, "Debug?")
)

func main() {
	flag.Parse()

	program := make([]int, 0)

	for _, digit := range strings.Split(*inputString, ",") {
		n, err := strconv.Atoi(digit)
		if err != nil {
			fmt.Printf("Couldn't parse %s: %e\n", digit, err)
			os.Exit(1)
		}
		program = append(program, n)
	}

	// Replace
	program[1] = 12
	program[2] = 2
	if *debug {
		fmt.Printf("Full program: %+v\n", program)
	}
	if !*partB {
		// part A
		// go through program
		cursor := 0
		for {
			var opcode, left, right, dest int
			opcode = program[cursor]
			if opcode != 99 {
				left = program[cursor+1]
				right = program[cursor+2]
				dest = program[cursor+3]
			}
			var result int
			switch opcode {
			case 1:

				result = program[left] + program[right]
				if *debug {
					fmt.Printf("[cursor: %d; [%d, %d, %d, %d]; n: %d]: %d + %d = %d.\n",
						cursor, opcode, left, right, dest, cursor+4, program[left], program[right], result)
				}
			case 2:
				result = program[left] * program[right]
				if *debug {
					fmt.Printf("[cursor: %d; [%d, %d, %d, %d]; n: %d]: %d: %d * %d = %d.\n",
						cursor, opcode, left, right, dest, cursor+4, program[left], program[right], result)
				}
			case 99:
				fmt.Printf("End of program. Position 0 value: %d\n", program[0])
				if *debug {
					fmt.Printf("Program at end: %+v\n", program)
				}
				os.Exit(0)
			}
			program[dest] = result
			cursor += 4
		}
	} else {
		// part B
	}

	os.Exit(0)
}
