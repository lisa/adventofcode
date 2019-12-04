package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	partB       = flag.Bool("partB", false, "Perform part B solution?")
	inputFile   = flag.String("inputFile", "inputs/day03a.txt", "Input File")
	inputString = flag.String("input", "", "Input string")
	debug       = flag.Bool("debug", false, "Debug?")
)

type partBResult struct {
	wire1Mag int
	wire2Mag int
	overlap  *coordinate
}

// where turns take place
type coordinate struct {
	x int
	y int
}

func (c *coordinate) distanceTo(o *coordinate) int {
	r := int(math.Abs(float64(c.x-o.x)) + math.Abs(float64(c.y-o.y)))
	if *debug {
		fmt.Printf("Distance from (%d,%d) to (%d,%d) is %d\n", c.x, c.y, o.x, o.y, r)
	}
	return r
}

type lineSegment struct {
	xStart, xEnd   int
	yStart, yEnd   int
	totalMagnitude int // total magnitude (distance) to (xStart,yStart)
	magnitude      int

	vertical bool
}

type grid struct {
	coordinates []*coordinate
	segments    []*lineSegment
	curX, curY  int
}

func newSegment(x1, y1, x2, y2 int) *lineSegment {

	seg := lineSegment{
		xStart: x1,
		xEnd:   x2,
		yStart: y1,
		yEnd:   y2,
	}

	if x1 == x2 {
		seg.vertical = true
	}

	return &seg
}

// does s include c?
func (s *lineSegment) includesCoordinate(c *coordinate) bool {
	includedX := false
	includedY := false
	if s.xStart > s.xEnd {
		// right to left segment
		includedX = c.x <= s.xStart && c.x >= s.xEnd
	} else if s.xStart < s.xEnd {
		includedX = c.x >= s.xStart && c.x <= s.xEnd
	} else {
		includedX = true
	}
	if s.yStart > s.yEnd {
		includedY = c.y <= s.yStart && c.y >= s.yEnd
	} else if s.yStart < s.yEnd {
		includedY = c.y >= s.yStart && c.y <= s.yEnd
	} else {
		includedY = true
	}
	return includedX && includedY
}

// do I overlap with another lineSegment?
// return the overlapping coordinate, if any. nil otherwise
func (s *lineSegment) overlapsWith(o *lineSegment) *coordinate {

	// If both are in the same orientation, we have division by zero.
	if s.vertical && o.vertical || !s.vertical && !o.vertical {
		return nil
	}
	// Create a candidate for overlapping
	// https://math.stackexchange.com/questions/375083/given-coordinates-of-beginning-and-end-of-two-intersecting-line-segments-how-do
	x := -1 * ((s.xStart-s.xEnd)*(o.xStart*o.yEnd-o.xEnd*o.yStart) - (o.xEnd-o.xStart)*(s.xEnd*s.yStart-s.xStart*s.yEnd)) / ((o.yStart-o.yEnd)*(s.xStart-s.xEnd) - (o.xEnd-o.xStart)*(s.yEnd-s.yStart))
	y := -1 * (o.xStart*o.yEnd*s.yStart - o.xStart*o.yEnd*s.yEnd - o.xEnd*o.yStart*s.yStart + o.xEnd*o.yStart*s.yEnd - o.yStart*s.xStart*s.yEnd + o.yStart*s.xEnd*s.yStart + o.yEnd*s.xStart*s.yEnd - o.yEnd*s.xEnd*s.yStart) / (-1*o.xStart*s.yStart + o.xStart*s.yEnd + o.xEnd*s.yStart - o.xEnd*s.yEnd + o.yStart*s.xStart - o.yStart*s.xEnd - o.yEnd*s.xStart + o.yEnd*s.xEnd)
	// we never match on (0,0)
	if x == 0 && y == 0 {
		return nil
	}
	test := coordinate{
		x: x,
		y: y,
	}
	if s.includesCoordinate(&test) && o.includesCoordinate(&test) {
		return &test
	}
	return nil
}

// executes a turn by creating a new line segment to represent that turn
// c represents the "current" position
func (g *grid) addCoordinate(c *coordinate, magnitude, totalMagnitude int) {
	segment := newSegment(g.curX, g.curY, c.x, c.y)
	segment.magnitude = magnitude
	segment.totalMagnitude = totalMagnitude

	g.coordinates = append(g.coordinates, c)
	g.segments = append(g.segments, segment)
	g.curX = c.x
	g.curY = c.y
}

