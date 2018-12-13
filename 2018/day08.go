package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	inputFile = flag.String("input", "inputs/day08.txt", "Input file")
	partB     = flag.Bool("partB", false, "Perform part B?")
	debug     = flag.Bool("debug", false, "Debug?")
)

const (
	parseChildren int = iota
	metadataCount
	metadataRead
)

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}

// parseInput - takes the node to which we append other nodes/data, and the piece of data to work with
func parseInput(dataSlice []int) []int {
	// we have to "pause" parsing the current Node to look, potentially, at its
	// children. we'll use this to know how far back we have to go to get to +root+
	depth := 0
	entries := make([]int, 0)

	metadataEntries := make(map[int]int)
	childrenToRead := make(map[int]int)
	state := parseChildren
	for i := 0; i < len(dataSlice); i++ {
		if *debug {
			fmt.Printf("[%5d/%5d], d=%d: state=%d\n", i, len(dataSlice)-1, dataSlice[i], state)
		}
		switch state {
		case parseChildren:
			// Get the number of child nodes
			if dataSlice[i] > 0 {
				// there are children to read in a lower layer, so send that down
			}
			childrenToRead[depth+1] = dataSlice[i]
			if *debug {
				fmt.Printf("  parseChildren: Expecting to read %d children in depth %d.\n", childrenToRead[depth+1], depth+1)
			}
			state = metadataCount
		case metadataCount:
			// how many metadata for this depth?
			metadataEntries[depth] = dataSlice[i]
			if *debug {
				fmt.Printf("  metadataCount: Expecting to read %d metadata entries for depth %d\n", metadataEntries[depth], depth)
			}
			if childrenToRead[depth+1] == 0 {
				// no children in the next depth down, so read our metadata
				state = metadataRead
			} else {
				// go down one level
				if *debug {
					fmt.Printf("Going from depth %d to %d\n", depth, depth+1)
				}
				depth++
				state = parseChildren
			}
		case metadataRead:
			// read metadata for this depth
			if *debug {
				fmt.Printf("  metadataRead: We have %d entries left to read at depth %d\n", metadataEntries[depth], depth)
			}

			if metadataEntries[depth] > 0 {
				entries = append(entries, dataSlice[i])
				metadataEntries[depth]--
				if *debug {
					fmt.Printf("  metadataRead: Reading a metadata entry %d at depth %d (%d left to read)\n", dataSlice[i], depth, metadataEntries[depth])
				}
			}
			// check to see if we're done with reading
			if metadataEntries[depth] == 0 {
				// we have successfully completed reading a child, so there's one less
				childrenToRead[depth]--
				if childrenToRead[depth] == 0 {
					if *debug {
						fmt.Printf("  metadataRead: There are no children left to read at depth %d, so next up is metadata\n", depth)
					}
					// read metadata at higher depth
					depth--
					state = metadataRead
				} else {
					if *debug {
						fmt.Printf("  metadaRead: There's more children at depth %d, so read them next\n", depth)
					}
					// there's more children right away
					state = parseChildren
				}
			}
		}
	}
	return entries
}

func main() {
	flag.Parse()
	input, err := os.Open(*inputFile)
	errorIf("Can't open input file", err)

	defer input.Close()
	lineReader := bufio.NewScanner(input)

	lineReader.Split(bufio.ScanWords)

	dataset := make([]int, 0)
	for lineReader.Scan() {
		number := lineReader.Text()
		d, err := strconv.Atoi(number)
		errorIf("Couldn't parse a digit\n", err)
		dataset = append(dataset, d)
	}

	res := parseInput(dataset)

	sum := 0
	for i := 0; i < len(res); i++ {
		sum = sum + res[i]
	}
	fmt.Printf("Sum: %d\n", sum)

}
