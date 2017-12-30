package main

/* Day 9 part A:

For a given input stream consisting of groups of characters bounded by {}
and optionally within {} separated by a comma, determine the total number of
groups. These may be infinitely nested.

There are special groups ("garbage") bounded by < and > within which are
excluded from the overall count. The ! character will "cancel" the following
character, including another !. Thus, !! is a noop. Garbage only occurs within
groups (bounded by {}).

Parsing State Machine:

! skip next
< begin garbage; ignore everything until >
{ begin group
, ok to start a new group when seeing another {.
} end of previous group

Scoring:
Each group is assigned a score which is one more than the score of the group
that immediately contains it. (The outermost group gets a score of 1.)

Part B:
Count the number of characters in garbage. <> don't count towards the count and
neither does ! nor the character(s) being cancelled by !.
*/

import (
	"flag"
	"fmt"
	"strings"
)

var input = flag.String("input", "{}", "Input string for the program")
var partB = flag.Bool("partB", false, "Perform part B solution")
var debug = flag.Bool("debug", false, "Debug output")

func main() {
	flag.Parse()
	// Loop over each character
	characters := strings.Split(*input, "")
	depth := 0             // Group depth; used with scoring
	score := 0             // Total score
	garbageCharacters := 0 // Cleaned up garbage characters
	ignoreNext := false
	ignoreGroups := false
	for i, token := range characters {
		if *debug {
			fmt.Printf("[%d/%d] score: %d, depth: %d, ignoreNext: %t, ignoreGroups: %t, token %s\n", i, len(*input)-1, score, depth, ignoreNext, ignoreGroups, token)
		}
		if ignoreNext {
			// Clear and skip this character
			ignoreNext = false
			continue
		}
		switch token {
		case "{":
			if !ignoreGroups {
				depth += 1
			} else {
				garbageCharacters += 1
			}
		case "<":
			if ignoreGroups {
				// Already ignoring, so we need to bump garbageCharacters
				garbageCharacters += 1
			}
			ignoreGroups = true // start of garbage
		case ">":
			ignoreGroups = false // end of garbage
		case "}":
			if !ignoreGroups {
				score += depth
				depth -= 1
			} else {
				garbageCharacters += 1
			}
		case "!":
			ignoreNext = true
		default:
			if ignoreGroups {
				garbageCharacters += 1
			}
		}
	} // EOF

	if *partB {
		fmt.Printf("Garbage characters: %d\n", garbageCharacters)
	} else {
		fmt.Printf("high score: %d\n", score)
	}
}
