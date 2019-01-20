package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/kr/pretty"
	. "github.com/logrusorgru/aurora"
)

var (
	inputFile    = flag.String("input", "inputs/day15.txt", "Input file")
	partB        = flag.Bool("partB", false, "Part B solution")
	debug        = flag.Bool("debug", false, "debug?")
	debug2       = flag.Bool("debug2", false, "require user input to continue")
	debug3       = flag.Bool("debug3", false, "..")
	debugpathing = flag.Bool("debugPathing", false, "debug pathing?")
)

type PathList struct {
	this   Tile
	parent *PathList
	score  int
}

// TileType - a specific kind of tile
type TileType int

const (
	WallTile TileType = iota
	GoblinTile
	ElfTile
	OpenTile
)

// DistanceTo - cartesian distance between Tiles
func DistanceTo(from, to Tile) int {
	fromX, fromY := from.Coords()
	toX, toY := to.Coords()

	return int(math.Abs(float64(fromX-toX)) + math.Abs(float64(fromY-toY)))
}

// Less - Comparison for two tiles, compare on health (if applicable) and then
// reading order less means less health and/or closer to top-left in reading
// order
func Less(left, right Tile) bool {
	leftHealth := left.GetHealth()
	rightHealth := right.GetHealth()

	// if health are the same (as in for Walls and Open spaces), check reading order
	if leftHealth == rightHealth {
		//reading order check
		leftX, leftY := left.Coords()
		rightX, rightY := left.Coords()
		if leftY < rightY {
			return true
		} else if leftY > rightY {
			return false
		} else {
			// Y are equal
			if leftX < rightX {
				return true
			}
			return false
		}
	}

	// different healths, so just check that
	return leftHealth < rightHealth
}

// Tile - represents a tile in the field
type Tile interface {
	Kind() TileType
	Coords() (int, int)
	Eql(Tile) bool
	GetHealth() int
	SetHealth(int) Tile
	GetPower() int
	UpdateCoords(int, int) // update the coords (DOES NOT CHECK FOR )
}
type Elf struct {
	X, Y, Power, Health int
}

// UpdateCoords - update the coords for this Tile to the specified x,y pair
func (t Elf) UpdateCoords(newX, newY int) {
	fmt.Printf("Elf updating coords to (%d,%d)", newX, newY)
	t.X = newX
	t.Y = newY
	fmt.Printf(" -> %+v\n", t)
}

// UpdateCoords - update the coords for this Tile to the specified x,y pair
func (t Wall) UpdateCoords(newX, newY int) {
	t.X = newX
	t.Y = newY
}

// UpdateCoords - update the coords for this Tile to the specified x,y pair
func (t Open) UpdateCoords(newX, newY int) {
	t.X = newX
	t.Y = newY
}

// UpdateCoords - update the coords for this Tile to the specified x,y pair
func (t Goblin) UpdateCoords(newX, newY int) {
	fmt.Printf("Goblin updating coords to (%d,%d)", newX, newY)
	t.X = newX
	t.Y = newY
	fmt.Printf(" -> %+v\n", t)
}

// Kind - what kind of Tile am I?
func (e Elf) Kind() TileType { return ElfTile }
func (e Elf) Coords() (int, int) {
	return e.X, e.Y
}
func (e Elf) Eql(t Tile) bool {
	x, y := t.Coords()
	return e.X == x && e.Y == y
}

func (e Elf) GetHealth() int {
	return e.Health
}
func (e Elf) GetPower() int {
	return e.Power
}

func (e Elf) Less(t Tile) bool {
	myX, myY := e.Coords()
	oX, oY := t.Coords()
	switch t.Kind() {
	// for these, just comparing reading order
	case WallTile, OpenTile:
		if myY < oY {
			return true
		} else if myY > oY {
			return false
		} else {
			// on same Y, check X
			if myX < oX {
				return true
			} else if myX > oX {
				return true
			} else {
				// technically this is equality
				return false
			}
		}
	case GoblinTile:
		// check health

	}
	return false
}

