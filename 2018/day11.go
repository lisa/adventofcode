package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	input = flag.Int("input", 1308, "puzzle input")
	partB = flag.Bool("partB", false, "do part b solution?")
	debug = flag.Bool("debug", false, "debug?")
)

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}

type FuelCell struct {
	Cells map[int]map[int]*Cell
}

func (fc *FuelCell) FindBestPowerBlock() (*Cell, int) {
	type BestBlock struct {
		cell  *Cell // Cell ID
		score int   // power of this range
		size  int   // Size of square out (n=0 is just the cell)
	}

	best := BestBlock{
		cell:  fc.GetCellByXY(1, 1),
		score: fc.GetCellByXY(1, 1).Power,
		size:  0,
	}
	for x := 1; x <= 300; x++ {
		for y := 1; y <= 300; y++ {
			for n := 0; ; n++ {
				if x+n > 300 || y+n > 300 {
					break
				}
				nBlockScore, err := fc.PowerNFromXY(x, y, n)

				if err != nil {
					// the proposed block out of the bounds? then w should move to the next x,y pair
					// and try again
					break
				}

				if nBlockScore > best.score {
					best.cell = fc.GetCellByXY(x, y)
					best.score = nBlockScore
					best.size = n
				}

			}
		}
	}
	return best.cell, best.size
}

// PowerNFromXY - Compute the power level n units out in each direction (eg a
// square) from (x,y). If n will exceed the bounds of the FuelCell (eg 300x300)
// a meaningless power level will be returned as well as an error.
// n can be 0, and if it is, it's the power of (x,y)
// n>0 means go to the right (increase x) and go down (increase y) by the same amount
func (fc *FuelCell) PowerNFromXY(x, y, n int) (int, error) {
	power := 0
	if n == 0 {
		return fc.GetCellByXY(x, y).Power, nil
	} else if n < 0 {
		return 0, errors.New("n must be >0")
	} else if x+n > 300 || y+n > 300 {
		return 0, errors.New("out of bounds")
	}

	for outX := x; outX < x+n; outX++ {
		for outY := y; outY < y+n; outY++ {
			p := fc.GetCellByXY(outX, outY).Power
			power += p
		}
	}

	return power, nil
}

func (fc *FuelCell) GetCellByXY(x, y int) *Cell {
	return fc.Cells[x][y]
}

func InitFuelCell() *FuelCell {
	f := make(map[int]map[int]*Cell)
	for x := 1; x <= 300; x++ {
		f[x] = make(map[int]*Cell)
		for y := 1; y <= 300; y++ {
			f[x][y] = NewCell(x, y)
		}
	}
	return &FuelCell{
		Cells: f,
	}
}

type Cell struct {
	X, Y  int
	Power int
}

func NewCell(x, y int) *Cell {
	c := &Cell{
		X: x,
		Y: y,
	}
	c.Power = c.PowerLevel()
	return c
}

func (c *Cell) RackID() int {
	return (c.X + 10)
}

// PowerLevel - get this cell's power level if we don't have it already
func (c *Cell) PowerLevel() int {
	return (((((c.RackID() * c.Y) + *input) * c.RackID() / 100) % 10) - 5)
}

func main() {
	flag.Parse()

	fuelCell := InitFuelCell()

	if !*partB {
		// Top left portion of the 3x3 section => its combined fuel power level
		score := make(map[*Cell]int)

		// Scan across X to section out 3x3 pieces. stop at x=297 since the 3x3 must be
		// within the bounds of the 300x300 field
		for x := 1; x < 297; x++ {
			for y := 1; y < 297; y++ {
				topLeft := fuelCell.GetCellByXY(x, y)
				power := topLeft.Power + fuelCell.GetCellByXY(x+1, y).Power + fuelCell.GetCellByXY(x+2, y).Power +
					fuelCell.GetCellByXY(x, y+1).Power + fuelCell.GetCellByXY(x+1, y+1).Power + fuelCell.GetCellByXY(x+2, y+1).Power +
					fuelCell.GetCellByXY(x, y+2).Power + fuelCell.GetCellByXY(x+1, y+2).Power + fuelCell.GetCellByXY(x+2, y+2).Power
				score[topLeft] = power
			}
		}

		highScore := fuelCell.GetCellByXY(1, 1).Power
		highScoreCell := fuelCell.GetCellByXY(1, 1)
		fmt.Printf("Starting high score: %d\n", highScoreCell)
		for cell, powerSum := range score {
			if powerSum > highScore {
				highScoreCell = cell
				highScore = powerSum
			}
		}

		fmt.Printf("Cell at top-left (%d,%d) has total power %d\n", highScoreCell.X, highScoreCell.Y, highScore)
	} else {
		bcell, bsize := fuelCell.FindBestPowerBlock()
		fmt.Printf("Part B: (%d,%d) n=%d\n", bcell.X, bcell.Y, bsize)
	}
}
