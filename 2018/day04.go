package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

var (
	partB     = flag.Bool("partB", false, "Perform part B solution?")
	inputFile = flag.String("input", "inputs/day04.txt", "Input")
	debug     = flag.Bool("debug", false, "Debug?")
	debug2    = flag.Bool("debug2", false, "more debug?")
	debug3    = flag.Bool("debug3", false, "even more debug?")

	// Match groups: 1-5: yyyy-mm-dd hh:mm. 6: action identifier part 1. 7: action identifier part 2
	// action identifier part 1 could be in [Guard, falls, wakes]. If part 1 is
	// Guard, then part 2 is that Guard's identifier, otherwise, it is discarded
	// Note, group 7 doesn't actually include the `#` before the guard ID, although
	// it is an optional part of the matchgroup (eg, this token may be [#\d{1,},
	// falls, wakes])
	datePicker = regexp.MustCompile(`\[(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2})\] (\w+) #?([a-z0-9]+)`)
)

// What is the guard doing?
type guardState int

const (
	beginShift guardState = iota
	beginSleep
	wakesUp
)

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}

//GuardAction - What's the guard doing?
// PreviousAction and NextAction have to do with the next in the log file, not this guard's actions
type GuardAction struct {
	GuardID   int
	Action    guardState
	Timestamp int64
}

// Guard - just this guard, you know?
type Guard struct {
	GuardID      int
	SleepTime    int   // total amount of time asleep
	SleepMinutes []int // on which minute do they fall asleep?
}

func (g Guard) AddSleepTime(t int) Guard {
	g.SleepTime += t
	return g
}

