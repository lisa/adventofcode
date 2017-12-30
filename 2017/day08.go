package main

/* Day 8 Part A
For a given set of assembler-like instructions that work with/modify registers
(which start at 0 value) determine the largest value to all of the registers
at the program's completion.

Logic Operands (boolean return):
>
<
>=
==
!=
<=


Modification Operands:
inc
dec

ex:

b inc 5 if a > 1
a inc 1 if b < 5
c dec -10 if a >= 1
c inc -20 if c == 10

At the end, the largest value in any register is 1.

*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "./inputs/day08-example.txt", "Input file")
var debug = flag.Bool("debug", false, "Trace execution?")
var step = flag.Bool("step", false, "Step execution?")

type LineOfCode struct {
	Register       string // 0
	Operation      string // 1
	ModifyBy       int    // 2
	LogicRegister  string // 4
	LogicOperation string // 5
	LogicTest      int    // 6

	RawLine string // the raw line from input
}

/* Returns if work was done and the (new) value, if any, of Register
Give current value of 'Register' and 'LogicRegister' as parameters
*/
func (loc LineOfCode) Execute(valueOfRegister, valueOfLogicRegister int) (bool, int) {
	ret := valueOfRegister
	doWork := false
	switch loc.LogicOperation {
	case ">":
		doWork = (valueOfLogicRegister > loc.LogicTest)
	case "<":
		doWork = (valueOfLogicRegister < loc.LogicTest)
	case ">=":
		doWork = (valueOfLogicRegister >= loc.LogicTest)
	case "<=":
		doWork = (valueOfLogicRegister <= loc.LogicTest)
	case "!=":
		doWork = (valueOfLogicRegister != loc.LogicTest)
	case "==":
		doWork = (valueOfLogicRegister == loc.LogicTest)
	default:
		fmt.Printf("Unknown operation %s. Aborting\n", loc.LogicOperation)
		os.Exit(1)
	}
	if doWork {
		switch loc.Operation {
		case "inc":
			ret = valueOfRegister + loc.ModifyBy
		case "dec":
			ret = valueOfRegister - loc.ModifyBy
		default:
			fmt.Printf("Unexpected operation %s. (Expected 'inc' or 'dec'.) Aborting\n", loc.Operation)
			os.Exit(1)
		}
	}

	return doWork, ret
}

func (loc LineOfCode) String() string {
	return fmt.Sprintf("%s %s %d if %s %s %d", loc.Register, loc.Operation, loc.ModifyBy, loc.LogicRegister, loc.LogicOperation, loc.LogicTest)
}

func PrintRegisters(registers map[string]int) {
	for register, value := range registers {
		fmt.Printf("%s=%v\n", register, value)
	}
}

func Step() {
	reader := bufio.NewReader(os.Stdin)

	reader.ReadString('\n')
	return
}
func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't read file: %s\n", err)
		os.Exit(1)
	}
	defer input.Close()
	lineReader := bufio.NewScanner(input)

	// Registers name => value
	registers := make(map[string]int)

	for lineReader.Scan() {
		// loop over tokens separated by spaces
		line := lineReader.Text()
		var lineOfCode LineOfCode
		lineOfCode.RawLine = line
		for n, token := range strings.Split(line, " ") {
			switch n {
			/*
						   0: register
						   1: operation
						   2: modifyBy
						   3: noop
						   4: logicRegister
						   5: logicOperation
				       6: logicTest
			*/
			case 0:
				lineOfCode.Register = token
			case 1:
				lineOfCode.Operation = token
			case 2:
				modifyBy, err := strconv.Atoi(token)
				if err != nil {
					fmt.Printf("Couldn't convert %s to a number: %s\n", token, err)
					os.Exit(1)
				}
				lineOfCode.ModifyBy = modifyBy
			case 3:
				continue //noop
			case 4:
				lineOfCode.LogicRegister = token
			case 5:
				lineOfCode.LogicOperation = token
			case 6:
				logicTest, err := strconv.Atoi(token)
				if err != nil {
					fmt.Printf("Couldn't convert %s to a number: %s\n", token, err)
					os.Exit(1)
				}
				lineOfCode.LogicTest = logicTest
			}
		} // Processed a line
		if *step {
			fmt.Printf("Registers before running\n")
			PrintRegisters(registers)
			Step()
		}

		if *debug {
			fmt.Printf("Executing %s with %s=%d, %s=%d", lineOfCode, lineOfCode.Register, registers[lineOfCode.Register], lineOfCode.LogicRegister, registers[lineOfCode.LogicRegister])
		}
		didWork, newValue := lineOfCode.Execute(registers[lineOfCode.Register], registers[lineOfCode.LogicRegister])
		if didWork {
			if *debug {
				fmt.Printf(" #=> %s = %d\n", lineOfCode.Register, newValue)
			}
			registers[lineOfCode.Register] = newValue
		} else {
			if *debug {
				fmt.Printf(" #=> No Change\n")
			}
		}
		if *step {
			fmt.Printf("Registers after running\n")
			PrintRegisters(registers)
			Step()
		}

	} // EOF
	highest := 0
	highestRegister := ""
	for register, value := range registers {
		if *debug {
			fmt.Printf("Final register %s=%d\n", register, value)
			fmt.Printf("Comparing %s to %s: %d > %d? => %t\n", register, highestRegister, value, highest, value > highest)
		}
		if value > highest {
			if *debug {
				fmt.Printf("%s usurped %s as highest. Replacing %d with %d\n", register, highestRegister, highest, value)
			}
			highest = value
			highestRegister = register
		}
	}
	fmt.Printf("Highest register is %s, value of %d\n", highestRegister, highest)
}
