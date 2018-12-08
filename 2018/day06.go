package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
)

var (
	partB     = flag.Bool("partB", false, "Perform part B solution?")
	inputFile = flag.String("input", "inputs/day06.txt", "Input")
	debug     = flag.Bool("debug", false, "Debug?")

	coords = regexp.MustCompile(`^(\d+), (\d+)$`)
)

// part A settings
var (
	minXPadding = 10
	maxXPadding = 10
	minYPadding = 10
	maxYPadding = 10
)

// Point is an (x,y) on the plane. (0,0) is top left
type Point struct {
	X, Y   int
	Claims int // how many coordinates does this Point claim?
}

// Plane represents the known points, and the min/max coord values
// north means Y values get smaller, south means Y values get larger
// west means X values get smaller, east means X values get larger
// (0,0) is at the top-left of the plane.
type Plane struct {
	KnownPoints []*Point
	MaxX        int // East-most
	MaxY        int // South-most
	MinY        int // North-most
	MinX        int // West-most

	checkMap            map[int]map[int][]CoordinateClaim // [x][y][]CoordinateClaim
	checkMapInitialized bool                              // have we done the initialization yet?
}

// NewPoint for the new Point
func NewPoint(x, y int) *Point {
	return &Point{
		X:      x,
		Y:      y,
		Claims: 0, // i always claim myself
	}
}

// IsEqual - am i equal to another?
func (p *Point) IsEqual(other *Point) bool {
	return p.X == other.X && p.Y == other.Y
}

// Distance - from myself to another Point
func (p *Point) Distance(other *Point) int {
	return p.DistanceToXY(other.X, other.Y)
}

// DistanceToXY - distance to an (x,y) coordinate
func (p *Point) DistanceToXY(x, y int) int {
	r := int(math.Abs(float64(p.X-x)) + math.Abs(float64(p.Y-y)))
	return r
}
func (p *Point) String() string {
	return fmt.Sprintf("(%d,%d; %d)", p.X, p.Y, p.Claims)
}

// InitializeCheckMap - Create the distance map from all points in the Plane to the given KnownPoints
func (p *Plane) InitializeCheckMap() {
	fmt.Printf("Initializing the plane\n")
	if p.checkMapInitialized {
		return
	}
	// init map. This computes the distances from every value in the (x,y) plane to every known point
	p.checkMap = make(map[int]map[int][]CoordinateClaim)
	for x := p.MinX - minXPadding; x <= p.MaxX+maxXPadding; x++ {
		p.checkMap[x] = make(map[int][]CoordinateClaim)
		for y := p.MinY - minYPadding; y <= p.MaxY+maxYPadding; y++ {
			p.checkMap[x][y] = make([]CoordinateClaim, len(p.KnownPoints))
			for i, kp := range p.KnownPoints {
				p.checkMap[x][y][i] = CoordinateClaim{
					P:               kp,
					DistanceToPoint: kp.Distance(NewPoint(x, y)),
				}
			}
			// now sort for distance
			sort.Slice(p.checkMap[x][y], func(i, j int) bool { return p.checkMap[x][y][i].DistanceToPoint < p.checkMap[x][y][j].DistanceToPoint })
		}
	}
	p.checkMapInitialized = true
}

