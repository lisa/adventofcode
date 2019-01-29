package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	. "github.com/logrusorgru/aurora"
)

var (
	inputFile        = flag.String("input", "inputs/day16-detection.txt", "input file")
	partB            = flag.Bool("partB", false, "do part b solution?")
	debug            = flag.Bool("debug", false, "debug?")
	debug2           = flag.Bool("debug2", false, "pause on part b execution?")
	detectionMatcher = regexp.MustCompile(`(?mU)Before:\s+\[(\d+),\s+(\d+),\s+(\d+),\s+(\d+)\]\n(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\nAfter:\s+\[(\d+),\s+(\d+),\s+(\d+),\s+(\d+)\]`)
)

func printAndColourize(opid OpcodeId, t bool) string {
	var ret string
	if t {
		ret = fmt.Sprintf("%s", Green(OpCodeTypeToString(opid)).Bold())
	} else {
		ret = fmt.Sprintf("%s", Red(OpCodeTypeToString(opid)).Bold())
	}
	return ret
}

// OpcodeId - numerical ID code for opcodes
type OpcodeId uint8

const (
	gtir OpcodeId = iota // (OpID 0; i=89) gt if val A > reg B { 1 -> reg C } else { 0 -> reg C }
	mulr                 // (OpID 1; i=21) multiply reg A * reg B -> reg C
	seti                 // (OpID 2; i=729) set reg C to val A (copy value)
	gtrr                 // (OpID 3; i=26) gt if reg A > reg B { 1 -> reg C } else { 0 -> reg C }
	bori                 // (OpID 4; i=667) bitwise or reg A | val B -> reg C
	borr                 // (OpID 5; i=664) bitwise or reg A | reg B -> reg C
	banr                 // (OpID 6; i=127) bitwise and reg A & reg b -> reg C
	eqri                 // (OpID 7; i=261) eq if reg A == val B { 1 -> reg C } else { 0 -> reg C }
	bani                 // (OpID 8; i=130) bitwise and reg A & val B -> reg C
	addr                 // (OpID 9; i=733) add reg A + reg C -> reg C
	addi                 // (OpID 10; id=5) add reg A + val B -> reg C
	eqrr                 // (OpID 11; i=264) eq if reg A == reg B { 1 -> reg C } else { 0 -> reg C }
	gtri                 // (OpID 12; i=20) gt if reg A > val B { 1 -> reg C } else { 0 -> reg C }
	eqir                 // (OpID 13; i=11) eq if val A == reg B { 1 -> reg C } else { 0 -> reg C }
	setr                 // (OpID 14; i=124) set reg C to reg A (copy contents)
	muli                 // (OpID 15; i=718) multiply reg A * val B -> reg C
)

// solve order: eqrr, eqri, eqir, gtri, gtrr, gtir, setr, banr, bani, seti, addr, borr, bori, muli, addi, mulr

// DetectionRecords - all of them
type DetectionRecords struct {
	Records []*DetectionRecord
}

// AddRecord - add a record
func (dr *DetectionRecords) AddRecord(b1, b2, b3, b4 uint16, opcode OpcodeId, a, b, c, a1, a2, a3, a4 uint16) {
	dr.Records = append(dr.Records, &DetectionRecord{
		Before: [4]uint16{b1, b2, b3, b4},
		Opcode: opcode,
		A:      a, B: b, C: c,
		After: [4]uint16{a1, a2, a3, a4},
	})
}

// Instruction - represents the instruction to perform
type Instruction struct {
	Opcode  OpcodeId // the opcode's ID that we'll be doing
	A, B, C uint16   //inputs A & B, output C
}

