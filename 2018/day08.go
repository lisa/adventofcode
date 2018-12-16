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
	debug2    = flag.Bool("debug2", false, "more debugging?")
	debug3    = flag.Bool("debug3", false, "require user input to advance parse loop?")
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

type Node struct {
	Children       []*Node
	MetadataCount  int // number of metadata entries in total
	Metadata       []int
	Parent         *Node
	Value          int
	MetadataToRead int // number of metadata entries left to read
}

func NewNode() *Node {
	return &Node{
		Children:       make([]*Node, 0),
		Metadata:       make([]int, 0),
		Parent:         nil,
		Value:          0,
		MetadataCount:  0,
		MetadataToRead: 0,
	}
}

// SumTreeMetadata - sum up my metadata and my childrens'
func (n *Node) SumTreeMetadata() int {
	s := n.MetadataSum()
	for _, c := range n.Children {
		s += c.SumTreeMetadata()
	}
	return s
}

// MetadataSum - sum up all my metadata
func (n *Node) MetadataSum() int {
	s := 0
	for _, v := range n.Metadata {
		s += v
	}
	return s
}

// PartBSum - compute the part B sum for day 8
func (n *Node) PartBSum() int {
	if *debug {
		fmt.Printf("Part B sum: %+v I have %d children ", n, len(n.Children))
	}
	s := 0
	if len(n.Children) == 0 {
		if *debug {
			fmt.Printf(" and that means my value is my metadata sum %d\n", n.MetadataSum())
		}
		s = n.MetadataSum()
	} else {
		if *debug {
			fmt.Printf("And that means my sum is the sum of my children at my indexes: %d\n", n.Metadata)
		}
		for _, i := range n.Metadata {
			if *debug {
				fmt.Printf("  metadata value: %d\n", i)
			}
			// okay to use len() and i here like this since they're both 1-based
			// however, in indexing into our Children slice we need to (i-1)
			if len(n.Children) >= i {
				s += n.Children[i-1].PartBSum()
			}
		}
	}
	return s
}

// recursively parse the input
// *i needs to start as -1
func parseInputIntoNodes(root *Node, i *int, dataSlice []int) {
	if *debug {
		fmt.Printf("data %d\n", dataSlice)
	}
	state := parseChildren
	for *i < len(dataSlice) {
		// first thing's first, advance the counter
		*i++
		if *debug {
			fmt.Println()
			fmt.Printf("[%5d/%5d], d=%d, state=%d", *i, len(dataSlice)-1, dataSlice[*i], state)
		}

		switch state {
		case parseChildren:
			// Allocate room for my children
			if *debug {
				fmt.Printf(" parseChildren: Allocating %d Children for (%p) root %+v\n", dataSlice[*i], root, root)
			}
			root.Children = make([]*Node, dataSlice[*i])
			// and now, add "stubs" for my children
			for c := range root.Children {
				n := NewNode()
				n.Parent = root
				root.Children[c] = n
			}
			state = metadataCount
		case metadataCount:
			// I have this many metadata entries in total
			if *debug {
				fmt.Printf(" metadataCount: I have to read %d entries\n", dataSlice[*i])
			}
			root.MetadataCount = dataSlice[*i]
			root.MetadataToRead = dataSlice[*i]
			// Allocate room for my metadata entries
			root.Metadata = make([]int, dataSlice[*i])

			// What should I do next? If I have children I need to read them in before I
			// read my own metadata
			if len(root.Children) > 0 {
				// I should loop through my kids to read them in
				if *debug {
					fmt.Printf("   metadataCount: %d children to read, so looping...\n", len(root.Children))
				}
				for _, child := range root.Children {
					// before passing control to the next iteration, need to advance to the next
					// item in the dataSlice
					if *debug {
						fmt.Printf("   metadataCount: About to pass control for (%p) child %+v\n", child, child)
					}
					parseInputIntoNodes(child, i, dataSlice)
					if *debug {
						fmt.Printf("   metadataCount: Finished processing (%p) child %+v\n", child, child)
					}
					// the final child will fall through the select statement and advance *i
					// afterwards.
				}
			}
			// At this point one of two things is true:
			// 1) I have finished reading all of my children and now need to read my own metadata
			// 2) I never had any children and now need to read my own metadata
			state = metadataRead
		case metadataRead:
			if *debug {
				fmt.Printf(" metadataRead: Reading metadata value %d\n", dataSlice[*i])
			}
			// read metadata for myself
			root.Metadata[root.MetadataCount-root.MetadataToRead] = dataSlice[*i]
			root.MetadataToRead--
			// Do I have more metadata?
			if *debug {
				fmt.Printf("   metadataRead: More to go if MetadataToRead > 0 (%d > 0)\n", root.MetadataToRead)
			}
			if root.MetadataToRead > 0 {
				// More metadata to read
				// (this state "change" isn't strictly necessary but it will help readability.)
				state = metadataRead
			} else {
				// no more metadata for myself, so I need to break out of the for loop.
				return
			}
		}
	}
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
	// must be -1 due to pre-increment
	i := -1
	root := NewNode()
	parseInputIntoNodes(root, &i, dataset)

	if !*partB {
		// sum all of the entries
		fmt.Printf("tree sum= %d\n", root.SumTreeMetadata())
	} else {
		// part B
		fmt.Printf("part b sum = %d\n", root.PartBSum())
	}

}
