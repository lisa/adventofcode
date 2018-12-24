package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

var (
	inputFile    = flag.String("input", "inputs/day12.txt", "input file")
	partB        = flag.Bool("partB", false, "do part b solution?")
	debug        = flag.Bool("debug", false, "debug?")
	inputMatcher = regexp.MustCompile(`(?m:^(initial state: (.*))|^(.*) => (.*)$)`)
)

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}

// Rule matches: five states, each of which may be t/f
// result is t for flower result, f for no flower
type Rule struct {
	pattern []bool // pattern to match: LLCRR
	result  bool   // result of pattern match
}

func (r *Rule) String() string {
	ret := ""

	for _, v := range r.pattern {
		if v {
			ret += "#"
		} else {
			ret += "."
		}
	}
	ret += " => "
	if r.result {
		ret += "#"
	} else {
		ret += "."
	}
	return ret
}

// Apply - apply this rule to the given pot with its neighbours
// +potWithNeighbours+ is the Current pot (`C`) with two pots on either side
// `LL` to the left, `RR` to the right, thus, LLCRR
// -1 is false, 0 is the given pot didn't match this rule (eg pass), 1 is true
func (r *Rule) Apply(potWithNeighbours []bool) int {
	matches := true
	if *debug {
		fmt.Printf("Applying rule %s (res=%t) to %s", r, r.result, boolSliceToString(potWithNeighbours))
	}
	for i := 0; i <= 4 && matches; i++ {
		matches = matches && (r.pattern[i] == potWithNeighbours[i])
	}

	if matches {
		switch r.result {
		case true:
			if *debug {
				fmt.Printf(" => #\n")
			}
			return 1
		case false:
			if *debug {
				fmt.Printf(" => .\n")
			}
			return -1
		}
	}
	if *debug {
		fmt.Printf(" => P\n")
	}
	return 0
}

// NewRule - return a rule
func NewRule(r []bool, result bool) *Rule {
	return &Rule{
		pattern: r,
		result:  result,
	}
}

// Generation - do a generation
func Generation(state map[int]bool, rules *Rule) map[int]bool {
	ret := make(map[int]bool)

	//for pot := range state {

	//	}
	return ret
}

func getMinMaxPot(pots *map[int]bool) (int, int) {
	var min, max int

	for p, _ := range *pots {
		if p > max {
			max = p
		}
		if p < min {
			min = p
		}
	}
	return min, max
}

func boolSliceToString(pots []bool) string {
	str := ""

	for i := range pots {
		if pots[i] {
			str += "#"
		} else {
			str += "."
		}
	}
	return str
}

func main() {
	flag.Parse()

	inputBuffer, err := ioutil.ReadFile(*inputFile)
	errorIf("couldn't read input file", err)

	// pot number => has a flower?
	state := make(map[int]bool)
	rules := make([]*Rule, 0)
	if *debug {
		fmt.Printf("Input = %s\n", string(inputBuffer))
	}

	matches := inputMatcher.FindAllStringSubmatch(string(inputBuffer), -1)
	minPotID := 0
	maxPotID := 0
	for i := 0; i < len(matches[0][2]); i++ {
		// hash (plant)
		switch matches[0][2][i] {
		case 35:
			state[i] = true
		case 46:
			// no need to set false, it is the default
		default:
			fmt.Printf("No idea what %s is in the initial state!", string(matches[0][2][i]))
			os.Exit(1)
		}
	}
	minPotID, maxPotID = getMinMaxPot(&state)

	for r := 1; r < len(matches); r++ {
		rule := make([]bool, 5)
		var ruleResult bool
		for i := 0; i < 5; i++ {
			switch matches[r][3][i] {
			case 35:
				rule[i] = true
			case 46:
				rule[i] = false
			default:
				fmt.Printf("no idea what %s is in rule %d\n", string(matches[r][3][i]), r-1)
				os.Exit(1)
			}
		}
		switch matches[r][4][0] {
		case 35:
			ruleResult = true
		case 46:
			ruleResult = false
		default:
			fmt.Printf("No idea what %s means in rule result for rule %d\n", string(matches[r][4][0]), r-1)
			os.Exit(1)
		}
		rules = append(rules, NewRule(rule, ruleResult))
	}
	if *debug {
		fmt.Printf("Initial State: %v\n", state)
		for _, r := range rules {
			fmt.Printf("Rule Pattern => %t. String: %s\n", r.pattern, r)
		}
	}

	// this wants to create the next generation based on the old generation.
	// loop through current, apply rules to pots [n-2..n+2] and the result of that becomes the next gen
	// will need to maintain the pot ID of the extrema to avoid searching the entire keyspace
	for g := 1; g <= 20; g++ {
		if *debug {
			fmt.Printf("[%d/20] State %+v\n", g, state)
		}

		nextGenState := make(map[int]bool)
		for p := minPotID - 2; p < maxPotID+2; p++ {
			var hasPot int
			neighbourPots := make([]bool, 5)

			neighbourPots[2] = state[p]
			neighbourPots[1] = state[p-1]
			neighbourPots[0] = state[p-2]
			neighbourPots[3] = state[p+1]
			neighbourPots[4] = state[p+2]
			if *debug {
				fmt.Printf("pot %d current state => %t\n", p, state[p])
			}
			for _, rule := range rules {
				hasPot = rule.Apply(neighbourPots)

				switch hasPot {
				case -1:
					// next generation state for this is false
					if *debug {
						fmt.Printf(" next gen => false\n")
					}
					break
				case 1:
					// have a plant
					if *debug {
						fmt.Printf(" next gen => true\n")
					}
					nextGenState[p] = true
					break
				case 0:
					// didn't match this rule
					continue
				}
			}

		}
		if *debug {
			fmt.Printf("[%d/20] Next Gen State = %+v\n", g, nextGenState)
		}
		state = nextGenState
		minPotID, maxPotID = getMinMaxPot(&state)
	}
	sum := 0
	for p := range state {
		sum += p
	}

	fmt.Printf("Pot sum = %d\n", sum)

}
