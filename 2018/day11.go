package main

import (
	"flag"
	"fmt"
)

var (
	input = flag.Int("input", 1308, "puzzle input")
	partB = flag.Bool("partB", false, "do part b solution?")
	debug = flag.Bool("debug", false, "debug?")
)

type FuelCell struct {
	Cells map[int]map[int]*Cell
}

func (fc *FuelCell) GetCellByXY(x, y int) *Cell {
	return fc.Cells[x][y]
}

func InitFuelCell() *FuelCell {
	f := make(map[int]map[int]*Cell)
	for x := 1; x < 300; x++ {
		f[x] = make(map[int]*Cell)
		for y := 1; y < 300; y++ {
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
}
