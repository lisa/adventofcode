package main

/*
Day 5 Part A

For a given list of instructions ("list of integers"), treat each instruction
as a `jmp n` where n is the value of the instruction ("integer"). Once the
instruction is processed, and before processing moves on to the next
instruction, it is stored in its same position with its value incremented one.

For example, given this instruction set set: [0, 3, 0, 1, -3], the steps go:

Positive jumps ("forward") move downward; negative jumps move upward. For
egibility in this example, these offset values will be written all on one
line, with the current instruction marked in parentheses. The following steps
would be taken before an exit is found:

(0) 3  0  1  -3  - before we have taken any steps.
(1) 3  0  1  -3  - jump with offset 0 (that is, don't jump at all). Fortunately, the instruction is then incremented to 1.
 2 (3) 0  1  -3  - step forward because of the instruction we just modified. The first instruction is incremented again, now to 2.
 2  4  0  1 (-3) - jump all the way to the end; leave a 4 behind.
 2 (4) 0  1  -2  - go back to where we just were; increment -3 to -2.
 2  5  0  1  -2  - jump 4 steps forward, escaping the maze.

*/
import (
	"bufio"
	"bytes"
	"container/ring"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

var inputFile = flag.String("inputFile", "./inputs/day05-example.txt", "Instructions Input File")

type Instruction struct {
	Step int
}

// https://stackoverflow.com/questions/24562942/golang-how-do-i-determine-the-number-of-lines-in-a-file-efficiently
// Quicker than iterating with bufio
func countLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't read file: %s\n", err)
		os.Exit(1)
	}
	defer input.Close()
	lines, err := countLines(input)
	if err != nil {
		fmt.Printf("Couldn't count the lines in %s: %s\n", inputFile, err)
		os.Exit(1)
	}
	input.Seek(0, 0)

	lineReader := bufio.NewScanner(input)
	fmt.Printf("Lines: %d\n", lines)

	instructions := ring.New(lines)

	// Build instruction set
	for lineReader.Scan() {
		instruction, err := strconv.Atoi(lineReader.Text())

		if err != nil {
			fmt.Printf("Couldn't convert to int: %s\n", err)
			os.Exit(1)
		}
		instructions.Value = Instruction{Step: instruction}
		instructions = instructions.Next()
	}

	outOfBounds := false
	// We'll be zero-based in our reckoning even though the Ring is one-based.
	instructionPosition := 0
	jumps := 0
	stepAmount := 0
	instruction := instructions.Value

	for !outOfBounds {
		instruction = instructions.Value.(Instruction)
		stepAmount = instruction.(Instruction).Step

		if instructionPosition+stepAmount > instructions.Len()+1 {
			outOfBounds = true
		} else {
			jumps += 1
			instructionPosition += stepAmount
			instructions.Value = Instruction{Step: stepAmount + 1}
			instructions = instructions.Move(stepAmount)
		}
	}

	fmt.Printf("Jumps needed: %d\n", jumps)
}
