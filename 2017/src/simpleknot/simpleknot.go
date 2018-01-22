package simpleknot

import (
	"fmt"
)

// append this suffix to the rawInput
var hashSufix = []byte{17, 31, 73, 47, 23}

type Hash struct {
	data          []byte // the hashing structure
	rawInput      []byte // input key, slice of ints representing bytes (we'll append the suffix)
	appendedInput []byte // input with appended suffix
}

func (h *Hash) String() string {
	return h.DenseHashToString()
}

func (h *Hash) DenseHashToString() string {
	return fmt.Sprintf("%02x", h.ComputeDenseHash())
}

// Return a 16 byte slice representing the hex digits of the dense hash based on the input
func (h *Hash) ComputeDenseHash() []byte {
	// do rounds
	workingList := h.data
	currentIndex := 0
	skipSize := 0
	ret := make([]byte, 16)
	for i := 0; i < 64; i++ {
		doRound(h.appendedInput, &currentIndex, &skipSize, &workingList)
	}

	currentIndex = 0
	// chunk it up into our dense hash
	for chunk := 0; chunk < 16; chunk++ {
		ret[chunk] = workingList[currentIndex]
		currentIndex += 1
		for digit := 1; digit < 16; digit++ {
			ret[chunk] ^= workingList[currentIndex]
			currentIndex += 1
		} // finsihed 16 digits
	} // done with the chunks

	return ret
}

// Reverse `l` consecutive bytes in `data`, starting from the `start`th index,
// with intelligent wrapping.
func reverseSlice(data []byte, start, l int) []byte {
	reverseSlice := make([]byte, 0)
	reversed := data

	// begin the index looping at (start+l) and perform l iterations backwards (ie, back to `start`)
	for i := (start + l) - 1; i >= start; i-- {
		reverseSlice = append(reverseSlice, data[i%len(data)])
	}
	// stitch back together
	for i := 0; i < len(reverseSlice); i++ {
		reversed[(start+i)%len(data)] = reverseSlice[i]
	}
	return reversed
}

// Do a round of manipulation.
func doRound(inputLengths []byte, currentIndex, skipSize *int, data *[]byte) []byte {
	workingData := *data
	for number := 0; number < len(inputLengths); number++ {
		workingData = reverseSlice(workingData, *currentIndex, int(inputLengths[number]))
		*currentIndex += (*skipSize + int(inputLengths[number])) % len(workingData)
		// increase skipSize by 1 since we've completed a round.
		*skipSize += 1
	} // done with the input
	return workingData
}

// Return a 256 byte-long slice with digits 0-255 inclusive, in order.
func initialize() []byte {
	ret := make([]byte, 256)
	for i := 0; i < 256; i++ {
		ret[i] = byte(i)
	}
	return ret
}

func New(input []byte) *Hash {
	return &Hash{
		data:          initialize(),
		rawInput:      input,
		appendedInput: append(input,hashSuffix...),
	}
}
