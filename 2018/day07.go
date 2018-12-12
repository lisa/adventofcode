package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/kr/pretty"
	"math"
	"os"
	"sort"
	"strings"
)


var (
	partB        = flag.Bool("partB", false, "Perform part B solution?")
	inputFile    = flag.String("input", "inputs/day07.txt", "Input")
	debug        = flag.Bool("debug", false, "Debug?")
	freeAtOffset int
)

func stepCompletionTimeLookup(s string) int {
	return map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
		"D": 4,
		"E": 5,
		"F": 6,
		"G": 7,
		"H": 8,
		"I": 9,
		"J": 10,
		"K": 11,
		"L": 12,
		"M": 13,
		"N": 14,
		"O": 15,
		"P": 16,
		"Q": 17,
		"R": 18,
		"S": 19,
		"T": 20,
		"U": 21,
		"V": 22,
		"W": 23,
		"X": 24,
		"Y": 25,
		"Z": 26,
	}[s]
}

// WorkerPool - Wrapper for all of the workers
type WorkerPool struct {
	Workers []*Worker
}

// NewWorkerPool - create the worker pool of +size+
func NewWorkerPool(size int) *WorkerPool {
	wp := &WorkerPool{Workers: make([]*Worker, 0)}
	for i := 0; i < size; i++ {
		wp.AddWorker()
	}
	return wp
}

// StartStep - If possible, begin working on +step+ having begun at +beginAt+
// If there are no free workers, false will be returned, otherwise true.
func (wp *WorkerPool) StartStep(step string, beginAt int) bool {
	w := wp.GetFreeWorker(beginAt)
	if w == nil {
		if *debug {
			fmt.Printf("t=%d Got nil back from GetFreeWorker(%d)\n", beginAt, beginAt)
		}
		return false
	}
	w.FreeAt = stepCompletionTimeLookup(step) + freeAtOffset + beginAt
	w.Step = step
	if *debug {
		fmt.Printf("time %d: StartStep(%s,%d): Scheduled on worker %d; completing at %d\n", beginAt, step, beginAt, w.ID, w.FreeAt)
	}
	return true
}

// AddWorker - Adds a worker to the pool and then returns it
// This will auto-generate the ID and set the Step to "" and FreeAt to 0 (free
// now)
func (wp *WorkerPool) AddWorker() *Worker {
	w := &Worker{
		ID:     len(wp.Workers) + 1,
		Step:   "",
		FreeAt: 0,
	}
	wp.Workers = append(wp.Workers, w)
	return w
}

