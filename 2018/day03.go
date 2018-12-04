package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	partB     = flag.Bool("partB", false, "Perform part B solution?")
	inputFile = flag.String("input", "inputs/day03.txt", "Input")
	debug     = flag.Bool("debug", false, "Debug?")
	debug2    = flag.Bool("debug2", false, "Second debug level?")

	picker = regexp.MustCompile(`^#(\d.*) @ (\d{1,}),(\d{1,}): (\d{1,})x(\d{1,})$`)
)

// Point - an (x,y) point on a plane.
type Point struct {
	X int
	Y int
}

// NewPoint - Create a pointer to a Point
func NewPoint(x, y int) *Point {
	return &Point{
		X: x,
		Y: y,
	}
}

// Claim - #123 @ 3,2: 5x4
// Note: LeftX is inches from the left _BEFORE_ our Claim
//     : TopY is the inches from the top _BEFORE_ our Claim
// Thus, we really exist within these bounds:
// (LeftX+1,TopY+1), (LeftX+1+Width,TopY+1), (LeftX+1,TopY+1+Height), (LeftX+1+Width,TopY+1+Height)
// Or: X plane [LeftX+1,LeftX+1+Width] & Y Plane [TopY+1,TopY+1+Height]
type Claim struct {
	ID       int      // #123
	Height   int      // 4
	Width    int      // 5
	LeftX    int      // 3
	TopY     int      // 2
	Points   []*Point // (X,Y) coords this Claim has
	Overlaps int      // how many other points does this overlap with?
}

// NewClaim - make a new claim on the fabric
func newClaim(id, leftx, topy, width, height int) *Claim {
	return &Claim{
		ID:       id,
		LeftX:    leftx,
		TopY:     topy,
		Width:    width,
		Height:   height,
		Overlaps: 0,
		Points:   make([](*Point), 0),
	}
}

// ParseClaim - parse the Claim from the raw input
func ParseClaim(rawLine string) *Claim {
	match := picker.FindAllStringSubmatch(rawLine, -1)
	id, err := strconv.Atoi(match[0][1])
	if err != nil {
		fmt.Printf("Couldn't parse from %s\n", rawLine)
		os.Exit(1)
	}
	leftX, err := strconv.Atoi(match[0][2])
	if err != nil {
		fmt.Printf("Couldn't parse from %s\n", rawLine)
		os.Exit(1)
	}
	topY, err := strconv.Atoi(match[0][3])
	if err != nil {
		fmt.Printf("Couldn't parse from %s\n", rawLine)
		os.Exit(1)
	}
	width, err := strconv.Atoi(match[0][4])
	if err != nil {
		fmt.Printf("Couldn't parse from %s\n", rawLine)
		os.Exit(1)
	}
	height, err := strconv.Atoi(match[0][5])
	if err != nil {
		fmt.Printf("Couldn't parse from %s\n", rawLine)
		os.Exit(1)
	}
	claim := newClaim(id, leftX, topY, width, height)
	for x := claim.LeftX + 1; x <= claim.LeftX+claim.Width; x++ {
		for y := claim.TopY + 1; y <= claim.TopY+claim.Height; y++ {
			claim.Points = append(claim.Points, NewPoint(x, y))
		}
	}
	return claim
}

// FindOverlapSize - Returns the total area of overlap between my Claim and the other Claim.
// If there is no overlap the return will be 0.
// Claims overlap if any region of the other Claim is within the bounds of our own Claim.
// FIXME: This is insufficient because (2,5) doesn't overlap (2,7). This check will count it as an overlap.
// Brute force method is to check each cell against each one another.
func (c *Claim) FindOverlapSize(other *Claim) []*Point {

	overlap := make([]*Point, 0)

	// Go through all our points, compare to their points.
	for _, p := range c.Points {
		//p is a *Point
		for _, o := range other.Points {
			if p.X != o.X || p.Y != o.Y {
				continue
			}
			if *debug2 {
				fmt.Printf("FindOverlapSize: My Point: (%d,%d) -> Their Point (%d,%d)\n", p.X, p.Y, o.X, o.Y)
			}

			if p.X == o.X && p.Y == o.Y {
				c.Overlaps++
				other.Overlaps++
				overlap = append(overlap, p)
			}
		}
	}

	if *debug2 {
		fmt.Printf("Overlaps %d\n", overlap)
	}
	return overlap

}

func main() {
	flag.Parse()
	fmt.Printf("Day 3\n")

	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't open %s: %v\n", *inputFile, err)
		os.Exit(1)
	}
	defer input.Close()
	lineReader := bufio.NewScanner(input)
	inputRows := make([]string, 0)
	for lineReader.Scan() {
		line := lineReader.Text()
		inputRows = append(inputRows, strings.ToLower(line))
	}

	allClaims := make([]*Claim, len(inputRows))
	for i, row := range inputRows {
		allClaims[i] = ParseClaim(row)
		if *debug2 {
			fmt.Printf("Added a Claim: index=%d, %+v\n", i, allClaims[i])
		}
	}
	if *debug2 {
		fmt.Printf("\n")
	}
	overlap := make([]*Point, 0)
	// Loop through all the Claims and compare them to one another.
	// We can be smart about this. If we've compared ID 1 to 2 we don't need to
	// compare 2 to 1. The way this looks is, if the right side is < the left side,
	// we skip.
	for left := 0; left < len(allClaims)-1; left++ {
		if *debug2 {
			fmt.Printf("Left: %+v\n", allClaims[left])
		}
		for right := left + 1; right < len(allClaims); right++ {
			if right < left {
				continue
			}
			if *debug2 {
				fmt.Printf("left=%d, right=%d\n", left, right)
			}
			if *debug2 {
				fmt.Printf("  Right: %+v\n", allClaims[right])
			}
			overlap = append(overlap, allClaims[left].FindOverlapSize(allClaims[right])...)
		}
	}
	if !*partB {

		uniqueOverlaps := make(map[Point]bool)
		// Deduplicate.
		for _, overlappingPoint := range overlap {
			uniqueOverlaps[*NewPoint(overlappingPoint.X, overlappingPoint.Y)] = true
		}
		if *debug {
			fmt.Printf("All overlapping points: %+v\n", uniqueOverlaps)
		}
		fmt.Printf("Total overlap: %d\n", len(uniqueOverlaps))
	} else {
		// Find the sole Claim without any overlaps (so special!)
		for _, c := range allClaims {
			if c.Overlaps == 0 {
				// god, i hope there's only one here
				fmt.Printf("The Special Claim is ID %d\n", c.ID)
			} else {
				if *debug {
					fmt.Printf("Claim %d has %d overlaps (rip)\n", c.ID, c.Overlaps)
				}
			}
		}
	}
}
