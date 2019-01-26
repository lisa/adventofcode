package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	. "github.com/logrusorgru/aurora"
)

var (
	inputFile        = flag.String("input", "inputs/day16-detection.txt", "input file")
	partB            = flag.Bool("partB", false, "do part b solution?")
	debug            = flag.Bool("debug", false, "debug?")
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
	addr OpcodeId = iota // add reg A + reg C -> reg C
	addi                 // add reg A + val B -> reg C
	mulr                 // multiply reg A * reg B -> reg C
	muli                 // multiply reg A * val B -> reg C
	banr                 // bitwise and reg A & reg b -> reg C
	bani                 // bitwise and reg A & val B -> reg C
	borr                 // bitwise or reg A | reg B -> reg C
	bori                 // bitwise or reg A | val B -> reg C
	setr                 // set reg C to reg A (copy contents)
	seti                 // set reg C to val A (copy value)
	gtir                 // gt if val A > reg B { 1 -> reg C } else { 0 -> reg C }
	gtri                 // gt if reg A > val B { 1 -> reg C } else { 0 -> reg C }
	gtrr                 // gt if reg A > reg B { 1 -> reg C } else { 0 -> reg C }
	eqir                 // eq if val A == reg B { 1 -> reg C } else { 0 -> reg C }
	eqri                 // eq if reg A == val B { 1 -> reg C } else { 0 -> reg C }
	eqrr                 // eq if reg A == reg B { 1 -> reg C } else { 0 -> reg C }
)

// DetectionRecords - all of them
type DetectionRecords struct {
	Records []*DetectionRecord
}

// AddRecord - add a record
func (dr *DetectionRecords) AddRecord(b1, b2, b3, b4 uint8, opcode OpcodeId, a, b, c, a1, a2, a3, a4 uint8) {
	dr.Records = append(dr.Records, &DetectionRecord{
		Before: [4]uint8{b1, b2, b3, b4},
		Opcode: opcode,
		A:      a, B: b, C: c,
		After: [4]uint8{a1, a2, a3, a4},
	})
}

// DetectionRecord - A record of what was read
// [4]Before{0,1,2,3} registers 0-3.
// Opcode A B C = Instruction
type DetectionRecord struct {
	Before  [4]uint8 // Registers
	Opcode  OpcodeId // Opcode we're doing
	A, B, C uint8    // Inputs A & B, output C
	After   [4]uint8 // Registers
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
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	inputBuffer, err := ioutil.ReadFile(*inputFile)
	errorIf("couldn't read input file", err)

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

		records.AddRecord(uint8(b1), uint8(b2), uint8(b3), uint8(b4),
			OpcodeId(op), uint8(a), uint8(b), uint8(c),
			uint8(a1), uint8(a2), uint8(a3), uint8(a4))
	}

	//fmt.Printf("records %+v\n", records.Records[0])
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

}