func (wp *WorkerPool) GetWorkerByID(id int) (*Worker, error) {
	for _, worker := range wp.Workers {
		if worker.ID == id {
			return worker, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No worker with ID %d\n", id))
}

func (wp *WorkerPool) GetWorkerByStep(step string) (*Worker, error) {
	for _, worker := range wp.Workers {
		if worker.Step == step {
			return worker, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No worker working on step %s\n", step))
}

// NextFreeAt - Return the time when the next worker is going to be free
func (wp *WorkerPool) NextFreeAt() int {
	t := math.MaxInt16
	for _, worker := range wp.Workers {
		if worker.FreeAt < t {
			t = worker.FreeAt
		}
	}
	return t
}

// HasFreeWorker - is there a free worker as of timestamp +at+?
func (wp *WorkerPool) HasFreeWorker(at int) bool {
	for _, worker := range wp.Workers {
		if *debug {
			fmt.Printf("time %d: Checking Worker %d to see if they're free (they claim free at %d) %d >= %d == %t\n", at, worker.ID, worker.FreeAt, at, worker.FreeAt, at >= worker.FreeAt)
		}
		if at >= worker.FreeAt {
			return true
		}
	}
	return false
}

// GetCompletedSteps - return a list of all steps which completed on this tick
func (wp *WorkerPool) GetCompletedSteps(at int, skipZero bool) []string {
	ret := make([]string, 0)
	// do no work if at == 0 && skipZero
	if at == 0 && skipZero {
		return []string{}
	}
	for _, worker := range wp.Workers {
		if worker.FreeAt == at {
			ret = append(ret, worker.Step)
		}
	}
	return ret
}

// AllWorkersDone - are all the workers done work as of +at+?
func (wp *WorkerPool) AllWorkersDone(at int) bool {
	finished := true
	for _, worker := range wp.Workers {

		finished = finished && (worker.FreeAt <= at)
		if *debug {
			fmt.Printf("t=%d AllWOrkersDone: worker %d free (on %s) at %d finished? %t\n", at, worker.ID, worker.Step, worker.FreeAt, finished)
		}
		if !finished {
			break
		}
	}
	return finished
}

// GetFreeWorker - Gets a free worker, if there is one. If there is no free
// worker then nil is returned instead. The caller must specify at which
// timestamp they want.
func (wp *WorkerPool) GetFreeWorker(at int) *Worker {

	for _, worker := range wp.Workers {
		if *debug {
			fmt.Printf("t=%d GetFreeWorker: worker %d is FreeAt %d\n", at, worker.ID, worker.FreeAt)
		}
		if at >= worker.FreeAt {
			return worker
		}
	}
	return nil
}

// Worker - represents an Elf that is working on a step (or me, the Human)
type Worker struct {
	ID     int    // what worker ID is this?
	Step   string // what step is this working on?
	FreeAt int    // when is this worker expected to be free?
}

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}

// SleepTiming - letter to offset: A -> 1, B -> 2, etc.
func SleepTiming(step string) int {
	return int(([]rune(strings.ToUpper(step))[0] - 64))
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
	allDeps := depmap[step]
	if *debug {
		fmt.Printf("IsStepReady: %s are all %s's deps, they must all be in processedList (%s) to schedule %s. Let's iterate and check\n", allDeps, step, processedList, step)
	}
	for i, dep := range depmap[step] {
		ready = ready && SliceContains(dep, processedList)
		if *debug {
			fmt.Printf("IsStepReady [dep %d/%d -> %s] pre-req done? %t\n", i, len(allDeps), dep, SliceContains(dep, processedList))
		}
		if !ready {
			// no need to continue
			// (that's a programming joke)
			break
		}
	}
	return ready
}

// ResolveRDepMapPartB - A reimagination for part B
/* Guidance & Note Each time through the loop represents a single second.
 * Each time we must ask:
 1. Has a worker just completed a step? If so:
	 1a. Add the step to the completed list
	 1b. Append the step to the sequence
	 1c. Its child dependencies should be added to the "ready" list (if their other deps are done)
 2. For each free worker,
	2a. Assign a single task until free workers or ready tasks is exhausted
 3. If all of the workers are idle and there are no tasks waiting, return the sequence
 4. Otherwise, increment the clock and iterate again.
*/
func ResolveRDepMapPartB(wp *WorkerPool, roots []string, rdepmap, depmap map[string][]string) (string, int) {
	t := 0
	sequence := ""
	toProcessList := make([]string, 0)
	sort.Strings(roots)
	toProcessList = append(toProcessList, roots...)
	processedList := make([]string, 0)

	for {
		if *debug {
			fmt.Println()
			fmt.Printf("Worker pool %# v\n", pretty.Formatter(wp))
			fmt.Printf("To process list: %s, processed list: %s\n", toProcessList, processedList)
		}
		completedNow := wp.GetCompletedSteps(t, true)
		// Sort, I guess?
		if *debug {
			fmt.Printf("t=%d Sequence=%s\n", t, sequence)
		}
		sort.Strings(completedNow)
		if len(completedNow) > 0 {
			// One or more workers completed a step on this tick
			if *debug {
				fmt.Printf("t=%d Completed steps %s\n", t, completedNow)
			}
			// Now we have to add any child deps from that step to the ready list, if they are ready
			for _, completedStep := range completedNow {
				// add to completed list
				if *debug {
					fmt.Printf("Finished a step!\n")
				}
				processedList = append(processedList, completedStep)
				sequence += completedStep
				for _, childDep := range rdepmap[completedStep] {

					if IsStepReady(childDep, depmap, processedList) {
						if *debug {
							fmt.Printf("We can schedule %s! Adding to the list (%s)\n", childDep, toProcessList)
						}
						toProcessList = append(toProcessList, childDep)
					} else {
						if *debug {
							fmt.Printf("We can't schedule %s now\n", childDep)
						}
					}
				}
			}
			// Dedupe and sort
			toProcessList = SliceUniqMap(toProcessList)
			sort.Strings(toProcessList)
		}
		// Can we schedule anything?
		if *debug {
			fmt.Printf("Let's check if we can schedule something (%s todo)\n", toProcessList)
		}
	tryToSchedule:
		for len(toProcessList) > 0 && wp.HasFreeWorker(t) {
			if *debug {
				fmt.Printf("We have work to do and workers to do it! Scheduling\n")
			}
			for {
				// Get the first item from the list unless there's no work left
				if len(toProcessList) == 0 {
					break tryToSchedule
				}
				todo := toProcessList[0]
				if ok := wp.StartStep(todo, t); !ok {
					if *debug {
						fmt.Printf("t=%d Couldn't schedule step %s\n", t, todo)
					}
					// if the first item in the list can't be processed for whatever reason none of
					// the subsequent items will be able either, so let's just bail out of the loop
					// now.
					break tryToSchedule
				} else {
					// scheduled okay, remove from list
					if *debug {
						fmt.Printf("We scheduled %s, so removing it from todo list %s\n", todo, toProcessList)
					}
					toProcessList = append(toProcessList[:0], toProcessList[0+1:]...)
				}
			}
			/*
				for _, todo := range toProcessList {
					if ok := wp.StartStep(todo, t); !ok {
						fmt.Printf("t=%d Couldn't schedule step %s\n", t, todo)
						break tryToSchedule
					} else {
						// it was scheduled, so remove it from the list
						fmt.Printf("We scheduled %s, so removing it from todo list %s\n", todo, toProcessList)
						toProcessList = append(toProcessList[:0], toProcessList[0+1:]...)
					}
				}
			*/
		}
		if *debug {
			fmt.Printf("end of loop. To Process list: %s, work pool %# v\n", toProcessList, pretty.Formatter(wp))
		}
		// Everything is scheduled now
		if wp.AllWorkersDone(t) && len(toProcessList) == 0 {
			if *debug {
				fmt.Printf("t=%d All done! :)\n", t)
			}
			return sequence, t
		}

		t++
	}
}

// ResolveRDepMap - for the given +rdepmap+ and starting point +roots+ return the
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
	input.Close()
	freeAtOffset = 0

	workers := 1
	if *partB {
		workers = 5
		freeAtOffset = 60
	}

	wp := NewWorkerPool(workers)

	var rootDeps []string

	// the roots are the only item in the dep map which has no corresponding key in the rdepmap
	for dep, _ := range rdepmap {
		if _, ok := depmap[dep]; !ok {
			// 'dep' is the root
			rootDeps = append(rootDeps, dep)
		}
	}
	if !*partB {
		fmt.Printf("RDEP PATH: %s\n", ResolveRDepMap(rootDeps, rdepmap, depmap))
	} else {
		sequence, worktime := ResolveRDepMapPartB(wp, rootDeps, rdepmap, depmap)
		fmt.Printf("Got %s back in %d seconds\n", sequence, worktime)
	}
}
