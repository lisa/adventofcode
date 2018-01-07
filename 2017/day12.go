package main

/*
Day 12, part A

Find the number of pipes connected to program 0 with input looking like this:

0 <-> 2
1 <-> 1
2 <-> 0, 3, 4
3 <-> 2, 4
4 <-> 2, 3, 6
5 <-> 6
6 <-> 4, 5

The input describes a bidirectional set of pipes between programs :

0 to/from 2
2 to/from 0, 3, 4
etc
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "./inputs/day12-example.txt", "Instructions Input File")
var partB = flag.Bool("partB", false, "Perform part B solution?")
var debug = flag.Bool("debug", false, "Debug")

type PipeMap map[int][]int

// How many pipes are in the group pid? ignore pid `ignore` (ie, 0)
func (p PipeMap) ProgramsOf(pid int, seenPids map[int]bool) {
	if *debug {
		fmt.Printf("ProgramsOf: pid=%d map=%v\n", pid, seenPids)
	}
	if !seenPids[pid] {
		seenPids[pid] = true
		for _, pipePids := range p[pid] {
			p.ProgramsOf(pipePids, seenPids)
		}
	}
	// base case
}

func main() {
	flag.Parse()
	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't read file: %s\n", err)
		os.Exit(1)
	}
	defer input.Close()
	lineReader := bufio.NewScanner(input)

	programs := make(PipeMap)

	for lineReader.Scan() {
		line := lineReader.Text()
		var program int
		pipes := make([]int, 0)
		for n, token := range strings.Split(line, " ") {

			switch n {
			case 0:
				program, err = strconv.Atoi(token)
				if err != nil {
					fmt.Printf("Couldn't convert %s to a program id: %s\n", token, err)
					os.Exit(1)
				}
			case 1:
				continue // <->
			default:
				endpoint, err := strconv.Atoi(strings.TrimRight(token, ", "))
				if err != nil {
					fmt.Printf("Couldn't convert %s to a pipe endpoint: %s\n", token, err)
					os.Exit(1)
				}
				pipes = append(pipes, endpoint)
			}
		} // EOL
		programs[program] = pipes
	} // EOF
	if *debug {
		fmt.Printf("programs: %v\n", programs)
	}

	seenPids := make(map[int]bool)
	programs.ProgramsOf(0, seenPids)
	fmt.Printf("Total number of programs in the group that contain program ID %d: %d", 0, len(seenPids))
	if *debug {
		fmt.Printf(" seenPids=%v\n", seenPids)
	} else {
		fmt.Printf("\n")
	}
}
