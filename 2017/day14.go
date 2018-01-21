package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"simpleknot"
)

var input = flag.String("input", "flqrgnkx", "Puzzle input")
var debug = flag.Bool("debug", false, "Debug?")

// For printing the fragmentation
func ActiveBitsToString(hash *simpleknot.Hash) string {
	ret := ""
	for _, char := range strings.Split(hash.DenseHashToString(), "") {
		hexDigit, err := strconv.ParseInt(char, 16, 8) // should be a
		if err != nil {
			fmt.Printf("Couldn't convert %s to a hex digit somehow...\n", char)
			os.Exit(1)
		}

		hexDigits := HexDigitToBits(hexDigit)

		for i := 0; i < 4; i++ {
			if hexDigits[i] == 1 {
				ret += "#"
			} else {
				ret += "."
			}
		}
	}
	return ret
}

func CountActiveBits(hash *simpleknot.Hash) int {
	ret := 0
	for _, char := range strings.Split(hash.DenseHashToString(), "") {
		hexDigit, err := strconv.ParseInt(char, 16, 8) // should be a
		if err != nil {
			fmt.Printf("Couldn't convert %s to a hex digit somehow...\n", char)
			os.Exit(1)
		}

		hexDigits := HexDigitToBits(hexDigit)

		for i := 0; i < 4; i++ {
			if hexDigits[i] == 1 {
				ret += 1
			}
		}
	}
	return ret
}

/* This is probably a crime somewhere... */
func HexDigitToBits(h int64) [4]int8 {
	bits := make(map[int64][4]int8)
	bits[0] = [4]int8{0, 0, 0, 0}
	bits[1] = [4]int8{0, 0, 0, 1}
	bits[2] = [4]int8{0, 0, 1, 0}
	bits[3] = [4]int8{0, 0, 1, 1}
	bits[4] = [4]int8{0, 1, 0, 0}
	bits[5] = [4]int8{0, 1, 0, 1}
	bits[6] = [4]int8{0, 1, 1, 0}
	bits[7] = [4]int8{0, 1, 1, 1}
	bits[8] = [4]int8{1, 0, 0, 0}
	bits[9] = [4]int8{1, 0, 0, 1}
	bits[10] = [4]int8{1, 0, 1, 0}
	bits[11] = [4]int8{1, 0, 1, 1}
	bits[12] = [4]int8{1, 1, 0, 0}
	bits[13] = [4]int8{1, 1, 0, 1}
	bits[14] = [4]int8{1, 1, 1, 0}
	bits[15] = [4]int8{1, 1, 1, 1}

	return bits[h]
}

func stringToChars(str string) []byte {
	ret := make([]byte, len(str))
	for i, char := range strings.Split(str, "") {
		ret[i] = byte(char[0])
	}
	return ret
}

func main() {
	flag.Parse()
	fmt.Printf("Input: %s\n", *input)
	hashes := make(map[int]*simpleknot.Hash) // Row and hash
	usedSquares := 0
	for i := 0; i < 128; i++ {
		hashes[i] = simpleknot.New(stringToChars(fmt.Sprintf("%s-%d", *input, i)))
		usedSquares += CountActiveBits(hashes[i])
		if *debug {
			fmt.Printf("[%03d/128] %s - %02x (usedSquares=%d)\n", i, ActiveBitsToString(hashes[i]), hashes[i].ComputeDenseHash(), usedSquares)
		}

	}
	fmt.Printf("Used squares: %d\n", usedSquares)

}