// IsPointInfinite - Is the point in question infinite?
// A point `t` is infinite if, for all points `c` along the outside perimeter of
// the bounding rectangle, `t` is the only closest point of all known points `kp`.
// To support this, we need a way to associate all `c` to `t`, and the distance for all `kp`s to `c`.
// this would look like:
// *Point[x][y]distance - from this we can sort sort the second map
// based on its distance value and check 0th and 1st element for a distance match
func (p *Plane) IsPointInfinite(point *Point) bool {
	ret := false
	p.InitializeCheckMap()

	// For the given +point+, do a circuit around the perimeter and check to see if
	// it is the only closest point. This will be evident because a each point we
	// will sort by distance, and the shortest[0] and shortest[1]
	for y := p.MinY - minYPadding; y <= p.MaxY+maxYPadding; y++ {
		for _, x := range []int{p.MinX - minXPadding, p.MaxX + maxXPadding} {
			// For this (x,y) coordinate we want to know if:
			// What are the two closest Points (p1, p2) aka (least and secondLeast)
			// (1) Is p1 or p2 the +point+ we care about?
			// (2) Is p1's distance to (x,y) equal to p2's distance to (x,y)?
			// (3a) If p1 is our point and (2) is true then p1 is NOT infinite
			// (3b) if p1 is our point and (2) is false then p1 is infinite
			// (3c) if p2 is our point and (2) is true then p2 is NOT infinite
			// (3d) if p2 is our point and (2) is true then p2 is NOT infinite
			least := p.checkMap[x][y][0]
			secondLeast := p.checkMap[x][y][1]
			if *debug {
				fmt.Printf("1: (%d,%d)\tp=%s\tl=%s;d=%d\ts=%s;d=%d\n", x, y, point, least.P, least.DistanceToPoint, secondLeast.P, secondLeast.DistanceToPoint)
			}
			if secondLeast.P.IsEqual(point) || !least.P.IsEqual(point) {
				// If secondLeast is ever the +point+ then it can never be infinite because if
				// it's second and the distance is the same, that means noninfinite. if the
				// distances are different, then least's distance < secondLeast's distance
				// 3c and 3d
				continue
			}
			// least is our +point+ and p1

			// is least or secondLeast our point?
			if least.DistanceToPoint == secondLeast.DistanceToPoint {
				if *debug {
					fmt.Printf(" * ld == sd (%d == %d)\n", least.DistanceToPoint, secondLeast.DistanceToPoint)
				}
				// 3a
				continue
			} else {
				// 3b
				if *debug {
					fmt.Printf(" * %s infinite because: least=%s (d=%d) is the closest point to (%d,%d) because secondLeast=%s (d=%d)\n",
						point, least.P, least.DistanceToPoint, x, y, secondLeast.P, secondLeast.DistanceToPoint)
				}
				return true
			}

		}
	}

	for x := p.MinX + minXPadding; x <= p.MaxX+maxXPadding; x++ {
		for _, y := range []int{p.MinY - minYPadding, p.MaxY + maxYPadding} {
			// For this (x,y) coordinate we want to know if:
			// What are the two closest Points (p1, p2) aka (least and secondLeast)
			// (1) Is p1 or p2 the +point+ we care about?
			// (2) Is p1's distance to (x,y) equal to p2's distance to (x,y)?
			// (3a) If p1 is our point and (2) is true then p1 is NOT infinite
			// (3b) if p1 is our point and (2) is false then p1 is infinite
			// (3c) if p2 is our point and (2) is true then p2 is NOT infinite
			// (3d) if p2 is our point and (2) is true then p2 is NOT infinite
			least := p.checkMap[x][y][0]
			secondLeast := p.checkMap[x][y][1]
			if *debug {
				fmt.Printf("2: (%d,%d)\tp=%s\tl=%s;d=%d\ts=%s;d=%d\n", x, y, point, least.P, least.DistanceToPoint, secondLeast.P, secondLeast.DistanceToPoint)
			}
			if secondLeast.P.IsEqual(point) || !least.P.IsEqual(point) {
				// If secondLeast is ever the +point+ then it can never be infinite because if
				// it's second and the distance is the same, that means noninfinite. if the
				// distances are different, then least's distance < secondLeast's distance
				// 3c and 3d
				continue
			}
			// least is our +point+ and p1

			// is least or secondLeast our point?
			if least.DistanceToPoint == secondLeast.DistanceToPoint {
				// 3a
				if *debug {
					fmt.Printf(" * ld == sd (%d == %d)\n", least.DistanceToPoint, secondLeast.DistanceToPoint)
				}
				continue
			} else {
				// 3b
				if *debug {
					fmt.Printf(" * %s infinite because: least=%s (d=%d) is the closest point to (%d,%d) because secondLeast=%s (d=%d)\n",
						point, least.P, least.DistanceToPoint, x, y, secondLeast.P, secondLeast.DistanceToPoint)
				}
				return true
			}

		}
	}

	return ret
}

// NewPlane - make a new plane to work with
func NewPlane() *Plane {
	return &Plane{
		KnownPoints:         make([]*Point, 0),
		MaxX:                -1 * math.MaxInt16,
		MaxY:                -1 * math.MaxInt16,
		MinX:                math.MaxInt16,
		MinY:                math.MaxInt16,
		checkMapInitialized: false,
	}
}

// AddPoint - add a Point to the Known Points (and do some bounds bookkeeping)
func (p *Plane) AddPoint(t *Point) {
	p.KnownPoints = append(p.KnownPoints, t)
	if p.MaxX < t.X {
		p.MaxX = t.X
	}
	if p.MinX > t.X {
		p.MinX = t.X
	}
	if p.MaxY < t.Y {
		p.MaxY = t.Y
	}
	if p.MinY > t.Y {
		p.MinY = t.Y
	}
	return
}

