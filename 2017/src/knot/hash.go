package knot

import (
	"container/ring"
)

// Append this to the input after it has been converted to ASCII.
var hashSufix = []int{17, 31, 73, 47, 23}

type Hash struct {
	ring          *ring.Ring // underlying ring structure
	rawInput      []int      // input key, slice of ints representing bytes (we'll append the suffix)
	appendedInput []int      // input with appended suffix
}

// Compute the Dense Hash
func (h *Hash) ComputeDenseHash() []byte {
	// Need to do the 64 rounds

	skipSize := 0
	totalSkips := 0
	for round := 0; round < 64; round++ {
		h.ring = doRound(h.appendedInput, &skipSize, &totalSkips, h.ring)
	}
	h.ring = h.ring.Move(-1 * totalSkips)

	ret := make([]int, 16)
	for chunk := 0; chunk < 16; chunk++ {
		// "Seed" the chunk bits with the first value of the 16 digits for ^=
		ret[chunk] = h.ring.Value.(int)

		h.ring = h.ring.Next()
		for digit := 1; digit < 16; digit++ {
			// The digit-th digit in the dense hash
			ret[chunk] ^= h.ring.Value.(int)
			h.ring = h.ring.Next()
		} // finsihed 16 digits
	} // done with the chunks

	// coerce to []byte
	byteRet := make([]byte, len(ret))
	for i := 0; i < len(ret); i++ {
		byteRet[i] = byte(ret[i])
	}
	return byteRet
}

func reverseRingSlice(r *ring.Ring, sliceLen int) *ring.Ring {
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

/* Performs a single round */
func doRound(inputLengths []int, skipSize, totalSkips *int, ring *ring.Ring) *ring.Ring {

	for number := 0; number < len(inputLengths); number++ {
		ring = reverseRingSlice(ring, inputLengths[number])
		ring = ring.Move(inputLengths[number] + *skipSize)
		// Save the total number of skips for later rewinding
		*totalSkips += inputLengths[number] + *skipSize
		// Then increase skipSize
		*skipSize += 1
	} // done with input
	return ring
}

func New(input []int) *Hash {
	appendedinputs := make([]int, 0)
	for i := 0; i < len(input); i++ {
		appendedinputs = append(appendedinputs, input[i])
	}
	for i := 0; i < len(hashSufix); i++ {
		appendedinputs = append(appendedinputs, hashSufix[i])
	}

	h := Hash{
		ring:          ring.New(256),
		rawInput:      input,
		appendedInput: appendedinputs,
	}
	for i := 0; i <= 255; i++ {
		h.ring.Value = i
		h.ring = h.ring.Next()
	}

	return &h
}