func (e Elf) SetHealth(newHealth int) Tile {
	e.Health = newHealth
	return e
}

// Kind - what kind of Tile am I?
func (w Wall) SetHealth(h int) Tile { return w }
func (w Wall) Kind() TileType       { return WallTile }
func (w Wall) Coords() (int, int) {
	return w.X, w.Y
}
func (w Wall) Eql(t Tile) bool {
	x, y := t.Coords()
	return w.X == x && w.Y == y
}

func (w Wall) GetHealth() int {
	return -1
}
func (w Wall) GetPower() int { return 0 }

// Kind - what kind of Tile am I?
func (o Open) Kind() TileType { return OpenTile }
func (o Open) Coords() (int, int) {
	return o.X, o.Y
}
func (o Open) Eql(t Tile) bool {
	x, y := t.Coords()
	return o.X == x && o.Y == y
}
func (o Open) GetHealth() int {
	return -1
}
func (o Open) GetPower() int        { return 0 }
func (o Open) SetHealth(h int) Tile { return o }

// Kind - what kind of Tile am I?
func (g Goblin) Kind() TileType { return GoblinTile }
func (g Goblin) Coords() (int, int) {
	return g.X, g.Y
}
func (g Goblin) Eql(t Tile) bool {
	x, y := t.Coords()
	return g.X == x && g.Y == y
}

func (g Goblin) SetHealth(newHealth int) Tile {
	g.Health = newHealth
	return g
}
func (g Goblin) GetHealth() int {
	return g.Health
}
func (g Goblin) GetPower() int { return g.Power }

type Goblin struct {
	X, Y, Power, Health int
}

type Wall struct {
	X, Y int
}

type Open struct {
	X, Y int
}

// Field - The playing Field
type Field struct {
	Tiles      map[int]map[int]Tile
	MaxX, MaxY int
}

func NewField() *Field {
	return &Field{
		Tiles: make(map[int]map[int]Tile),
		MaxX:  0,
		MaxY:  0,
	}
}

func tileTypeToString(t TileType) string {
	switch t {
	case OpenTile:
		return "Open"
	case WallTile:
		return "Wall"
	case ElfTile:
		return "Elf"
	case GoblinTile:
		return "Goblin"
	}
	return "Unknown Type"
}

// AddTile - adds a Tile to the Field
func (f *Field) AddTile(x, y int, t TileType) {
	if _, ok := f.Tiles[x]; !ok {
		f.Tiles[x] = make(map[int]Tile)
	}
	switch t {
	case WallTile:
		f.Tiles[x][y] = Wall{X: x, Y: y}
		// wall
	case OpenTile:
		f.Tiles[x][y] = Open{X: x, Y: y}
		// open
	case ElfTile:
		f.Tiles[x][y] = Elf{X: x, Y: y, Health: 200, Power: 3}
		// elf
	case GoblinTile:
		// goblin
		f.Tiles[x][y] = Goblin{X: x, Y: y, Health: 200, Power: 3}
	}
	if x > f.MaxX {
		f.MaxX = x
	}
	if y > f.MaxY {
		f.MaxY = y
	}
}

func (f *Field) GetNorthTile(t Tile) (Tile, error) {
	x, y := t.Coords()

	if y-1 < 0 {
		return t, errors.New("No tile north of given tile")
	}
	return f.Tiles[x][y-1], nil
}

func (f *Field) GetWestTile(t Tile) (Tile, error) {
	x, y := t.Coords()

	if x-1 < 0 {
		return t, errors.New("No tile west of given tile")
	}
	return f.Tiles[x-1][y], nil
}
func (f *Field) GetEastTile(t Tile) (Tile, error) {
	x, y := t.Coords()

	if x+1 > f.MaxX {
		return t, errors.New("No tile east of given tile")
	}
	return f.Tiles[x+1][y], nil
}
func (f *Field) GetSouthTile(t Tile) (Tile, error) {
	x, y := t.Coords()

	if y+1 > f.MaxY {
		return t, errors.New("No tile south of given tile")
	}
	return f.Tiles[x][y+1], nil
}

