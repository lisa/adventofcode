package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"

	. "github.com/logrusorgru/aurora"
)

var (
	inputFile = flag.String("input", "inputs/day13.txt", "Input data")
	partB     = flag.Bool("partB", false, "Perform part B solution?")
	debug     = flag.Bool("debug", false, "debug")
	debug2    = flag.Bool("debug2", false, "Force keypress to advance ticks")
)

// SegmentType - type of segment (track)
type SegmentType int8

// Direction of travel for the Shuttles (carts)
type Direction int8

// Turn is the last movement choice made at an intersection
type Turn int8

// far left: x=0, top most: y=0, thus:
// east  -> x+1
// west  -> x-1
// north -> y-1
// south -> y+1
const (
	EastWest     SegmentType = iota // -
	NorthSouth                      // |
	TopLeftDiag                     // \ N to E or S to W
	TopRightDiag                    // / S to E or N to W
	Intersection                    // +

	// Direction of Movement
	East Direction = iota
	South
	West
	North

	Left Turn = iota
	Straight
	Right
)

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}
func segTypeToString(t SegmentType) string {
	switch t {
	case EastWest:
		return "East-West"
	case NorthSouth:
		return "North-South"
	case TopLeftDiag:
		return "Top-Left Diagonal"
	case TopRightDiag:
		return "Top-Right Diagonal"
	case Intersection:
		return "Intersection"
	}
	return "Unknown Segment Type"
}

func turnToString(t Turn) string {
	switch t {
	case Left:
		return "Left"
	case Straight:
		return "Straight"
	case Right:
		return "Right"
	case 0:
		return "None made yet"
	}
	return "Unknown turn type"
}

type Segment struct {
	X, Y int
	Type SegmentType

	North, South, West, East *Segment // possible connections
}

type Shuttle struct {
	X, Y              int
	CurrentSegment    *Segment
	DirectionOfTravel Direction
	LastTurn          Turn
}

// TurnLeft - If I would turn left, what would the segment be?
func (s *Shuttle) TurnLeft() Direction {
	switch s.DirectionOfTravel {
	case East:
		return North
	case South:
		return East
	case West:
		return South
	case North:
		return West
	}
	// should never get here!
	fmt.Printf("Shuttle at (%d,%d) going in direction %v tried to make an illegal Left Turn! Segment dump %+v\n", s.X, s.Y, s.DirectionOfTravel, s.CurrentSegment)
	os.Exit(1)
	return 0
}

// MoveStraight - if I would move straight on, what would the segment be?
func (s *Shuttle) MoveStraight() Direction {
	switch s.DirectionOfTravel {
	case East:
		return East
	case South:
		return South
	case West:
		return West
	case North:
		return North
	}
	// should never get here!
	fmt.Printf("Shuttle at (%d,%d) going in direction %v tried to move straight! Segment dump %+v\n", s.X, s.Y, s.DirectionOfTravel, s.CurrentSegment)
	os.Exit(1)
	return 0
}

// TurnRight - if I would turn Right, what would the next segment be?
func (s *Shuttle) TurnRight() Direction {
	switch s.DirectionOfTravel {
	case East:
		return South
	case South:
		return West
	case West:
		return North
	case North:
		return East
	}
	// should never get here!
	fmt.Printf("Shuttle at (%d,%d) going in direction %v tried to make an illegal Right Turn! Segment dump %+v\n", s.X, s.Y, s.DirectionOfTravel, s.CurrentSegment)
	os.Exit(1)
	return 0
}

type Field struct {
	Segments   map[int]map[int]*Segment // (X,Y) -> *Segment
	Shuttles   []*Shuttle               // all the shuttles in progress
	MaxX, MaxY int
}

func NewField() *Field {
	return &Field{
		Segments: make(map[int]map[int]*Segment),
		Shuttles: make([]*Shuttle, 0),
		MaxX:     0,
		MaxY:     0,
	}
}