// Execute the instruction against the registers pointed to by *r, modifies
// those registers according to the result of the instruction
func (n *Instruction) Execute(r *[4]uint16) string {

	var debugDesc string

	switch n.Opcode {
	case addr:
		debugDesc = fmt.Sprintf("reg A (%d) + reg B (%d) -> reg C (%d) [%d + %d = ", n.A, n.B, n.C, (*r)[n.A], (*r)[n.B])
		(*r)[n.C] = (*r)[n.A] + (*r)[n.B]
		debugDesc = fmt.Sprintf("%s%d]\n", debugDesc, (*r)[n.C])
	case addi:
		debugDesc = fmt.Sprintf("reg A (%d) + val B (%d) -> reg C (%d) [%d + %d = ", n.A, n.B, n.C, (*r)[n.A], n.B)
		(*r)[n.C] = (*r)[n.A] + n.B
		debugDesc = fmt.Sprintf("%s%d]\n", debugDesc, (*r)[n.C])
	case mulr:
		debugDesc = fmt.Sprintf("reg A (%d) * reg B (%d) -> reg C (%d) [%d * %d = ", n.A, n.B, n.C, (*r)[n.A], (*r)[n.B])
		(*r)[n.C] = (*r)[n.A] * (*r)[n.B]
		debugDesc = fmt.Sprintf("%s%d]\n", debugDesc, (*r)[n.C])
	case muli:
		debugDesc = fmt.Sprintf("reg A (%d) * val B (%d) -> reg C (%d) [%d * %d = ", n.A, n.B, n.C, (*r)[n.A], n.B)
		(*r)[n.C] = (*r)[n.A] * n.B
		debugDesc = fmt.Sprintf("%s%d]\n", debugDesc, (*r)[n.C])
	case banr:
		debugDesc = fmt.Sprintf("reg A (%d) & reg B (%d) -> reg C (%d) [%04b & %04b = ", n.A, n.B, n.C, (*r)[n.A], (*r)[n.B])
		(*r)[n.C] = (*r)[n.A] & (*r)[n.B]
		debugDesc = fmt.Sprintf("%s %04b (%d)]\n", debugDesc, (*r)[n.C], (*r)[n.C])
	case bani:
		debugDesc = fmt.Sprintf("reg A (%d) & val B (%d) -> reg C (%d) [%04b & %04b = ", n.A, n.B, n.C, (*r)[n.A], n.B)
		(*r)[n.C] = (*r)[n.A] & n.B
		debugDesc = fmt.Sprintf("%s %04b (%d)]\n", debugDesc, (*r)[n.C], (*r)[n.C])
	case borr:
		debugDesc = fmt.Sprintf("reg A (%d) | reg B (%d) -> reg C (%d) [%04b | %04b =", n.A, n.B, n.C, (*r)[n.A], (*r)[n.B])
		(*r)[n.C] = (*r)[n.A] | (*r)[n.B]
		debugDesc = fmt.Sprintf("%s %04b (%d)]\n", debugDesc, (*r)[n.C], (*r)[n.C])
	case bori:
		debugDesc = fmt.Sprintf("reg A (%d) | val B (%d) -> reg C (%d) [%04b | %04b = ", n.A, n.B, n.C, (*r)[n.A], n.B)
		(*r)[n.C] = (*r)[n.A] | n.B
		debugDesc = fmt.Sprintf("%s %04b (%d)]\n", debugDesc, (*r)[n.C], (*r)[n.C])
	case setr:
		(*r)[n.C] = (*r)[n.A]
		debugDesc = fmt.Sprintf("copy reg A (%d) contents to reg C (%d) [%d -> ", n.A, n.C, (*r)[n.A])
		debugDesc = fmt.Sprintf("%s%d]\n", debugDesc, (*r)[n.C])
	case seti:
		debugDesc = fmt.Sprintf("copy val A (%d) to reg C (%d) [%d -> ", n.A, n.C, n.A)
		(*r)[n.C] = n.A
		debugDesc = fmt.Sprintf("%s%d]\n", debugDesc, (*r)[n.C])
	case gtir:
		debugDesc = fmt.Sprintf("if val A (%d) > reg B (%d); reg C %d = 1, else reg C (%d) = 0 [%d > %d -> ", n.A, n.B, n.C, n.C, n.A, (*r)[n.B])
		if n.A > (*r)[n.B] {
			(*r)[n.C] = 1
		} else {
			(*r)[n.C] = 0
		}
		debugDesc = fmt.Sprintf("%s%d]\n", debugDesc, (*r)[n.C])
	case gtri:
		debugDesc = fmt.Sprintf("if reg A (%d) > val B (%d); reg C %d = 1, else reg C (%d) = 0 [%d > %d -> ", n.A, n.B, n.C, n.C, (*r)[n.A], n.B)
		if (*r)[n.A] > n.B {
			(*r)[n.C] = 1
		} else {
			(*r)[n.C] = 0
		}
		debugDesc = fmt.Sprintf("%s%d]\n", debugDesc, (*r)[n.C])
	case gtrr:
		debugDesc = fmt.Sprintf("if reg A (%d) > reg B (%d); reg C %d = 1, else reg C (%d) = 0 [%d > %d -> ", n.A, n.B, n.C, n.C, (*r)[n.A], (*r)[n.B])
		if (*r)[n.A] > (*r)[n.B] {
			(*r)[n.C] = 1
		} else {
			(*r)[n.C] = 0
		}
		debugDesc = fmt.Sprintf("%s%d]\n", debugDesc, (*r)[n.C])
	case eqir:
		debugDesc = fmt.Sprintf("if val A (%d) == reg B (%d); reg C %d = 1, else reg C (%d) = 0 [%d == %d -> ", n.A, n.B, n.C, n.C, n.A, (*r)[n.B])
		if n.A == (*r)[n.B] {
			(*r)[n.C] = 1
		} else {
			(*r)[n.C] = 0
		}
		debugDesc = fmt.Sprintf("%s%d]\n", debugDesc, (*r)[n.C])
	case eqri:
		debugDesc = fmt.Sprintf("if reg A (%d) == val B (%d); reg C %d = 1, else reg C (%d) = 0 [%d == %d -> ", n.A, n.B, n.C, n.C, (*r)[n.A], n.B)
		if (*r)[n.A] == n.B {
			(*r)[n.C] = 1
		} else {
			(*r)[n.C] = 0
		}
		debugDesc = fmt.Sprintf("%s%d]\n", debugDesc, (*r)[n.C])
	case eqrr:
		debugDesc = fmt.Sprintf("if reg A (%d) == reg B (%d); reg C %d = 1, else reg C (%d) = 0 [%d == %d -> ", n.A, n.B, n.C, n.C, (*r)[n.A], (*r)[n.B])
		if (*r)[n.A] == (*r)[n.B] {
			(*r)[n.C] = 1
		} else {
			(*r)[n.C] = 0
		}
		debugDesc = fmt.Sprintf("%s%d]\n", debugDesc, (*r)[n.C])
	}
	return debugDesc
}