// OpenNeighboursOf - return all of the open neighbours of the given tile.
// Tiles will be returned in order: north, west, east, south. If a direction has
// no valid neighbouring tile it will simply not be present; this isn't always a
// 4 length list.
func (f *Field) OpenNeighboursOf(t Tile) []Tile {
	ret := make([]Tile, 0)
	if north, err := f.GetNorthTile(t); err == nil {
		if north.Kind() == OpenTile {
			ret = append(ret, north.(Open))
		}
	}
	if west, err := f.GetWestTile(t); err == nil {
		if west.Kind() == OpenTile {
			ret = append(ret, west.(Open))
		}
	}
	if east, err := f.GetEastTile(t); err == nil {
		if east.Kind() == OpenTile {
			ret = append(ret, east.(Open))
		}
	}
	if south, err := f.GetSouthTile(t); err == nil {
		if south.Kind() == OpenTile {
			ret = append(ret, south.(Open))
		}
	}
	return ret
}
func intSliceContains(needle int, haystack []int) bool {
	for i := 0; i < len(haystack); i++ {
		if haystack[i] == needle {
			return true
		}
	}
	return false
}

func (f *Field) printPathMap(p PathList) {
	poi := make(map[int][]int)
	x, y := p.this.Coords()
	startX, startY := p.this.Coords()
	endX, endY := p.this.Coords()
	poi[x] = append(poi[x], y)
	for parent := p.parent; parent != nil; parent = parent.parent {
		x, y = parent.this.Coords()
		// pathmap is reversed
		startX, startY = parent.this.Coords()
		poi[x] = append(poi[x], y)
	}
	for fieldY := 0; fieldY <= f.MaxY; fieldY++ {
		for fieldX := 0; fieldX <= f.MaxX; fieldX++ {
			t := f.Tiles[fieldX][fieldY]
			var symbol string
			switch t.Kind() {
			case OpenTile:
				symbol = "."
			case WallTile:
				symbol = "#"
			case ElfTile:
				symbol = "E"
			case GoblinTile:
				symbol = "G"
			}
			if _, ok := poi[fieldX]; ok && intSliceContains(fieldY, poi[fieldX]) {
				if fieldX == startX && fieldY == startY {
					// start is green
					fmt.Printf("%s", Green(symbol).Bold())
				} else if fieldX == endX && fieldY == endY {
					// end is red
					fmt.Printf("%s", Red(symbol).Bold())
				} else {
					fmt.Printf("%s", Blue(symbol).Bold())
				}
				// draw mark
			} else {
				fmt.Printf(symbol)
			}
		}
		fmt.Println()
	}

}

func printToParent(p PathList) {
	x, y := p.this.Coords()
	fmt.Printf("Path to parent (score=%d): (%d,%d)", p.score, x, y)
	for parent := p.parent; parent != nil; parent = parent.parent {
		x, y := parent.this.Coords()
		fmt.Printf(" -> (%d,%d)", x, y)
	}
	fmt.Println()
}

func pathListHasTile(needle Tile, haystack *[]Tile) bool {
	for i := 0; i < len(*haystack); i++ {
		if (*haystack)[i].Eql(needle) {
			return true
		}
	}
	return false
}

// ChangeTileCoords - try to move the given tile's coordinates to the given x,y coords
// this will fail if the move is to an non-open tile
func (f Field) ChangeTileCoords(from Tile, toX, toY int) bool {
	if f.Tiles[toX][toY].Kind() != OpenTile {
		return false
	}

	oldX, oldY := from.Coords()
	fmt.Printf("Updating %+v to be at (%d,%d)\n", from, toX, toY)
	f.Tiles[oldX][oldY] = Open{X: oldX, Y: oldY}
	from.UpdateCoords(toX, toY)
	fmt.Printf("From after update coords -> %+v\n", from)
	f.Tiles[toX][toY] = from
	return true
}