type CoordinateClaim struct {
	P               *Point
	DistanceToPoint int
}

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	if *partB {
		minXPadding = 10000
		maxXPadding = 10000
		minYPadding = 10000
		maxYPadding = 10000
	}

	input, err := os.Open(*inputFile)
	if err != nil {
		errorIf(fmt.Sprintf("Couldn't open %s: %v", *inputFile, err), err)
		os.Exit(1)
	}
	defer input.Close()
	plane := NewPlane()
	lineReader := bufio.NewScanner(input)

	for lineReader.Scan() {
		matches := coords.FindAllStringSubmatch(lineReader.Text(), -1)
		x, err := strconv.Atoi(string(matches[0][1]))
		errorIf("Couldn't parse X coordinate", err)
		y, err := strconv.Atoi(string(matches[0][2]))
		errorIf("Couldn't parse Y coordinate", err)

		plane.AddPoint(NewPoint(x, y))
	}
	if *debug {
		fmt.Printf("Bounding rectangle of the Plane (padded): (%d,%d), (%d,%d), (%d,%d), (%d,%d)\n",
			plane.MinX-minXPadding, plane.MinY-minYPadding, plane.MaxX+maxXPadding, plane.MinY-minYPadding, plane.MaxX+maxXPadding, plane.MaxY+maxYPadding, plane.MinX-minXPadding, plane.MaxY+maxYPadding)
	}
	// the checkmap tells us the distance from every known point to the outer
	// reaches of the bounding rectangle. thus, we can go through it and eliminate
	// all the points which have an infinite area, and, with the remainder, find
	// their area
	finitePoints := make([]*Point, 0)
	if !*partB {
		plane.InitializeCheckMap()

		for _, candidatePoint := range plane.KnownPoints {
			if plane.IsPointInfinite(candidatePoint) {
				if *debug {
					fmt.Printf("Skipping %s because it's infinite\n", candidatePoint)
				}
				continue
			} else {
				finitePoints = append(finitePoints, candidatePoint)
				if *debug {
					fmt.Printf("%s is finite, using it\n", candidatePoint)
				}
			}
		}
	}

	if !*partB {
		var least, secondLeast CoordinateClaim
		for y := plane.MinY - minYPadding; y <= plane.MaxY+maxYPadding; y++ {
			for x := plane.MinX - minXPadding; x <= plane.MaxX+maxXPadding; x++ {
				least = plane.checkMap[x][y][0]
				secondLeast = plane.checkMap[x][y][1]

				for _, fp := range finitePoints {
					// finite point gets a score if it's the closest to (x,y) without there being a tie
					// if least is the same as the finite point we're looking at here, then +1,
					// unless secondLeast has the same distance
					if least.DistanceToPoint == secondLeast.DistanceToPoint {
						continue
					}
					// there's not a tie, so, are we the leader?
					if least.P.IsEqual(fp) {
						fp.Claims++
					}

				}
			}
		}

		sort.Slice(finitePoints, func(i, j int) bool { return finitePoints[i].Claims > finitePoints[j].Claims })
		if *debug {
			fmt.Printf("Total Scores: %s\n", finitePoints)
		}
		fmt.Printf("High Score: %s\n", finitePoints[0])
	} else {
		//part B
		// What is the size of the region containing all locations which have a total
		// distance to all input Points of less than 10000? Note, you can't use the
		// checkMap or it will run you out of memory!
		//
		// Make the boundary padding 10000 on each side, then go through and select all (x,y) points with a sum <10000
		var regionPoints, sum int
		regionPoints = 0
		for y := plane.MinY - minYPadding; y <= plane.MaxY+maxYPadding; y++ {
			for x := plane.MinX - minXPadding; x <= plane.MaxX+maxXPadding; x++ {
				sum = 0
				for _, kp := range plane.KnownPoints {
					sum += kp.DistanceToXY(x, y)
				}
				if sum < 10000 {
					if *debug {
						fmt.Printf("(%d,%d) k %d.\n", x, y, sum)
					}
					regionPoints++
				} else {
					if *debug {
						fmt.Printf("(%d,%d) r %d.\n", x, y, sum)
					}
				}
			} // end x
		}
		fmt.Printf("Points in the region: %d\n", regionPoints)
	}

}