func (f *Field) PrintField() {
	for y := 0; y <= f.MaxY; y++ {
		for x := 0; x <= f.MaxX; x++ {
			shuttle, err := f.GetShuttleByXY(x, y)
			if err != nil {
				fmt.Printf("%s", Inverse("X"))
				continue
			} else if shuttle != nil {
				switch shuttle.DirectionOfTravel {
				case East:
					fmt.Printf("%s", Cyan(">").Bold())
				case West:
					fmt.Printf("%s", Green("<").Bold())
				case South:
					fmt.Printf("%s", Red("v").Bold())
				case North:
					fmt.Printf("%s", Magenta("^").Bold())
				}
				continue
			}
			if seg := f.GetSegmentByXY(x, y); seg != nil {
				switch seg.Type {
				case EastWest:
					fmt.Printf("-")
				case NorthSouth:
					fmt.Printf("|")
				case TopLeftDiag:
					fmt.Printf("\\")
				case TopRightDiag:
					fmt.Printf("/")
				case Intersection:
					fmt.Printf("+")
				}
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
}

// GetShuttleByXY - Get a shuttle by a given (x,y). If there's a collision, return an error
func (f *Field) GetShuttleByXY(x, y int) (*Shuttle, error) {
	seen := false
	var ret *Shuttle
	for _, shuttle := range f.Shuttles {
		if shuttle.X == x && shuttle.Y == y {
			if seen {
				return nil, errors.New("collision")
			} else {
				seen = true
				ret = shuttle
			}
		}
	}
	return ret, nil
}

// GetShuttlesByY - Get a shuttle by the Y value, in order of X:
func (f *Field) GetShuttlesByY(y int) []*Shuttle {
	ret := make([]*Shuttle, 0)

	for _, shuttle := range f.Shuttles {
		if shuttle.Y == y {
			ret = append(ret, shuttle)
		}
	}

	// sort to ensure left-to-right order based on X.
	sort.Slice(ret, func(i, j int) bool { return f.Shuttles[i].X < f.Shuttles[j].X })
	return ret
}

// AddSegment - adds a Segment into the Field. Expects its x,y position and segment type
func (f *Field) AddSegment(x, y int, segType SegmentType) {
	// Initialize field[x] if needed
	if *debug {
		fmt.Printf("Adding segment at (%d,%d) with type %s\n", x, y, segTypeToString(segType))
	}
	if _, ok := f.Segments[x]; !ok {
		f.Segments[x] = make(map[int]*Segment)
	}
	newSegment := &Segment{
		X:    x,
		Y:    y,
		Type: segType,
	}
	f.Segments[x][y] = newSegment

	// Try to stitch in this segment to its neighbour(s)
	switch segType {
	case EastWest:
		// I go east-west, which means I have partners to my west and to my east.
		// check (x-1,y), (x+1,y)
		// TODO: Does this work?
		if west := f.GetSegmentByXY(x-1, y); west != nil {
			west.East = newSegment
			newSegment.West = west
		}
		if east := f.GetSegmentByXY(x+1, y); east != nil {
			east.West = newSegment
			newSegment.East = east
		}
	case NorthSouth:
		if north := f.GetSegmentByXY(x, y-1); north != nil {
			north.South = newSegment
			newSegment.North = north
		}
		if south := f.GetSegmentByXY(x, y+1); south != nil {
			south.North = newSegment
			newSegment.South = south
		}
	// FIXME: check (S to W) or (N to E)
	// Also need to validate that the connection is valid:
	//
	// ----\/----
	// ----\+---
	// ----||
	//
	// In the first row, for example, 'East' isn't a valid connection so just
	// checking for that isn't sufficient. In fact, if there's a valid West
	// connection there may only be a South connection and if there's a valid North
	// connection there may only be an East connection.
	case TopLeftDiag:
		var west, east, north, south *Segment
		west = f.GetSegmentByXY(x-1, y)
		south = f.GetSegmentByXY(x, y+1)
		north = f.GetSegmentByXY(x, y-1)
		east = f.GetSegmentByXY(x+1, y)
		if *debug {
			fmt.Printf("Added a TopLeftDiag at (%d,%d). N/S/E/W: %+v / %+v / %+v / %+v\n", x, y, north, south, east, west)
		}
		if ((west != nil) && (west.Type == Intersection || west.Type == EastWest)) || ((south != nil) && (south.Type == Intersection || south.Type == NorthSouth)) {
			// it's a S to W
			if west != nil {
				west.East = newSegment
				newSegment.West = west
			}
			if south != nil {
				south.North = newSegment
				newSegment.South = south
			}
		} else if ((east != nil) && (east.Type == Intersection || east.Type == EastWest)) || ((north != nil) && (north.Type == Intersection || north.Type == NorthSouth)) {
			//it's a N to E
			if north != nil {
				if north.Type == Intersection || north.Type == NorthSouth {
					north.South = newSegment
					newSegment.North = north
				}
			}
			if east != nil {
				if east.Type == Intersection || east.Type == EastWest {
					east.West = newSegment
					newSegment.East = east
				}
			}
		}
		// check (x-1,y), (x,y+1)
	case TopRightDiag:
		// E to S or W to N
		var west, east, north, south *Segment
		west = f.GetSegmentByXY(x-1, y)
		south = f.GetSegmentByXY(x, y+1)
		north = f.GetSegmentByXY(x, y-1)
		east = f.GetSegmentByXY(x+1, y)
		if *debug {
			fmt.Printf("Added a TopRightDiag at (%d,%d). N/S/E/W: %+v / %+v / %+v / %+v\n", x, y, north, south, east, west)
		}
		if ((west != nil) && (west.Type == Intersection || west.Type == EastWest)) || ((north != nil) && (north.Type == Intersection || north.Type == NorthSouth)) {
			// it's a S to W
			if west != nil {
				if *debug {
					fmt.Printf("  West exists, setting its east to me and my west to them\n")
				}
				west.East = newSegment
				newSegment.West = west
			}
			if north != nil {
				if *debug {
					fmt.Printf("  North exists, setting its south to me and my west to them\n")
				}
				north.South = newSegment
				newSegment.North = north
			}
		} else if ((east != nil) && (east.Type == Intersection || east.Type == EastWest)) || ((south != nil) && (south.Type == Intersection || south.Type == NorthSouth)) {
			//it's a N to E
			if east != nil {
				east.West = newSegment
				newSegment.East = east
			}
			if south != nil {
				south.North = newSegment
				newSegment.South = south
			}
		}
		// check (x,y+1), (x+1,y)
	case Intersection:
		if west := f.GetSegmentByXY(x-1, y); west != nil {
			west.East = newSegment
			newSegment.West = west
		}
		if east := f.GetSegmentByXY(x+1, y); east != nil {
			east.West = newSegment
			newSegment.East = east
		}
		if north := f.GetSegmentByXY(x, y-1); north != nil {
			north.South = newSegment
			newSegment.North = north
		}
		if south := f.GetSegmentByXY(x, y+1); south != nil {
			south.North = newSegment
			newSegment.South = south
		}
	}
	if x > f.MaxX {
		f.MaxX = x
	}
	if y > f.MaxY {
		f.MaxY = y
	}
}

// AddShuttle - Adds a shuttle (okay, it's a Cart, but I like to think of the
// Elves as shuttling around.)
// This will also add the appropriate Segment.
func (f *Field) AddShuttle(x, y int, dir Direction) {
	var segType SegmentType
	switch dir {
	case West, East:
		segType = EastWest
	case North, South:
		segType = NorthSouth
	default:
		fmt.Printf("Unknown segment type for direction %d at (%d,%d)\n", dir, x, y)
		os.Exit(1)
	}
	f.AddSegment(x, y, segType)
	seg := f.GetSegmentByXY(x, y)

	s := &Shuttle{
		X:                 x,
		Y:                 y,
		DirectionOfTravel: dir,
		CurrentSegment:    seg,
	}
	f.Shuttles = append(f.Shuttles, s)
}

// GetSegmentByXY - return a segment, if it exists, identified by a specific
// (x,y) coordinate pair. If the segment does not exist, nil will be returne.
func (f *Field) GetSegmentByXY(x, y int) *Segment {
	if _, ok := f.Segments[x]; ok {
		return f.Segments[x][y]
	}
	return nil
}

func (f *Field) UpdateShuttlePositions() {
	for _, shuttle := range f.Shuttles {
		shuttle.X = shuttle.CurrentSegment.X
		shuttle.Y = shuttle.CurrentSegment.Y
	}
}

// HasCollision - does the Field have any collisions? If it does, return the
// segment at which it occurred, otherwise nil.
func (f *Field) HasCollision() *Segment {
	seen := make(map[int]map[int]bool)
	for _, shuttle := range f.Shuttles {
		if _, ok := seen[shuttle.CurrentSegment.X]; !ok {
			// init y->bool
			seen[shuttle.CurrentSegment.X] = make(map[int]bool)
		}
		if seen[shuttle.CurrentSegment.X][shuttle.CurrentSegment.Y] {
			if *debug {
				fmt.Printf("  %p at (%d,%d) has a collision. Segment type %s\n", shuttle, shuttle.CurrentSegment.X, shuttle.CurrentSegment.Y, segTypeToString(shuttle.CurrentSegment.Type))
			}
			fmt.Println()
			return shuttle.CurrentSegment
		} else {
			seen[shuttle.CurrentSegment.X][shuttle.CurrentSegment.Y] = true
		}

	}
	return nil
}

func (f *Field) HasCollisionAtXY(x, y int) bool {
	occupied := make(map[int]map[int]bool)
	for _, shuttle := range f.Shuttles {
		if shuttle.X != x || shuttle.Y != y {
			continue
		}
		fmt.Printf("HasCollisionAtXY(%d,%d): Examining shuttle (p=%p) %+v. Map -> %+v\n", x, y, shuttle, shuttle, occupied)
		if _, ok := occupied[shuttle.X]; !ok {
			occupied[shuttle.X] = make(map[int]bool)
		}
		if occupied[shuttle.X][shuttle.Y] {
			fmt.Printf("HasCollisionAtXY(%d,%d) Shuttle %+v has been seen? Map -> %+v\n", x, y, shuttle, occupied)
			return true
		} else {
			occupied[shuttle.X][shuttle.Y] = true
			fmt.Printf("HasCollisionAtXY(%d,%d) map (Setting seen) -> %+v\n", x, y, occupied)
		}
	}
	return false
}

// Tick - move all shuttles one by one checking for collisions along the way
// If there is a collision in this tick, return it, otherwise return nil.
func (f *Field) Tick() *Segment {

	for testY := 0; testY <= f.MaxY; testY++ {
		for _, shuttle := range f.GetShuttlesByY(testY) {
			if *debug {
				fmt.Printf("Moving shuttle %p from (%d,%d; segType=%s)\n", shuttle, shuttle.X, shuttle.Y, segTypeToString(shuttle.CurrentSegment.Type))
			}
			// The direction of movement for this tick was set by the previous tick.
			// That means that we will just move in that direction!
			switch shuttle.DirectionOfTravel {
			case East:
				shuttle.CurrentSegment = shuttle.CurrentSegment.East
				if *debug {
					fmt.Printf("    Moving East to (%d,%d) segType=%s\n", shuttle.CurrentSegment.X, shuttle.CurrentSegment.Y, segTypeToString(shuttle.CurrentSegment.Type))
				}
			case South:
				shuttle.CurrentSegment = shuttle.CurrentSegment.South
				if *debug {
					fmt.Printf("    Moving South to (%d,%d) segType=%s\n", shuttle.CurrentSegment.X, shuttle.CurrentSegment.Y, segTypeToString(shuttle.CurrentSegment.Type))
				}
			case West:
				shuttle.CurrentSegment = shuttle.CurrentSegment.West
				if *debug {
					fmt.Printf("    Moving West to (%d,%d) segType=%s\n", shuttle.CurrentSegment.X, shuttle.CurrentSegment.Y, segTypeToString(shuttle.CurrentSegment.Type))
				}
			case North:
				shuttle.CurrentSegment = shuttle.CurrentSegment.North
				if *debug {
					fmt.Printf("    Moving North to (%d,%d) segType=%s\n", shuttle.CurrentSegment.X, shuttle.CurrentSegment.Y, segTypeToString(shuttle.CurrentSegment.Type))
				}
			}
			if *debug {
				fmt.Printf("       Cheking for collisions\n")

			}
			collides := f.HasCollision()
			if collides != nil {
				return collides
			}
			if *debug {
				fmt.Printf("       Updating direction of travel for next Tick()...\n")
			}
			// Update the direction of travel for the next iteration
			switch shuttle.CurrentSegment.Type {
			case EastWest, NorthSouth:
				if *debug {
					fmt.Printf("          Straightaway, no change\n")
				}
				// do nothing because we can't change directions
			case TopLeftDiag:
				if *debug {
					fmt.Printf("          Top-left Diagional ")
				}
				switch shuttle.DirectionOfTravel {
				case South:
					shuttle.DirectionOfTravel = East
					if *debug {
						fmt.Printf("          Top-left Diagional South -> East\n")
					}
				case North:
					shuttle.DirectionOfTravel = West
					if *debug {
						fmt.Printf("          Top-left Diagional North -> West\n")
					}
				case West:
					shuttle.DirectionOfTravel = North
					if *debug {
						fmt.Printf("          Top-left Diagional West -> North\n")
					}
				case East:
					shuttle.DirectionOfTravel = South
					if *debug {
						fmt.Printf("          Top-left Diagional East -> South\n")
					}
				}
			case TopRightDiag:
				if *debug {
					fmt.Printf("          Top-right Diagional ")
				}
				switch shuttle.DirectionOfTravel {
				case East:
					shuttle.DirectionOfTravel = North
					if *debug {
						fmt.Printf("          Top-right Diagional East -> North\n")
					}
				case South:
					shuttle.DirectionOfTravel = West
					if *debug {
						fmt.Printf("          Top-right Diagional South -> West\n")
					}
				case North:
					shuttle.DirectionOfTravel = East
					if *debug {
						fmt.Printf("          Top-right Diagional North -> East\n")
					}
				case West:
					shuttle.DirectionOfTravel = South
					if *debug {
						fmt.Printf("          Top-right Diagional West -> South\n")
					}
				}
			case Intersection:
				// Order for this: Left, Straight, Right
				if *debug {
					fmt.Printf("          Intersection ")
				}
				switch shuttle.LastTurn {
				case Left:
					// last turn left, go straight ("None")
					shuttle.DirectionOfTravel = shuttle.MoveStraight()
					shuttle.LastTurn = Straight
					if *debug {
						fmt.Printf("          Intersection Went Straight -> Right next\n")
					}
				case Straight:
					// went straight last time, turn right
					shuttle.DirectionOfTravel = shuttle.TurnRight()
					shuttle.LastTurn = Right
					if *debug {
						fmt.Printf("          Intersection Went Right -> Straight next\n")
					}
				case Right:
					// went right, or made it is the first turn, turn left
					shuttle.DirectionOfTravel = shuttle.TurnLeft()
					shuttle.LastTurn = Left
					if *debug {
						fmt.Printf("          Intersection Went Left -> Straight next\n")
					}
				case 0:
					shuttle.DirectionOfTravel = shuttle.TurnLeft()
					shuttle.LastTurn = Left
					if *debug {
						fmt.Printf("          Intersection First encounter (turned Left) -> Straight next\n")
					}
				}
			}
		}
	}
	// Update out of the loop to avoid moving a shuttle twice in the above loop
	f.UpdateShuttlePositions()
	// Check for a collision - if we have one, return the Segment on which it occurred.
	return f.HasCollision()
}

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	errorIf("Can't open input file", err)

	defer input.Close()
	lineReader := bufio.NewScanner(input)

	field := NewField()
	y := 0
	for lineReader.Scan() {
		line := lineReader.Text()
		for x, r := range line {
			if *debug {
				fmt.Printf("coords = (%d,%d) ", x, y)
				fmt.Printf("r=%v\n", r)
			}

			var segType SegmentType
			var dir Direction
			hasShuttle := false
			switch r {
			case '/':
				// TopRightDiag
				segType = TopRightDiag
			case '\\':
				// TopLeftDiag
				segType = TopLeftDiag
			case '-':
				// EastWest
				segType = EastWest
			case '|':
				// NorthSouth
				segType = NorthSouth
			case '+':
				// Intersection
				segType = Intersection
			case '>':
				dir = East
				hasShuttle = true
				// Shuttle East (and EastWest)
			case '<':
				// Shuttle West (and EastWest)
				dir = West
				hasShuttle = true
			case 'v':
				// Shuttle South (and NorthSouth)
				dir = South
				hasShuttle = true
			case '^':
				// Shuttle North (and NorthSouth)
				dir = North
				hasShuttle = true
			default:
				// everything else, try the next character
				continue
			}

			// Now add a shuttle, or a segment
			if hasShuttle {
				field.AddShuttle(x, y, dir)
			} else {
				// adding a segment
				field.AddSegment(x, y, segType)
			}
			// next character
		}
		// Next line
		y++
	}

	iterations := 0
	var d string
	for {
		if *debug {
			field.PrintField()
		}
		if *debug2 {
			fmt.Printf("Press enter to advance")
			fmt.Scanf("\n", &d)
		}
		if collide := field.Tick(); collide == nil {
			if *debug {
				fmt.Printf("No collision. Shuttle positions\n")
			}
		} else {
			fmt.Printf("collision at (%d,%d)\n", collide.X, collide.Y)
			break
		}
		iterations++
		if *debug {
			fmt.Printf("Done %d iterations\n", iterations)
		}
	}
	if *debug {
		field.PrintField()
	}
	if *debug {
		fmt.Printf("Shuttles\n %+v\n", field.Shuttles)
	}
}