// nextToDestination - is the current tile next to the destination? If so, return true
func (f *Field) nextToDestination(current, destination Tile) bool {
	if t, err := f.GetNorthTile(current); err == nil && t.Eql(destination) {
		return true
	}
	if t, err := f.GetWestTile(current); err == nil && t.Eql(destination) {
		return true
	}
	if t, err := f.GetEastTile(current); err == nil && t.Eql(destination) {
		return true
	}
	if t, err := f.GetSouthTile(current); err == nil && t.Eql(destination) {
		return true
	}
	return false
}

// MoveTo - Move from (fromX,fromY) to (toX,toY). Return the list of
// Tiles that achieves the destination or an error if it isn't possible.
// Movement Rules: Prefer "reading order" (top to bottom, left to right) when
// possible. "(toX,toY)" is likely to be a tile adjacent to an enemy.
func (f *Field) MoveTo(fromX, fromY, toX, toY int) ([]Tile, error) {

	openList := make([]PathList, 0)
	closedList := make([]Tile, 0)
	path := make([]PathList, 0)

	root := f.Tiles[fromX][fromY]
	final := f.Tiles[toX][toY]
	if *debugpathing {
		fmt.Printf("Root = %+v, Final = %+v. Path list= %+v\n", root, final, path)
	}
	//	var s string
	// list of Tiles to consider on our way to the destination
	openList = append(openList, PathList{this: root, score: 0.0})
outer:
	for len(openList) > 0 {
		if *debugpathing {
			fmt.Printf("Open list (before slicing) len=%d Best score=%f; worst=%f", len(openList), openList[0].score, openList[len(openList)-1].score)
			fmt.Printf("  Closed list (size=%d) - Press enter to continue\n", len(closedList))
		}
		// Slice off the first Tile in our consideration list
		currentTile := openList[0]
		openList = openList[1:]
		currX, currY := currentTile.this.Coords()
		if *debugpathing {
			fmt.Printf("Current Tile(%d,%d) = %+v (p=%p)\n", currX, currY, currentTile, &currentTile)
		}

		if *debugpathing {
			printToParent(currentTile)
		}

		// have we visited this Tile already and ruled it out?
		for i := 0; i < len(closedList); i++ {
			if closedList[i].Eql(currentTile.this) {
				continue outer
			}
		}
		closedList = append(closedList, currentTile.this)

		if *debugpathing {
			fmt.Printf("Are we there yet with (%d,%d) ", currX, currY)
		}
		// Are we there yet?
		if currentTile.this.Eql(final) {
			// TODO: rewind through parents til nil
			returnPath := make([]Tile, 0)
			returnPath = append(returnPath, currentTile.this)
			parent := currentTile.parent
			for parent != nil {
				returnPath = append(returnPath, parent.this)
				parent = parent.parent
			}

			if *debugpathing || *debug3 {
				if !*debug3 {
					fmt.Printf(" Yes ")
				}
				printToParent(currentTile)
				fmt.Println("Map for above route follows:")
				f.printPathMap(currentTile)
			}

			return returnPath, nil
		}
		if *debugpathing {
			fmt.Printf("No, so continue.\n")
		}

		// Yay! A new node! Let's explore its neighbours :) We can only move on open
		// tiles so we will only consider them.
		//fmt.Printf("Getting neighbours of (%d,%d) (type=%d) ", x, y, currentTile.this.Kind())
		neighbours := f.OpenNeighboursOf(currentTile.this)
		// We need to do some special considerations since our neighbour might be the destination, so let's check.
		if f.nextToDestination(currentTile.this, final) {
			// fake that it's an Open node for safety sake
			// TODO: fix this
			neighbours = append(neighbours, final)
		}

		for _, openNeighbouringTile := range neighbours {
			if pathListHasTile(openNeighbouringTile, &closedList) {
				continue
			}
			newTile := PathList{
				this:   openNeighbouringTile,
				parent: &currentTile,
				score:  0,
			}
			if closedListHas(&newTile, &closedList) {
				fmt.Printf("Skipping %+v\n", openNeighbouringTile)
				continue
			}

			newTile.score = f.movementScore(&newTile, final)

			openList = append(openList, newTile)
		}
		// TODO Fix this sort so that on equal scores, it will check reading order
		sort.Slice(openList, func(i, j int) bool {
			if openList[i].score == openList[j].score {
				//reading order check
				leftX, leftY := openList[i].this.Coords()
				rightX, rightY := openList[j].this.Coords()
				if leftY < rightY {
					return true
				} else if leftY > rightY {
					return false
				} else {
					// Y are equal
					if leftX < rightX {
						return true
					}
					return false
				}

			} else {
				return openList[i].score < openList[j].score
			}
		})
		//		if *debug2 {
		//			fmt.Printf("After sorting: %+v\n\n\n", openList)
		//		}
	}
	if *debug {
		fmt.Printf("End of open list")
	}

	return []Tile{}, errors.New("Couldn't find a path")
}

