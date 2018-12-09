package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	partB     = flag.Bool("partB", false, "Perform part B solution?")
	inputFile = flag.String("input", "inputs/day07.txt", "Input")
	debug     = flag.Bool("debug", false, "Debug?")
)

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}

// SliceUniqueMap - cribbed from https://www.reddit.com/r/golang/comments/5ia523/idiomatic_way_to_remove_duplicates_in_a_slice/db6qa2e/
func SliceUniqMap(s []string) []string {
	seen := make(map[string]struct{}, len(s))
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

// SliceContains - does the +needle+ exist in the +haystack+?
func SliceContains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

func IsStepReady(step string, depmap map[string][]string, processedList []string) bool {
	ready := true
	for _, dep := range depmap[step] {
		if *debug {
			fmt.Printf("checking %s to see if %s is present therein\n", processedList, dep)
		}
		ready = ready && SliceContains(dep, processedList)
		if !ready {
			// no need to continue
			// (that's a programming joke)
			break
		}
	}
	return ready
}

// ResolveRDepMap - for the given +rdepmap+ and starting point +root+ return the
// sequence to get from +root+ to the end of the tree, with no more dependencies
func ResolveRDepMap(roots []string, rdepmap, depmap map[string][]string) string {
	sequence := ""
	toProcessList := make([]string, 0)
	sort.Strings(roots)
	toProcessList = append(toProcessList, roots...)
	processedList := make([]string, 0)
	workingOn := toProcessList[0]

	if *debug {
		fmt.Printf("ResolveRDepMap before loop toProcessList: %s\n", toProcessList)
	}
	for {
		sequence += toProcessList[0]
		// Remove the current dep from the list
		workingOn = toProcessList[0]
		processedList = append(processedList, workingOn)
		toProcessList = append(toProcessList[:0], toProcessList[0+1:]...)

		if *debug {
			fmt.Printf(" * loop working on: %s.  Full list=%s\n", workingOn, toProcessList)
			fmt.Printf(" * %s's rdeps: %s\n", workingOn, rdepmap[workingOn])
		}
		for _, childDep := range rdepmap[workingOn] {
			if *debug {
				fmt.Printf("Have all of %s's deps been met?\n", childDep)
			}
			if IsStepReady(childDep, depmap, processedList) {
				toProcessList = append(toProcessList, childDep)
			}
		}
		// deduplicate the list (ya never know)
		toProcessList = SliceUniqMap(toProcessList)
		// sort it
		sort.Strings(toProcessList)
		if len(toProcessList) == 0 {
			break
		}
	}
	return sequence
}

func main() {
	flag.Parse()
	input, err := os.Open(*inputFile)
	errorIf("Can't open input file", err)

	defer input.Close()
	lineReader := bufio.NewScanner(input)

	depmap := make(map[string][]string)
	// backwards looking
	rdepmap := make(map[string][]string)

	for lineReader.Scan() {
		// get our steps, convert to integers (for easier sorting later on)
		tokens := strings.Split(lineReader.Text(), " ")
		rdep := tokens[1] // this step must be done
		dep := tokens[7]  // before this step

		depmap[dep] = append(depmap[dep], rdep)
		rdepmap[rdep] = append(rdepmap[rdep], dep)
	}

	var rootDeps []string

	// the roots are the only item in the dep map which has no corresponding key in the rdepmap
	for dep, _ := range rdepmap {
		if _, ok := depmap[dep]; !ok {
			// 'dep' is the root
			rootDeps = append(rootDeps, dep)
		}
	}
	fmt.Printf("RDEP PATH: %s\n", ResolveRDepMap(rootDeps, rdepmap, depmap))
}