// add a single minute
func (g Guard) AddSleepMinutes(t int) Guard {
	g.SleepMinutes = append(g.SleepMinutes, t)
	return g
}

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't open %s: %v\n", *inputFile, err)
		os.Exit(1)
	}
	defer input.Close()
	lineReader := bufio.NewScanner(input)

	guardActions := make([]GuardAction, 0)
	// Parse the log file:
	// Pick out the date and guard info.
	// We'll maintain a sorted list of all the timestamps (`timestamps`) that we'll
	// use to add guard actions to a linked list
	for lineReader.Scan() {
		line := lineReader.Text()

		lineParts := datePicker.FindAllStringSubmatch(line, -1)[0]
		year, err := strconv.Atoi(lineParts[1])
		errorIf("Couldn't parse the year", err)
		month, err := strconv.Atoi(lineParts[2])
		errorIf("Couldn't parse the month", err)
		day, err := strconv.Atoi(lineParts[3])
		errorIf("Couldn't parse the day", err)
		hour, err := strconv.Atoi(lineParts[4])
		errorIf("Couldn't parse the hour", err)
		minute, err := strconv.Atoi(lineParts[5])
		errorIf("Couldn't parse the minute", err)

		actionTime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC).Unix()

		var guard GuardAction

		guard.Timestamp = actionTime

		// parse what it's doing
		switch lineParts[6] {
		case "Guard":
			guardID, err := strconv.Atoi(lineParts[7])
			errorIf("Got a Guard, but couldn't parse its ID", err)
			guard.GuardID = guardID
			guard.Action = beginShift
		case "falls":
			//guard fell asleep
			guard.Action = beginSleep
		case "wakes":
			//guard woke up
			guard.Action = wakesUp
		default:
			// wut?
			fmt.Printf("Couldn't figure out WTF this guard was doing: %s\n", lineParts[6])
			os.Exit(1)
		}
		// Add to the list and then sort right away
		guardActions = append(guardActions, guard)
		sort.Slice(guardActions, func(i, j int) bool { return guardActions[i].Timestamp < guardActions[j].Timestamp })
	}

	// At this point guardActions may have GuardAction "objects" which lack
	// GuardIDs, so we need to be mindful going forward to add them.
	// But it is a sorted list

	if *debug3 {
		fmt.Printf("Some are missing guard IDs\n")
		for index, ga := range guardActions {
			fmt.Printf("Index %d -> (Timestamp: [%d-%02d-%02d %02d:%02d]) = %+v\n", index,
				time.Unix(ga.Timestamp, 0).Year(), time.Unix(ga.Timestamp, 0).Month(),
				time.Unix(ga.Timestamp, 0).Day(), time.Unix(ga.Timestamp, 0).Hour(),
				time.Unix(ga.Timestamp, 0).Minute(),
				ga)
		}
	}

	var currentGuard int

	// Summary
	var sleepBeginsAt time.Time
	// guard id -> sleepcount
	allGuards := make(map[int]Guard)
	mostSleptID := -1
	mostSleptMinutes := -1
	for index, guard := range guardActions {
		timestamp := time.Unix(guard.Timestamp, 0)
		switch guard.Action {
		case beginShift:
			currentGuard = guard.GuardID
			if *debug {
				fmt.Printf("[%d-%02d-%02d %02d:%02d] Guard %d began shift\n",
					timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), currentGuard)
			}
			if _, ok := allGuards[currentGuard]; !ok {
				if *debug {
					fmt.Printf("init guard %d\n", currentGuard)
				}
				allGuards[currentGuard] = Guard{
					GuardID:      currentGuard,
					SleepTime:    0,
					SleepMinutes: make([]int, 0),
				}
			}
		case beginSleep:
			guardActions[index].GuardID = currentGuard
			if *debug {
				fmt.Printf("[%d-%02d-%02d %02d:%02d] Guard %d went to sleep\n",
					timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), currentGuard)
			}
			sleepBeginsAt = timestamp
			allGuards[currentGuard] = allGuards[currentGuard].AddSleepMinutes(timestamp.Minute())
		case wakesUp:
			guardActions[index].GuardID = currentGuard
			if *debug {
				fmt.Printf("[%d-%02d-%02d %02d:%02d] Guard %d woke up\n",
					timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), currentGuard)
			}
			// tally how long were they asleep?
			sleepDuration := int(timestamp.Sub(sleepBeginsAt).Minutes())
			allGuards[currentGuard] = allGuards[currentGuard].AddSleepTime(sleepDuration)

			// add all the minutes they were asleep to the list.
			// take each minute in the range (now - start] and add to the list
			for i := sleepBeginsAt.Unix(); i < timestamp.Unix(); i += 60 {
				if *debug {
					fmt.Printf(" * Adding minute %d\n", time.Unix(i, 0).Minute())
				}
				allGuards[currentGuard] = allGuards[currentGuard].AddSleepMinutes(time.Unix(i, 0).Minute())
			}

			if allGuards[currentGuard].SleepTime > mostSleptMinutes {
				mostSleptMinutes = allGuards[currentGuard].SleepTime
				mostSleptID = currentGuard
			}
		}

	}

	if *debug2 {
		fmt.Printf("Should now all have guardIDs that make sense\n")
		for index, ga := range guardActions {
			fmt.Printf("Index %d -> (Timestamp: [%d-%02d-%02d %02d:%02d]) = %+v\n", index,
				time.Unix(ga.Timestamp, 0).Year(), time.Unix(ga.Timestamp, 0).Month(),
				time.Unix(ga.Timestamp, 0).Day(), time.Unix(ga.Timestamp, 0).Hour(),
				time.Unix(ga.Timestamp, 0).Minute(),
				ga)
		}
	}
	if *debug {
		fmt.Printf("Most sleepy guard: %d. Their actions: %+v\n", mostSleptID, allGuards[mostSleptID])
	}
	// histogram: minute -> count
	sleepMinuteHistogram := make(map[int]int)
	for _, minute := range allGuards[mostSleptID].SleepMinutes {
		if *debug {
			fmt.Printf("Incrementing minute %d count to %d\n", minute, sleepMinuteHistogram[minute]+1)
		}
		sleepMinuteHistogram[minute]++
	}
	if *debug {
		fmt.Printf("Init most common minute to %d\n", allGuards[mostSleptID].SleepMinutes[0])
	}
	mostCommonMinute := allGuards[mostSleptID].SleepMinutes[0]
	for minute, count := range sleepMinuteHistogram {
		if count > sleepMinuteHistogram[mostCommonMinute] {
			if *debug {
				fmt.Printf("new high score: %d minutes at %d\n", minute, count)
			}
			mostCommonMinute = minute
		}
	}
	fmt.Printf("Guard %d slept the most for a total of %d minutes. They slept most on minute %d, which puts the math at %d * %d = %d\n",
		mostSleptID, allGuards[mostSleptID].SleepTime, mostCommonMinute, mostSleptID, mostCommonMinute, mostSleptID*mostCommonMinute)
}
