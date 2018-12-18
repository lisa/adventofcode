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
	inputFile = flag.String("input", "inputs/day09.txt", "input file")
	partB     = flag.Bool("partB", false, "do part b solution?")
	debug     = flag.Bool("debug", false, "debug?")
	debug2    = flag.Bool("debug2", false, "more debug")
	removed   = make(map[int]int)
)

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}

// InitCircle - Start the loop
func InitCircle() *Marble {
	r := NewMarble(0)
	r.Left = r
	r.Right = r
	return r
}

// NewMarble - make a new Marble
func NewMarble(v int) *Marble {
	return &Marble{
		Value: v,
		Left:  nil,
		Right: nil,
	}
}

// Marble - the Marble object; a doubly linked list.
type Marble struct {
	Right *Marble // clockwise
	Left  *Marble // counterclockwise
	Value int     // value of this marble
}

// Move - Traverse +offset+ Marbles.
// If +offset+ is positive, move that many to the right (clockwise)
// If +offset+ is negative, move that many to the left (counter-clockwise)
func (m *Marble) Move(offset int) *Marble {
	ret := m
	if offset == 0 {
		return ret
	}
	for i := 0; i < int(math.Abs(float64(offset))); i++ {
		if offset < 0 {
			ret = ret.Left
			//fmt.Printf("Moving left -> %+v\n", ret)
		} else {
			ret = ret.Right
			//fmt.Printf("Moving right -> %+v\n", ret)
		}
	}
	return ret
}

// Delete - Delete the current Marble; its Value is returned. Left and Right associations are carried over
func (m *Marble) Delete() int {
	ret := m.Value
	l := m.Left
	r := m.Right
	r.Left = l
	l.Right = r
	return ret
}

// InsertRight - insert a Marble with value +v+ to the Right (clockwise) of myself. Return the new Marble
// new Marble is to my right
func (m *Marble) InsertRight(v int) *Marble {
	n := NewMarble(v)

	// new marble and old "right of me"
	m.Right.Left = n
	n.Right = m.Right

	n.Left = m
	m.Right = n

	return n
}

// InsertLeft - insert a Marble with value +v+ to the Left (counterclockwise) of myself. Return the new Marble.
func (m *Marble) InsertLeft(v int) *Marble {
	n := NewMarble(v)

	m.Left.Right = n
	n.Left = m.Left

	m.Left = n
	n.Right = m

	return n
}

// PrintCircle - Print a representation of the Marble circle
func PrintCircle(current *Marble) {
	currentIntVal := current.Value
	currentPtrVal := current

	for {
		current = current.Right
		if current.Value == currentIntVal {
			fmt.Printf("(%d)\n", current.Value)
			break
		} else {
			fmt.Printf("%d ", current.Value)
		}
	}
	current = currentPtrVal
}

// PartA - do part A and return the scores
func PartA(players, lastMarbleValue int) map[int]int {
	current := InitCircle()

	// player -> score
	score := make(map[int]int)

	var del *Marble
	for v := 1; v <= lastMarbleValue; v++ {
		if v%23 == 0 {
			del = current.Move(-7)
			current = del.Right
			score[v%players] += v + del.Delete()
		} else {
			current = current.Right.InsertRight(v)
		}
	}
	return score

}
func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	errorIf("Can't open input file", err)

	defer input.Close()
	lineReader := bufio.NewScanner(input)
	var players, lastMarbleValue int
	if lineReader.Scan() {
		words := strings.Split(lineReader.Text(), " ")
		players, err = strconv.Atoi(words[0])
		errorIf("Couldn't parse the number of players\n", err)
		lastMarbleValue, err = strconv.Atoi(words[6])
		errorIf("Couldn't parse last marble score\n", err)
	}

	fmt.Printf("Player count %d, highest marble value %d\n", players, lastMarbleValue)

	if !*partB {
		score := PartA(players, lastMarbleValue)

		highScore := -1
		for i := range score {
			if score[i] > highScore {
				highScore = score[i]
			}
		}
		fmt.Printf("Part A high score %d\n", highScore)
	} else {
		highScore := -1
		score := PartA(players, lastMarbleValue*100)
		for i := range score {
			if score[i] > highScore {
				highScore = score[i]
			}
		}
		fmt.Printf("Part B high score %d\n", highScore)
	}
}