// Program - a program to run
type Program struct {
	Instructions []*Instruction // list of instructions
	Registers    [4]uint16      // registers for the program
}

// NewProgram - create a new Program structure
func NewProgram() *Program {
	return &Program{
		Instructions: make([]*Instruction, 0),
		Registers:    [4]uint16{0, 0, 0, 0},
	}
}

func (p *Program) Execute() {
	var s, debugDesc string
	for i, instruction := range p.Instructions {
		if *debug {
			fmt.Printf("[%04d/%04d] Executing (%02d) %s A=%d, B=%d, C=%d Input registers = %d ", i+1, len(p.Instructions), instruction.Opcode, OpCodeTypeToString(instruction.Opcode), instruction.A, instruction.B, instruction.C, p.Registers)
		}
		debugDesc = instruction.Execute(&(p.Registers))
		if *debug {
			fmt.Printf("Output registers = %d\n", p.Registers)
		}
		if *debug {
			fmt.Printf("↑ %s", debugDesc)
		}
		if *debug2 {
			fmt.Scanf("\n", &s)
		}
	}
}

// AddInstruction - add an instruction to the program
func (p *Program) AddInstruction(opcode OpcodeId, a, b, c uint16) {
	p.Instructions = append(p.Instructions, &Instruction{
		Opcode: opcode,
		A:      a, B: b, C: c,
	})
}

// DetectionRecord - A record of what was read
// [4]Before{0,1,2,3} registers 0-3.
// Opcode A B C = Instruction
type DetectionRecord struct {
	Before  [4]uint16 // Registers
	Opcode  OpcodeId  // Opcode we're doing
	A, B, C uint16    // Inputs A & B, output C
	After   [4]uint16 // Registers
}

