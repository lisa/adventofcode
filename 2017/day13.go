package main

/*
Day 13, Part A

Given a list of inputs in a file in the format `layer number: depth`, each on
a new line, which represent the number of firewall layers in position. If
there is no layer specified, it has no depth.

Each cycle a "scanner" will progress down in depth of each layer. When the
"scan" reaches the bottom, start over at the top.

After the scanner moves, progress through each layer to the end (at the top)
and record each layer where the scanner runs into us.

Move the position, then move the scanner. If we run into an active scan then
we're detected.

*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "./inputs/day13-example.txt", "Input File")
var debug = flag.Bool("debug", false, "Debug?")

type Firewall struct {
	Rules             map[int]int  // layer # -> depth
	Positions         map[int]int  // current position in each layer
	MovementDirection map[int]bool // true=down, false=up
}

func (fw *Firewall) AddRuleAtPos(layer, depth int) {
	fw.Rules[layer] = depth
	fw.Positions[layer] = 0
	fw.MovementDirection[layer] = true
}

// how many layers in total?
func (fw *Firewall) Len() int {
	return len(fw.Rules)
}

/*
Advance the scanner once. Return a map of what the current position is for
each layer. If there is a gap between layers there will not be an entry.
*/
func (fw *Firewall) Advance() map[int]int {
	for layerNumber, _ := range fw.Positions {
		// Skip if there's no rules, or it has 0 length.
		if fw.Rules[layerNumber] == 0 {
			continue
		}
		/*
			If we're previously going down:
			 if i can go down, go down. otherwise, flip direction; go up
			If we're previously going up:
			 if i can go up, go up. otherwise, flip direction; go down
		*/

		if fw.MovementDirection[layerNumber] {
			// Going down
			if fw.Positions[layerNumber]+1 >= fw.Rules[layerNumber] {
				// can't keep going down, so flip around and go up
				fw.MovementDirection[layerNumber] = false
				fw.Positions[layerNumber] -= 1
			} else {
				// keep going down
				fw.Positions[layerNumber] += 1
			}
		} else {
			// Going up
			if fw.Positions[layerNumber]-1 < 0 {
				// can't keep going up, so flip around and go down
				fw.MovementDirection[layerNumber] = true
				fw.Positions[layerNumber] += 1
			} else {
				// keep going up
				fw.Positions[layerNumber] -= 1
			}
		} // end: which way am i going?
	} // out of layers
	return fw.Positions
}

func (fw *Firewall) HighestLayer() int {
	// highest layer number
	highest := -1
	for layerNumber, _ := range fw.Rules {
		if layerNumber > highest {
			highest = layerNumber
		}
	}
	return highest
}

// It's possible to have gaps, so fill them in with 0 depth layers
func (fw *Firewall) FillInGaps() {
	for i := 0; i < fw.HighestLayer(); i++ {
		if fw.Rules[i] == 0 {
			fw.Rules[i] = 0
			fw.MovementDirection[i] = true
		}
	}
}

// true if the position of the scan is at 0 in layer layerPosition.
func (fw *Firewall) CheckCollision(layerPosition int) bool {
	return fw.Rules[layerPosition] > 0 && fw.Positions[layerPosition] == 0
}

// What's the cost of being caught in layer layerPosition?
func (fw *Firewall) CollisionCost(layerPosition int) int {
	return layerPosition * fw.Rules[layerPosition]
}

func NewFirewall() *Firewall {
	return &Firewall{
		Rules:             make(map[int]int),
		Positions:         make(map[int]int),
		MovementDirection: make(map[int]bool),
	}
}

func main() {
	flag.Parse()
	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't read file: %s\n", err)
		os.Exit(1)
	}
	defer input.Close()
	firewall := NewFirewall()

	lineReader := bufio.NewScanner(input)
	for lineReader.Scan() {
		line := lineReader.Text()
		var layer, depth int
		for n, token := range strings.Split(line, ":") {
			switch n {
			case 0:
				layer, err = strconv.Atoi(strings.Trim(token, " "))
				if err != nil {
					fmt.Printf("Couldn't convert %s to layer number.\n", token)
					os.Exit(1)
				}
			case 1:
				depth, err = strconv.Atoi(strings.Trim(token, " "))
				if err != nil {
					fmt.Printf("Couldn't convert %s to depth.\n", token)
					os.Exit(1)
				}
			default:
				fmt.Printf("Unknown item found at %d (%s) in %s\n", n, token, line)
				os.Exit(1)
			}
		} // EOL
		firewall.AddRuleAtPos(layer, depth)
	} // EOF
	firewall.FillInGaps()
	if *debug {
		fmt.Printf("Firewall after creation: %+v\n", firewall)
		fmt.Println()
	}
	collisionCost := 0
	// check initial condition
	if firewall.CheckCollision(0) {
		collisionCost += firewall.CollisionCost(0)
	}
	if *debug {
		fmt.Printf("picosecond %d\n", 0)
		fmt.Printf("Firewall: %+v\n", firewall)
		fmt.Printf("Collision at %d? => %t\n", 0, firewall.CheckCollision(0))
		fmt.Println()
	}
	for position := 1; position <= firewall.HighestLayer(); position++ {
		firewall.Advance()
		if *debug {
			fmt.Printf("picosecond %d\n", position)
			fmt.Printf("Firewall: %+v\n", firewall)
			fmt.Printf("Collision at %d? => %t\n", position, firewall.CheckCollision(position))
		}
		if firewall.CheckCollision(position) {
			collisionCost += firewall.CollisionCost(position)
		}
		if *debug {
			fmt.Println()
		}
	}
	fmt.Printf("Made it! But at what cost...? Collision Cost: %d\n", collisionCost)
}
