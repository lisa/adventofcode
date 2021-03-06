package main

/* Part A:

Hat tip to @lizthegrey for pointing me in the right direction

Given a "spiral pattern" (see below) what is the Manhattan Distance from cell 1 to the specified cell (368078)?
Data from cell 1 is 0 steps, 12 is 3 steps, 23 is 2, 1024 is 31 steps. ex.


(Y= 2) 17  16  15  14  13
(Y= 1) 18   5   4   3  12
(Y= 0) 19   6   1   2  11
(Y=-1) 20   7   8   9  10
(Y=-2) 21  22  23---> ...
X=     -2  -1   0   1   2
^
| Axis
Y  X->

NOTE: There is no Y axis "column" for 0.

Continually build out to know how to get back, but before taking a step (creating a coordinate), inspect the
"places I've been" list to see if I can make a left turn or if I need to keep going straight.


Movement Rules:
Rule 1. if !seen(last.X-1,last.Y) && !seen(last.X,last.Y-1) (use (last.X,last.Y-1) unless (last.X+1,last.Y) is unseen in which case follow rule 4)
Rule 2. if !seen(last.X-1,last.Y) && seen(last.X,last.Y-1) (use (last.X-1,last.Y))
Rule 3. if !seen(last.X,last.Y+1) (use (last.X, last.Y+1))
Rule 4. if !seen(last.X+1,last.Y) (use (last.X+1,last.Y))
*/

import (
	"flag"
	"fmt"
	"math"
	"os"
)

var input = flag.Int("input", 23, "Find steps to access this data square")
var partB = flag.Bool("partB", false, "Perform part B solution?")

type Cell struct {
	X, Y int
}

/*
  Create the next Cell for data.
  Rules: Always want to turn left, where possible, which means that we need to look at the last cell's position
  relative to our own because it will inform our direction ("two points is a line") since sometimes a "left turn"
  means decrementing X and sometimes incrementing X.
*/

func CreateCell(seen *map[Cell]bool, lastCell *Cell) *Cell {
	var ret *Cell
	// Rule 1
	if (!(*seen)[Cell{X: lastCell.X - 1, Y: lastCell.Y}]) && (!(*seen)[Cell{X: lastCell.X, Y: lastCell.Y - 1}]) {
		// Sub-check
		if !(*seen)[Cell{X: lastCell.X + 1, Y: lastCell.Y}] {
			// This is rule 4
			ret = &Cell{
				X: lastCell.X + 1,
				Y: lastCell.Y,
			}
		} else {
			// Rule 1 primary success
			ret = &Cell{
				X: lastCell.X,
				Y: lastCell.Y - 1,
			}

		}
	} else if (!(*seen)[Cell{X: lastCell.X - 1, Y: lastCell.Y}]) && ((*seen)[Cell{X: lastCell.X, Y: lastCell.Y - 1}]) {
		ret = &Cell{
			X: lastCell.X - 1,
			Y: lastCell.Y,
		}
	} else if (!(*seen)[Cell{X: lastCell.X, Y: lastCell.Y + 1}]) {
		ret = &Cell{
			X: lastCell.X,
			Y: lastCell.Y + 1,
		}
	} else if (!(*seen)[Cell{X: lastCell.X + 1, Y: lastCell.Y}]) {
		ret = &Cell{
			X: lastCell.X + 1,
			Y: lastCell.Y,
		}
	} else {
		// this should probably return error...
		ret = nil
	}

	return ret
}

func (c Cell) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

func (c Cell) GetNeighbours(seen *map[Cell]bool) []Cell {
	ret := make([]Cell, 0) // grow as needed
	for x := -1; x <= 1; x++ {
	y:
		for y := -1; y <= 1; y++ {
			t := Cell{X: c.X + x, Y: c.Y + y}
			// skip myself
			if c == t {
				continue y
			}
			if (*seen)[t] {
				ret = append(ret, t)
			}
		}
	}
	return ret
}

func main() {
	flag.Parse()
	seenLargest := false
	largestI := 0
	largestValue := 0
	candidateNumber := *input

	seenCells := make(map[Cell]bool)
	cells := make(map[Cell]int)

	// initialize the state so we start with 1 and 2 in place already.
	seenCells[Cell{X: 0, Y: 0}] = true // 1
	seenCells[Cell{X: 1, Y: 0}] = true // 2
	cells[Cell{X: 0, Y: 0}] = 1
	if *partB {
		cells[Cell{X: 1, Y: 0}] = 1 // Sum is 1
	} else {
		cells[Cell{X: 1, Y: 0}] = 2 // Second digit
	}

	lastCell := &Cell{X: 1, Y: 0}

	for i := 3; i <= candidateNumber; i++ {
		r := CreateCell(&seenCells, lastCell)
		if r == nil {
			fmt.Printf("Couldn't figure out a rule for %d (wtf?)!\n", i)
			os.Exit(1)
		}
		if *partB {
			sum := 0
			for _, neighbourCell := range r.GetNeighbours(&seenCells) {
				sum += cells[neighbourCell]
			}

			cells[*r] = sum
			if !seenLargest && sum > candidateNumber {
				fmt.Printf("(Iteration %d) %d is larger than input (which was %d)\n", i, sum, candidateNumber)
				seenLargest = true
				largestValue = sum
				largestI = i
			}
		} else {
			cells[*r] = i
		}
		seenCells[*r] = true
		lastCell = r

	}
	if *partB {
		fmt.Printf("On iteration %d we saw a sum of %d, which was larger than input of %d\n", largestI, largestValue, candidateNumber)
	} else {
		fmt.Printf("Distance from %s and %s: %f\n", lastCell, Cell{X: 0, Y: 0}, math.Abs(float64(lastCell.X))+math.Abs(float64(lastCell.Y)))
	}

}