func OpCodeTypeToString(t OpcodeId) string {
	switch t {
	case addr:
		return "addr"
	case addi:
		return "addi"
	case mulr:
		return "mulr"
	case muli:
		return "muli"
	case banr:
		return "banr"
	case bani:
		return "bani"
	case borr:
		return "borr"
	case bori:
		return "bori"
	case setr:
		return "setr"
	case seti:
		return "seti"
	case gtir:
		return "gtir"
	case gtri:
		return "gtri"
	case gtrr:
		return "gtrr"
	case eqir:
		return "eqir"
	case eqri:
		return "eqri"
	case eqrr:
		return "eqrr"
	}
	return "Unknown type"
}

// TryAll - try all opcodes with the input to see if this record behaves like 3
// or more opcodes
func (d *DetectionRecord) TryAll() bool {
	// did opcode match this record?
	opcodeMatches := make(map[OpcodeId]bool)
	matches := 0

	// Try all of the opcodes
	for i := 0; i < 16; i++ {
		switch OpcodeId(i) {
		case addr:
			opcodeMatches[OpcodeId(i)] = (d.After[d.C] == d.Before[d.A]+d.Before[d.B])
			if *debug {
				fmt.Printf("\t%s: reg A (%d) + reg B (%d) -> reg C (%d) [does %d + %d = %d]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C,
					d.Before[d.A], d.Before[d.B], d.After[d.C],
				)
			}
		case addi:
			opcodeMatches[OpcodeId(i)] = (d.After[d.C] == d.Before[d.A]+d.B)
			if *debug {
				fmt.Printf("\t%s: reg A (%d) + val B (%d) -> reg C (%d) [does %d + %d = %d]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C,
					d.Before[d.A], d.B, d.After[d.C])
			}
		case mulr:
			opcodeMatches[OpcodeId(i)] = (d.After[d.C] == d.Before[d.A]*d.Before[d.B])
			if *debug {
				fmt.Printf("\t%s: reg A (%d) * reg B (%d) -> reg C (%d) [does %d * %d = %d]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C,
					d.Before[d.A], d.Before[d.B], d.After[d.C])
			}
		case muli:
			opcodeMatches[OpcodeId(i)] = (d.After[d.C] == d.Before[d.A]*d.B)
			if *debug {
				fmt.Printf("\t%s: reg A (%d) * val B (%d) -> reg C (%d) [does %d * %d = %d] \n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C,
					d.Before[d.A], d.B, d.After[d.C])
			}
		case banr:
			opcodeMatches[OpcodeId(i)] = (d.After[d.C] == d.Before[d.A]&d.Before[d.B])
			if *debug {
				fmt.Printf("\t%s: reg A (%d) & reg B (%d) -> reg C (%d) [does %04b & %04b = %04b]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C,
					d.Before[d.A], d.Before[d.B], d.After[d.C])
			}
		case bani:
			opcodeMatches[OpcodeId(i)] = (d.After[d.C] == d.Before[d.A]&d.B)
			if *debug {
				fmt.Printf("\t%s: reg A (%d) & val B (%d) -> reg C (%d) [does %04b & %04b = %04b]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C,
					d.Before[d.A], d.B, d.After[d.C])
			}
		case borr:
			opcodeMatches[OpcodeId(i)] = (d.After[d.C] == d.Before[d.A]|d.Before[d.B])
			if *debug {
				fmt.Printf("\t%s: reg A (%d) | reg B (%d) -> reg C (%d) [does %04b | %04b = %04b]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C,
					d.Before[d.A], d.Before[d.B], d.After[d.C])
			}
		case bori:
			opcodeMatches[OpcodeId(i)] = (d.After[d.C] == d.Before[d.A]|d.B)
			if *debug {
				fmt.Printf("\t%s: reg A (%d) | val B (%d) -> reg C (%d) [does %04b | %04b = %04b]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C,
					d.Before[d.A], d.B, d.After[d.C])
			}
		case setr:
			opcodeMatches[OpcodeId(i)] = (d.After[d.C] == d.Before[d.A])
			if *debug {
				fmt.Printf("\t%s: copy reg A (%d) contents to reg C (%d) [does %d = %d]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.C,
					d.After[d.C], d.Before[d.A])
			}
		case seti:
			opcodeMatches[OpcodeId(i)] = (d.After[d.C] == d.A)
			if *debug {
				fmt.Printf("\t%s: copy val A (%d) to reg C (%d) [does %d = %d]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.Before[d.A], d.C,
					d.After[d.C], d.A)
			}
		case gtir:
			opcodeMatches[OpcodeId(i)] = false
			if d.A > d.Before[d.B] {
				if d.After[d.C] == 1 {
					opcodeMatches[OpcodeId(i)] = true
				}
			} else {
				if d.After[d.C] == 0 {
					opcodeMatches[OpcodeId(i)] = true
				}
			}
			// opcodeMatches[OpcodeId(i)] = (d.After[d.C] == 1 && d.A > d.Before[d.B])
			if *debug {
				fmt.Printf("\t%s: if val A (%d) > reg B (%d), reg C (%d) = 1, else reg C (%d) = 0 [is %d > %d? reg C val=%d]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C, d.C,
					d.A, d.Before[d.B], d.After[d.C])
			}
		case gtri:
			opcodeMatches[OpcodeId(i)] = false
			if d.Before[d.A] > d.B {
				if d.After[d.C] == 1 {
					opcodeMatches[OpcodeId(i)] = true
				}
			} else {
				if d.After[d.C] == 0 {
					opcodeMatches[OpcodeId(i)] = true
				}
			}
			if *debug {
				fmt.Printf("\t%s: if reg A (%d) > val B (%d), reg C (%d) = 1, else reg C (%d) = 0 [is %d > %d? reg C val=%d]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C, d.C,
					d.Before[d.A], d.B, d.After[d.C])
			}
		case gtrr:
			opcodeMatches[OpcodeId(i)] = false
			if d.Before[d.A] > d.Before[d.B] {
				if d.After[d.C] == 1 {
					opcodeMatches[OpcodeId(i)] = true
				}
			} else {
				if d.After[d.C] == 0 {
					opcodeMatches[OpcodeId(i)] = true
				}
			}
			if *debug {
				fmt.Printf("\t%s: if reg A (%d) > reg B (%d), reg C (%d) = 1, else reg C (%d) = 0 [is %d > %d? reg C val=%d]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C, d.C,
					d.Before[d.A], d.Before[d.B], d.After[d.C])
			}
		case eqir:
			opcodeMatches[OpcodeId(i)] = false
			if d.A == d.Before[d.B] {
				if d.After[d.C] == 1 {
					opcodeMatches[OpcodeId(i)] = true
				}
			} else {
				if d.After[d.C] == 0 {
					opcodeMatches[OpcodeId(i)] = true
				}
			}
			if *debug {
				fmt.Printf("\t%s: if val A (%d) == reg B (%d), reg C (%d) = 1, else reg C (%d) = 0 [is %d == %d? reg C=%d]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C, d.C,
					d.A, d.Before[d.B], d.After[d.C])
			}
		case eqri:
			opcodeMatches[OpcodeId(i)] = false
			if d.Before[d.A] == d.B {
				if d.After[d.C] == 1 {
					opcodeMatches[OpcodeId(i)] = true
				}
			} else {
				if d.After[d.C] == 0 {
					opcodeMatches[OpcodeId(i)] = true
				}
			}
			if *debug {
				fmt.Printf("\t%s: if reg A (%d) == val B (%d), reg C (%d) = 1, else reg C (%d) = 0 [is %d == %d? reg C=%d]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C, d.C,
					d.Before[d.A], d.B, d.After[d.C])
			}
		case eqrr:
			opcodeMatches[OpcodeId(i)] = false
			if d.Before[d.A] == d.Before[d.B] {
				if d.After[d.C] == 1 {
					opcodeMatches[OpcodeId(i)] = true
				}
			} else {
				if d.After[d.C] == 0 {
					opcodeMatches[OpcodeId(i)] = true
				}
			}
			if *debug {
				fmt.Printf("\t%s: if reg A (%d) == reg B (%d), reg C (%d) = 1, else reg C (%d) = 0 [is %d == %d? reg C=%d]\n", printAndColourize(OpcodeId(i), opcodeMatches[OpcodeId(i)]),
					d.A, d.B, d.C, d.C,
					d.Before[d.A], d.Before[d.B], d.After[d.C])
			}
		}
		// how many opcodes did this record behave like?
		if opcodeMatches[OpcodeId(i)] {
			matches++
		}
	}
	if *debug {
		fmt.Printf("\tMatches = %d\n", matches)
	}
	return matches >= 3
}

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n", e.Error())
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	inputBuffer, err := ioutil.ReadFile(*inputFile)
	errorIf("couldn't read input file", err)

	if !*partB {
		matches := detectionMatcher.FindAllStringSubmatch(string(inputBuffer), -1)
		records := new(DetectionRecords)

		for matchGroupIdx := 0; matchGroupIdx < len(matches); matchGroupIdx++ {

			if *debug {
				fmt.Printf("Match %d\n", matchGroupIdx)
			}
			b1, err := strconv.Atoi(matches[matchGroupIdx][1])
			errorIf("Couldn't parse b1", err)
			b2, err := strconv.Atoi(matches[matchGroupIdx][2])
			errorIf("Couldn't parse b2", err)
			b3, err := strconv.Atoi(matches[matchGroupIdx][3])
			errorIf("Couldn't parse b3", err)
			b4, err := strconv.Atoi(matches[matchGroupIdx][4])
			errorIf("Couldn't parse b4", err)

			op, err := strconv.Atoi(matches[matchGroupIdx][5])
			errorIf("Couldn't parse op", err)
			a, err := strconv.Atoi(matches[matchGroupIdx][6])
			errorIf("Couldn't parse a", err)
			b, err := strconv.Atoi(matches[matchGroupIdx][7])
			errorIf("Couldn't parse b", err)
			c, err := strconv.Atoi(matches[matchGroupIdx][8])
			errorIf("Couldn't parse c", err)

			a1, err := strconv.Atoi(matches[matchGroupIdx][9])
			errorIf("Couldn't parse a1", err)
			a2, err := strconv.Atoi(matches[matchGroupIdx][10])
			errorIf("Couldn't parse a2", err)
			a3, err := strconv.Atoi(matches[matchGroupIdx][11])
			errorIf("Couldn't parse a3", err)
			a4, err := strconv.Atoi(matches[matchGroupIdx][12])
			errorIf("Couldn't parse a4", err)
			records.AddRecord(uint16(b1), uint16(b2), uint16(b3), uint16(b4),
				OpcodeId(op), uint16(a), uint16(b), uint16(c),
				uint16(a1), uint16(a2), uint16(a3), uint16(a4))
		}

		total := 0
		for i, record := range records.Records {
			if *debug {
				fmt.Printf("[%05d/%05d] Before:[%d, %d, %d, %d]; OpID=%d, A=%d, B=%d, C=%d; After:[%d, %d, %d, %d]\n", i+1, len(records.Records),
					record.Before[0], record.Before[1], record.Before[2], record.Before[3], record.Opcode, record.A, record.B, record.C, record.After[0], record.After[1], record.After[2], record.After[3])
			}
			if record.TryAll() {
				total++
			}
		}
		fmt.Printf("Total hits %d\n", total)
	} else {
		//part B
		program := NewProgram()
		for _, line := range strings.Split(string(inputBuffer), "\n") {
			var op, a, b, c int
			var err error
			for i, token := range strings.Split(line, " ") {
				switch i {
				case 0:
					op, err = strconv.Atoi(token)
					errorIf("Couldn't parse op", err)
				case 1:
					a, err = strconv.Atoi(token)
					errorIf("Couldn't parse a", err)
				case 2:
					b, err = strconv.Atoi(token)
					errorIf("Couldn't parse b", err)
				case 3:
					c, err = strconv.Atoi(token)
					errorIf("Couldn't parse c", err)
				}
			}
			program.AddInstruction(OpcodeId(op), uint16(a), uint16(b), uint16(c))
		}
		if *debug {
			fmt.Printf("Program has %d lines\n", len(program.Instructions))
		}
		program.Execute()
		fmt.Printf("Program complete. Registers: %d\n", program.Registers)
	}
}