func dedupTileSlice(s []Tile) []Tile {
	seen := make(map[Tile]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

// closedListHas - does the closed list have this item yet?
func closedListHas(needle *PathList, haystack *[]Tile) bool {
	for _, item := range *haystack {
		if item.Eql(needle.this) {
			return true
		}
	}
	return false
}

// Tick - Do one Round for each unit
// 1. For each unit in reading order for the units' STARTING POSITIONS for this Tick:
//  If no enemies are in range, move to an enemyand then attack
//  If enemies are in range, attack
func (f *Field) Tick(c *int) {
	// Create a list of all the goblins and elves in reading order and then process them

	units := make([]Tile, 0)

	if *debug {
		fmt.Printf("Collecting all elves and goblins\n")
	}
	elves := 0
	goblins := 0

	for fieldY := 0; fieldY <= f.MaxY; fieldY++ {
		for fieldX := 0; fieldX <= f.MaxX; fieldX++ {
			t := f.Tiles[fieldX][fieldY]
			if t.Kind() == ElfTile {
				elves++
				units = append(units, t)
			}
			if t.Kind() == GoblinTile {
				goblins++
				units = append(units, t)
			}
		}
	}
	if elves == 0 || goblins == 0 {
		return
	}

	var s string
	previousX := -1
	previousY := -1
	// With the list of moveables in reading order (top to bottom, left to right) ...
	for i, unit := range units {
		myX, myY := unit.Coords()
		if *debug3 {
			fmt.Printf("Press enter to continue\n")
			fmt.Scanf("\n", &s)

			fmt.Printf("Before working with Unit %+v map is:\n", unit)
			f.printField2(myX, myY, previousX, previousY)
		}

		if *debug {
			fmt.Printf("\n[%d/%d] Tick loop unit=%+v (real tile=%+v; type=%s)\n", i+1, len(units), unit, f.Tiles[myX][myY], tileTypeToString(f.Tiles[myX][myY].Kind()))
		}
		if f.Tiles[myX][myY].Kind() == OpenTile {
			if *debug {
				fmt.Printf("I'm dead :'(\n")
			}
			previousX = myX
			previousY = myY
			continue
		}

		didCombat := f.CombatFor(unit)

		var opposite TileType

		// If I'm an Elf, I care to match for Goblins, and vice versa
		switch unit.Kind() {
		case ElfTile:
			opposite = GoblinTile
		case GoblinTile:
			opposite = ElfTile
		}

		// did we have anything in range before having to move?
		if didCombat {
			previousX = myX
			previousY = myY
			continue
		} else {
			// did not have anything in range at the start of our turn, so we'll try to move to find something.
			// So let's get all open tiles around each opposite unit, and of all those open
			// tiles, find the shortest path.

			// so must loop through all the units again
			possibleMoves := make([]Tile, 0)
			for _, moveUnit := range units {
				if moveUnit.Kind() == opposite {
					possibleMoves = append(possibleMoves, f.OpenNeighboursOf(moveUnit)...)
				}
			}
			possibleMoves = dedupTileSlice(possibleMoves)
			if *debug {
				fmt.Printf("  No combat for me, so need to move to attack. Possible moves = %+v\n", possibleMoves)
			}

			// for each of the possible tiles we could move (adjacent to an enemy unit), which has the shortest path? so, path to them.
			allRoutes := make([][]Tile, 0)
			for i, possibleMove := range possibleMoves {
				pX, pY := possibleMove.Coords()
				if *debug {
					fmt.Printf("[%02d/%02d] Trying to route (%d,%d) -> (%d,%d) ", i+1, len(possibleMoves), myX, myY, pX, pY)
				}
				// Try to reverse route just in case the destination is blocked, we will uncover it quicker.
				//if route, err := f.MoveTo(pX, pY, myX, myY); err == nil {
				if route, err := f.MoveTo(myX, myY, pX, pY); err == nil {

					if *debug {
						fmt.Printf("success -> len=%d\n", len(route))
						fmt.Printf("  Success -> %+v\n", route)
					}
					allRoutes = append(allRoutes, route)
				} else {
					if *debug {
						fmt.Printf("No path :(\n")
					}
				}
			}

			if *debug {
				fmt.Printf("  All possible routes = %# v\n", pretty.Formatter(allRoutes))
			}

			sort.Slice(allRoutes, func(i, j int) bool {
				// Pick out the preferred route. Normally this would be whichever is shortest,
				// but if one is shortest we should prefer the one in reading order first.
				if len(allRoutes[i]) == len(allRoutes[j]) {
					// reading order of destination
					leftX, leftY := allRoutes[i][0].Coords()
					rightX, rightY := allRoutes[j][0].Coords()
					if leftY < rightY {
						return true
					} else if leftY > rightY {
						return false
					} else {
						// Y are equal
						if leftX < rightX {
							return true
						}
						return false
					}
				}
				return len(allRoutes[i]) < len(allRoutes[j])
			})

			if *debug {
				fmt.Printf("  All possible routes (sorted) = %# v\n", pretty.Formatter(allRoutes))
				fmt.Printf("  All possible destinations = %+v\n", possibleMoves)

				fmt.Printf("Pausing after displaying all sorted routes. Press enter to continue...\n")
				//fmt.Scanf("\n", &s)
			}
			if len(allRoutes) == 0 {
				if *debug3 {
					fmt.Printf("No valid moves for %s at (%d,%d). Its turn is forfeit.\n", tileTypeToString(unit.Kind()), myX, myY)
				}
				// we have no path to any enemy units so we have forfeit our turn
				previousX = myX
				previousY = myY
				continue
			}

			shortest := allRoutes[0]

			// we want to move to `shortest` by taking the first square to it, which is allRoutes[0][0].

			// for forward searching (last element is `unit`, second to last is where it moves next)
			firstMove := shortest[len(shortest)-2 : len(shortest)-1][0]

			// reverse search
			//firstMove := shortest[1]
			firstX, firstY := firstMove.Coords()
			if *debug3 {
				fmt.Printf("Moved %s at (%d,%d) to (%d,%d) long shortest=%+v\n", tileTypeToString(unit.Kind()), myX, myY, firstX, firstY, shortest)
			}
			if *debug {
				fmt.Printf("Shortest path is %+v, first move is (%d,%d)\n", shortest, firstX, firstY)
			}
			switch unit.Kind() {
			case ElfTile:
				f.Tiles[firstX][firstY] = Elf{
					X: firstX, Y: firstY,
					Health: unit.GetHealth(),
					Power:  unit.GetPower(),
				}
			case GoblinTile:
				f.Tiles[firstX][firstY] = Goblin{
					X: firstX, Y: firstY,
					Health: unit.GetHealth(),
					Power:  unit.GetPower(),
				}
			}
			f.AddTile(myX, myY, OpenTile)
			if *debug {
				fmt.Printf("Checking combat after move\n")
			}
			unit = f.Tiles[firstX][firstY]
			f.CombatFor(unit)

			elves := 0
			goblins := 0

			for fieldY := 0; fieldY <= f.MaxY; fieldY++ {
				for fieldX := 0; fieldX <= f.MaxX; fieldX++ {
					t := f.Tiles[fieldX][fieldY]
					if t.Kind() == ElfTile {
						elves++
						units = append(units, t)
					}
					if t.Kind() == GoblinTile {
						goblins++
						units = append(units, t)
					}
				}
			}
			if elves == 0 || goblins == 0 {
				return
			}

			previousX = myX
			previousY = myY

			// end "did combat" else
		}
	}
	*c++

}

// CombatFor - Perform combat for a given tile:
// For the Tile +unit+, do all combat actions for its attacks. If there are no
// enemies in range, return false, otherwise return true. This method does not
func (f Field) CombatFor(unit Tile) bool {
	// Identify all possible targets for this unit.
	// Goblins want Elves and Elves want Goblins.
	var opposite TileType

	// If I'm an Elf, I care to match for Goblins, and vice versa
	switch unit.Kind() {
	case ElfTile:
		opposite = GoblinTile
	case GoblinTile:
		opposite = ElfTile
	}

	if *debug {
		fmt.Printf("CombatFor Opposite tile type to me (%s -> %s)\n", tileTypeToString(unit.Kind()), tileTypeToString(opposite))
	}

	// Are there any enemy units in range of me?
	nearTiles := make([]Tile, 0)
	if north, err := f.GetNorthTile(unit); err == nil && north.Kind() == opposite {
		nearTiles = append(nearTiles, north)
	}
	if west, err := f.GetWestTile(unit); err == nil && west.Kind() == opposite {
		nearTiles = append(nearTiles, west)
	}
	if east, err := f.GetEastTile(unit); err == nil && east.Kind() == opposite {
		nearTiles = append(nearTiles, east)
	}
	if south, err := f.GetSouthTile(unit); err == nil && south.Kind() == opposite {
		nearTiles = append(nearTiles, south)
	}
	if *debug {
		fmt.Printf("CombatFor have %+v nearTiles\n", nearTiles)
	}
	if len(nearTiles) == 0 {
		// nothing in range, so we do nothing
		return false
	}

	var lowestPower int
	var lowestTile Tile
	// we are in range of something to attack, but let's figure out which it is
	sort.Slice(nearTiles, func(i, j int) bool {
		if Less(nearTiles[i], nearTiles[j]) {
			// new low score
			lowestPower = nearTiles[i].GetHealth()
			lowestTile = nearTiles[i]
			return true
		}
		return false
	})
	if *debug {
		fmt.Printf("CombatFor sorted nearTiles -> %+v\n", nearTiles)
		// we're in range! ATTAAAAAACCKKK!!!!!!!!!
		fmt.Printf("Weakest is %v\n", nearTiles[0])
	}

	// deal damage
	nearX, nearY := nearTiles[0].Coords()
	// what will the new health of our victim be?
	newHealth := nearTiles[0].GetHealth() - unit.GetPower()
	// is it dead?
	if newHealth <= 0 {
		// it dead
		if *debug {
			fmt.Printf("CombatFor - victim is ded\n")
		}
		f.AddTile(nearX, nearY, OpenTile)
	} else {
		// not dead yet!
		if *debug {
			fmt.Printf("CombatFor - not dead yet\n")
		}
		f.Tiles[nearX][nearY] = nearTiles[0].SetHealth(newHealth)
		if *debug {
			fmt.Printf("CombatFor - not dead yet -> %+v\n", f.Tiles[nearX][nearY])
		}
	}
	return true
}

// movementScore - what's the score to move from a to b?
// computes g and h scores, returns a float32 with the score
func (f *Field) movementScore(from *PathList, to Tile) int {
	g := 0
	if *debugpathing {
		fmt.Printf("\tComputing score from %+v to %+v\n", from, to)
	}

	// g - distance to our starting point
	parent := from.parent
	if parent != nil {
		if *debugpathing {
			fmt.Printf("\t\t(g) Adding parent score %f\n", parent.score)
		}
		g += parent.score + 1
	}

	// h
	// to only do Djikstra's, h is 0 (eg, only consider g)
	h := DistanceTo(from.this, to)

	if *debugpathing {
		fmt.Printf("\t\tFinal score=g+h; %d=%d+%d\n", g+h, g, h)
	}
	return g + h
}

func (f Field) printField() {
	f.printField2(-1, -1, -1, -1)
}

func (f Field) printField2(x, y, ox, oy int) {
	for fieldY := 0; fieldY <= f.MaxY; fieldY++ {
		for fieldX := 0; fieldX <= f.MaxX; fieldX++ {
			t := f.Tiles[fieldX][fieldY]
			var symbol string
			switch t.Kind() {
			case OpenTile:
				symbol = "."
			case WallTile:
				symbol = "#"
			case ElfTile:
				symbol = "E"
			case GoblinTile:
				symbol = "G"
			}
			if x == fieldX && y == fieldY {
				fmt.Printf("%s", Cyan(symbol).Bold().Inverse())
			} else if ox == fieldX && oy == fieldY {
				fmt.Printf("%s", Red(symbol))
			} else {
				fmt.Printf(symbol)
			}
		}
		for fieldX := 0; fieldX <= f.MaxX; fieldX++ {
			t := f.Tiles[fieldX][fieldY]
			switch t.Kind() {
			case ElfTile:
				fmt.Printf(" E(%d),", t.GetHealth())
			case GoblinTile:
				fmt.Printf(" G(%d),", t.GetHealth())
			}
		}
		fmt.Println()
	}
}

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
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
		for x, c := range line {
			var t TileType
			switch c {
			case '#':
				// wall
				t = WallTile
			case '.':
				// empty
				t = OpenTile
			case 'G':
				// goblin
				t = GoblinTile
			case 'E':
				// elf
				t = ElfTile
			default:
				// unknown
				fmt.Printf("Unknown tile type %s\n", c)
			}
			field.AddTile(x, y, t)
		}
		// EOL
		y++
	}

	field.printField()
	fmt.Println()

	//var s string
	called := 0
	for {
		totalHealth := 0
		totalElves := 0
		totalGoblins := 0

		field.Tick(&called)
		fmt.Printf("Finished turn %d (Map directly after this line), press enter to see map for turn %d\n", called, called+1)
		//if *debug {

		field.printField()
		fmt.Println()
		//fmt.Scanf("\n", &s)

		//		}
		for fieldY := 0; fieldY <= field.MaxY; fieldY++ {
			for fieldX := 0; fieldX <= field.MaxX; fieldX++ {
				t := field.Tiles[fieldX][fieldY]
				switch t.Kind() {
				case ElfTile:
					totalElves++
					totalHealth += t.GetHealth()
					if *debug {
						fmt.Printf(" Post-Tick Unit ->  %+v (type %s)\n", t, tileTypeToString(t.Kind()))
					}
				case GoblinTile:
					totalGoblins++
					totalHealth += t.GetHealth()
					if *debug {
						fmt.Printf(" Post-Tick Unit ->  %+v (type %s)\n", t, tileTypeToString(t.Kind()))
					}

				}
			}
		}
		if totalElves > 0 && totalGoblins > 0 {
			//still some alive
			continue
		} else {
			fmt.Printf("Done after %d turns. %d elves and %d goblins left. tally=%d (totalHealth=%d)\n", called, totalElves, totalGoblins, called*totalHealth, totalHealth)
			break
		}
	}
}
