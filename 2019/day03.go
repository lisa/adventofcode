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

// where turns take place
type coordinate struct {
	x int
	y int
}

func (c *coordinate) distanceTo(o *coordinate) int {
	r := int(math.Abs(float64(c.x-o.x)) + math.Abs(float64(c.y-o.y)))
	return r
}

type lineSegment struct {
	xStart, xEnd int
	yStart, yEnd int
	vertical     bool
}

type grid struct {
	coordinates            []*coordinate
	segments               []*lineSegment
	maxX, minX, maxY, minY int
	curX, curY             int
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
func (g *grid) addCoordinate(c *coordinate) {
	if c.x > g.maxX {
		g.maxX = c.x
	}
	if c.x < g.minX {
		g.minX = c.x
	}
	if c.y > g.maxY {
		g.maxY = c.y
	}
	if c.y < g.minY {
		g.minY = c.y
	}

	segment := newSegment(g.curX, g.curY, c.x, c.y)

	g.coordinates = append(g.coordinates, c)
	g.segments = append(g.segments, segment)
	g.curX = c.x
	g.curY = c.y
}

func (g *grid) findOverlapsWith(o *grid) *[]coordinate {
	ret := make([]coordinate, 0)

	for _, gs := range g.segments {
		for _, os := range o.segments {
			if *debug {
				fmt.Printf("Checking line segment (%d,%d), (%d,%d) with (%d,%d), (%d,%d)\n",
					gs.xStart, gs.yStart, gs.xEnd, gs.yEnd,
					os.xStart, os.yStart, os.xEnd, os.yEnd)
			}
			if overlap := gs.overlapsWith(os); overlap != nil {
				if *debug {
					fmt.Printf("Overlap at (%d,%d)\n", overlap.x, overlap.y)
				}
				ret = append(ret, *overlap)
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

		for _, token := range strings.Split(line, ",") {
			moveAmount, err := strconv.Atoi(token[1:])
			if err != nil {
				fmt.Printf("Couldn't parse %s into a vector %s\n", token, err.Error())
				os.Exit(1)
			}
			if *debug {
				fmt.Printf("[%s] Wire %d at (%03d,%03d), Moving ", token, read, curX, curY)
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
				fmt.Printf(" to (%03d,%03d)\n", curX, curY)
			}

			grids[read].addCoordinate(&coordinate{x: curX, y: curY})
		}
		if *debug {
			fmt.Println()
		}
		read += 1
	}
	// end reading input lines
	overlaps := grids[0].findOverlapsWith(&grids[1])
	origin := coordinate{x: 0, y: 0}
	minDist := math.MaxInt64
	var minPoint coordinate
	for _, overlap := range *overlaps {
		if *debug{
		fmt.Printf("Overlap at (%d,%d)\n", overlap.x, overlap.y)}
		if d := overlap.distanceTo(&origin); d < minDist {
			minDist = d
			minPoint = overlap
		}
	}
	fmt.Printf("Closest overlap (%d,%d) = %d\n", minPoint.x, minPoint.y, minDist)
	os.Exit(0)
}
