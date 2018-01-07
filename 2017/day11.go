package main

/*
Day 11 Part A

For a given comma-separated list of directions [n,ne,se,s,sw,nw] determine how
far net movement is from the origin. For example:

ne,ne,ne is 3 steps away
ne,ne,sw,sw is 0 away (back at the start)
ne,ne,s,s is 2 steps away (se,se)
se,sw,se,sw,sw is 3 steps away (s,s,sw)

Hexagon:

  \ n  /
nw +--+ ne
  /    \
-+      +-
  \    /
sw +--+ se
  / s  \

Possible moves with offsets:
N : X+1, Y+0
NE: X+1, Y-1
SE: X+0, Y-1
S : X-1, Y+0
SW: X-1, Y+1
NW: X+0, Y+1
*/

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
)

var input = flag.String("input", "ne,ne,ne", "Puzzle Input")
var debug = flag.Bool("debug", false, "Debug")

// Hex Stuff
type Hex struct {
	X, Y int
}

func (h *Hex) String() string {
	return fmt.Sprintf("(%d,%d)", h.X, h.Y)
}

// Move in the direction of n/ne/se/s/sw/nw.
// If that direction doesn't exist, create it
// Major credit to http://keekerdc.com/2011/03/hexagon-grids-coordinate-systems-and-distance-calculations/
// for helping me understand the hexagon layout vis-a-vis how X and Y change.
func (h *Hex) Move(dir string) (*Hex, error) {
	var offsetX, offsetY int
	switch dir {
	case "n":
		offsetX, offsetY = 1, 0
	case "ne":
		offsetX, offsetY = 1, -1
	case "se":
		offsetX, offsetY = 0, -1
	case "s":
		offsetX, offsetY = -1, 0
	case "sw":
		offsetX, offsetY = -1, 1
	case "nw":
		offsetX, offsetY = 0, 1
	default:
		return nil, errors.New(fmt.Sprintf("Unexpected direction %s, can't move there", dir))
	} // end switch
	return NewHex(h.X+offsetX, h.Y+offsetY), nil
}

func NewHex(x, y int) *Hex {
	return &Hex{
		X: x,
		Y: y,
	}
}

func main() {
	flag.Parse()

	currentHex := NewHex(0, 0)
	if *debug {
		fmt.Printf("Initial hex %s\n", currentHex)
	}
	var err error

	if *debug {
		fmt.Printf("currentHex: %s\n", currentHex)
	}
	for _, direction := range strings.Split(*input, ",") {
		currentHex, err = currentHex.Move(direction)
		if err != nil {
			fmt.Printf("Got error: %s\n", err)
			os.Exit(1)
		} else {
			if *debug {
				fmt.Printf("currentHex: %s (moved %s)\n", currentHex, direction)
			}
		}
	}
	if *debug {
		fmt.Println()
	}
	var moves int
	fmt.Printf("After making the moves the location is: %+v\n", currentHex)
	if currentHex.X == 0 {
		moves = int(math.Abs(float64(currentHex.Y)))
	} else if currentHex.X < 0 {
		if currentHex.Y == 0 {
			moves = int(math.Abs(float64(currentHex.X)))
		} else if currentHex.Y > 0 {
			moves = int(math.Abs(float64(currentHex.X))) + (currentHex.Y - int(math.Abs(float64(currentHex.X))))
		} else {
			moves = int(math.Abs(float64(currentHex.X + currentHex.Y)))
		}
	} else {
		// X > 0
		if currentHex.Y == 0 {
			moves = int(math.Abs(float64(currentHex.X)))
		} else if currentHex.Y > 0 {
			moves = currentHex.Y + (int(math.Abs(float64(currentHex.X))) - currentHex.Y)
		} else {
			moves = int(math.Abs(float64(currentHex.Y))) + (int(math.Abs(float64(currentHex.X))) - int(math.Abs(float64(currentHex.Y))))
		}
	}

	fmt.Printf("Moves: %d\n", moves)
}