func (g *grid) findOverlapsWith(o *grid) *[]partBResult {
	ret := make([]partBResult, 0)
	for _, gs := range g.segments {
		for _, os := range o.segments {
			if *debug {
				fmt.Printf("Checking line segment (%d,%d), (%d,%d) (t=%d; m=%d) with (%d,%d), (%d,%d) (t=%d; m=%d).\n",
					gs.xStart, gs.yStart, gs.xEnd, gs.yEnd, gs.totalMagnitude, gs.magnitude,
					os.xStart, os.yStart, os.xEnd, os.yEnd, os.totalMagnitude, os.magnitude)
			}
			if overlap := gs.overlapsWith(os); overlap != nil {
				if *debug {
					fmt.Printf("Overlap at (%d,%d)\n", overlap.x, overlap.y)
					fmt.Printf("  Line 1 total mag: %d, Line 2 total mag: %d\n", gs.totalMagnitude, os.totalMagnitude)
					fmt.Printf("    Partial for line 1: %d, line 2: %d\n",
						overlap.distanceTo(&coordinate{x: gs.xStart, y: gs.yStart}),
						overlap.distanceTo(&coordinate{x: os.xStart, y: os.yStart}))
				}
				// These need to be adjusted because it includes the full magnitude of this
				// segment, and it ought to only include the amount to overlap.
				ret = append(ret, partBResult{
					wire1Mag: gs.totalMagnitude + overlap.distanceTo(&coordinate{x: gs.xStart, y: gs.yStart}),
					wire2Mag: os.totalMagnitude + overlap.distanceTo(&coordinate{x: os.xStart, y: os.yStart}),
					overlap:  overlap,
				})
			}
		}
		if *debug {
			fmt.Println()
		}
	}

	return &ret
}

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't open %s: %v\n", *inputFile, err)
		os.Exit(1)
	}

	lineReader := bufio.NewScanner(input)
	grids := make([]grid, 0)
	read := 0
	for lineReader.Scan() {
		line := lineReader.Text()
		grids = append(grids, grid{})

		curX := 0
		curY := 0
		totalMagnitude := 0

		for _, token := range strings.Split(line, ",") {
			moveAmount, err := strconv.Atoi(token[1:])
			if err != nil {
				fmt.Printf("Couldn't parse %s into a vector %s\n", token, err.Error())
				os.Exit(1)
			}
			if *debug {
				fmt.Printf("[%s] Wire %d at (%d,%d), Moving ", token, read, curX, curY)
			}
			switch token[0:1] {
			case "U":
				//up
				if *debug {
					fmt.Printf("up %d", moveAmount)
				}
				curY += moveAmount
			case "R":
				//right
				if *debug {
					fmt.Printf("right %d", moveAmount)
				}
				curX += moveAmount
			case "D":
				//down
				if *debug {
					fmt.Printf("down %d", moveAmount)
				}
				curY -= moveAmount
			case "L":
				//left
				if *debug {
					fmt.Printf("left %d", moveAmount)
				}
				curX -= moveAmount
			}
			if *debug {
				fmt.Printf(" to (%d,%d)\n", curX, curY)
			}

			grids[read].addCoordinate(&coordinate{x: curX, y: curY}, moveAmount, totalMagnitude)
			totalMagnitude += moveAmount
		}
		if *debug {
			fmt.Println()
		}
		read += 1
	}
	// end reading input lines
	partBResults := grids[0].findOverlapsWith(&grids[1])
	origin := coordinate{x: 0, y: 0}
	minDist := math.MaxInt64
	minMag := math.MaxInt64
	var minDistPoint coordinate
	var minMagPoint coordinate
	for _, bRes := range *partBResults {
		if *debug {
			fmt.Printf("Overlap at (%d,%d)\n", bRes.overlap.x, bRes.overlap.y)
		}
		d := bRes.overlap.distanceTo(&origin)

		if d < minDist {
			minDist = d
			minDistPoint = *bRes.overlap
		}
		if *debug {
			fmt.Printf("(minMag: %d) Wire 1 Mag: %d, wire 2 mag: %d. Sum = %d\n", minMag, bRes.wire1Mag, bRes.wire2Mag, bRes.wire1Mag+bRes.wire2Mag)
		}
		if bRes.wire1Mag+bRes.wire2Mag < minMag {
			minMag = bRes.wire1Mag + bRes.wire2Mag
			minMagPoint = *bRes.overlap
		}

	}
	if !*partB {
		fmt.Printf("Closest overlap (%d,%d) = %d\n", minDistPoint.x, minDistPoint.y, minDist)
	} else {
		fmt.Printf("Closest overlap with lowest magnitude is (%d,%d) = %d\n", minMagPoint.x, minMagPoint.y, minMag)
	}
	os.Exit(0)
}
