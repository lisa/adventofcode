package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	inputFile = flag.String("input", "inputs/day09.txt", "input file")
	partB     = flag.Bool("partB", false, "do part b solution?")
	debug     = flag.Bool("debug", false, "debug?")
	debug2    = flag.Bool("debug2", false, "more debug")
	removed   = make(map[int]int)
)

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}

// based on https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
// why oh why does Go behave differently :'(
// ruby: -1 % 4 == 3
// go:   -1 % 4 == -1
func rubyMod(d, m int) int {
	r := d % m
	if (r < 0 && m > 0) || (r > 0 && m < 0) {
		r += m
	}
	return r
}

// Len - how many marbles in the circle?
func (c *Circle) Len() int {
	return int(len(c.Marbles))
}

// Circle - where all the marbles go
type Circle struct {
	CurrentMarbleIndex int
	Marbles            []int
}

func NewCircle() *Circle {
	c := Circle{
		CurrentMarbleIndex: 0,
		Marbles:            []int{0},
	}
	c.Marbles[0] = 0
	return &c
}

// InsertMarble - inserts a marble according to the rules
func (c *Circle) InsertMarble(value int) {
	// Insert the marble between the marbles 1 and 2 clockwise (left) of the current index
	// [a] [b] [c] (d) [e] [f] [g] d = current
	// Insert between [e] and [f]
	oneEdge := rubyMod(c.CurrentMarbleIndex+1, c.Len())
	twoEdge := rubyMod(c.CurrentMarbleIndex+2, c.Len())
	var nextIndex int
	var leftPart, rightPart []int
	switch len(c.Marbles) {
	case 1:
		c.Marbles = []int{0, 1}
		c.CurrentMarbleIndex = 1
		return
	case 2:
		leftPart = []int{c.Marbles[0]}
		rightPart = []int{c.Marbles[1]}
		nextIndex = 1
	case 3:
		leftPart = []int{c.Marbles[0], c.Marbles[1], c.Marbles[2]}
		rightPart = []int{}
		nextIndex = 3
	default:
		if twoEdge < oneEdge {
			leftPart = c.Marbles[0:twoEdge]
			rightPart = c.Marbles[twoEdge : oneEdge+1]
			nextIndex = twoEdge
		} else {
			leftPart = c.Marbles[0 : oneEdge+1]
			rightPart = c.Marbles[oneEdge+1:]
			nextIndex = rubyMod(oneEdge+1, c.Len()+1)
		}
	}
	// insert the new marble
	c.Marbles = append(leftPart, append([]int{value}, rightPart...)...)
	c.CurrentMarbleIndex = nextIndex
	return
}

// RemoveAtIndex - removes the marble at index i and returns its value.
// The marble to the right (clockwise) of the removed marble becomes the current index
func (c *Circle) RemoveAtIndex(i int) int {
	v := c.Marbles[i]
	newIndex := c.FindIndexByOffset(-6)
	if newIndex == 0 {
		// we're removing the last element so let's help out a little since the math
		// below doesn't work
		c.Marbles = c.Marbles[0 : c.Len()-1]
		c.CurrentMarbleIndex = 0
		return v
	}
	c.CurrentMarbleIndex = newIndex
	// now cut it out
	leftSide := c.Marbles[:rubyMod(c.CurrentMarbleIndex-1, c.Len())]
	rightSide := c.Marbles[c.CurrentMarbleIndex:]
	c.Marbles = append(leftSide, rightSide...)
	c.MoveLeft()
	return v
}

// FindIndexByOffset - find the index if we would move +offset+.
func (c *Circle) FindIndexByOffset(offset int) int {
	if offset == 0 {
		return c.CurrentMarbleIndex
	}
	return rubyMod(c.CurrentMarbleIndex+offset, int(len(c.Marbles)))
}

// Move offset % len(circle), move right +, move left -
// "Move" means to move the CurrentMarbleIndex by that many.
func (c *Circle) Move(offset int) {
	c.CurrentMarbleIndex = c.FindIndexByOffset(offset)
}

// MoveLeft - move the index one to the left (counter-clockwise)
// oddly we don't ever need to go Right.
func (c *Circle) MoveLeft() {
	c.Move(-1)
}

// Mod23 - Do the work necessary to handle when the added value is a multiple of
// 23. Remove the marble 7 counter-clockwise (left) of the current marble
// (return it) The marble to the right (clockwise) of the removed marble is the
// new current.
func (c *Circle) Mod23() int {
	// [a] [b] [c] [d] [e] [f] [g] [h] [i] [j] (k) [l] [m] [n]
	// Remove [d]
	// New index is [e]
	v := c.RemoveAtIndex(c.FindIndexByOffset(-7))
	return v
}

func (c *Circle) String() string {
	ret := ""
	var i int
	for i = 0; i < c.Len(); i++ {
		if i == c.CurrentMarbleIndex {
			ret += fmt.Sprintf(" (%d) ", c.Marbles[i])
		} else {
			ret += fmt.Sprintf(" %d ", c.Marbles[i])
		}
	}
	return ret
}

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	errorIf("Can't open input file", err)

	defer input.Close()
	lineReader := bufio.NewScanner(input)
	var players, lastMarbleValue int
	if lineReader.Scan() {
		words := strings.Split(lineReader.Text(), " ")
		players, err = strconv.Atoi(words[0])
		errorIf("Couldn't parse the number of players\n", err)
		lastMarbleValue, err = strconv.Atoi(words[6])
		errorIf("Couldn't parse last marble score\n", err)
	}

	fmt.Printf("Player count %d, highest marble value %d\n", players, lastMarbleValue)

	score := make([]int, int(players))

	circle := NewCircle()

	var v int
	for v = 1; v <= int(lastMarbleValue); v++ {
		if v%23 == 0 {

			seventh := circle.Mod23()

			removed[v] = seventh
			score[v%int(players)] += seventh + v
		} else {
			circle.InsertMarble(v)
		}
	}

	highest := -1
	for i := range score {
		if score[i] > highest {
			highest = score[i]
		}
	}
	fmt.Printf("high score = %d\n", highest)
}
