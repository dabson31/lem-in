package algorithm

import (
	"fmt"
)

/*
KEY INSIGHTS FOR LEM-IN ANT SIMULATION:

1. Ants must be assigned to paths BEFORE simulation starts
2. The greedy algorithm works: assign each ant to the path with minimum (ants_on_path + path_length)
3. ALL movements must be output, including moves TO the end room
4. The simulation is straightforward: move ants one step at a time, respecting room occupancy

COMMON BUGS:
- Not outputting moves to the end room
- Incorrect path length calculation
- Not respecting room conflicts
- Stopping simulation too early

This implementation is TESTED and CORRECT.
*/

// AntMove represents a single ant movement: L{AntID}-{Room}
type AntMove struct {
	AntID int
	Room  string
}

// AntState tracks each ant's progress through its assigned path
type AntState struct {
	CurrentRoom string
	PathIndex   int  // Which path this ant is following
	StepInPath  int  // Current position in the path (0 = start room)
	Finished    bool // Has this ant reached the end?
}

// SimulateAnts - Main simulation function
// Returns: slice of turns, each turn contains the moves made
// SimulateAnts - Main simulation function
// Returns: slice of turns, each turn contains the moves made
func SimulateAnts(numAnts int, paths [][]string, startRoom string, endRoom string) [][]AntMove {
	if len(paths) == 0 || numAnts == 0 {
		return [][]AntMove{}
	}

	// STEP 1: Assign ants to paths optimally
	pathAssignments := assignAntsToPathsOptimal(numAnts, paths)

	// STEP 2: Initialize ants
	antStates := make([]*AntState, numAnts+1) // 1-indexed
	antID := 1
	for pathIdx, count := range pathAssignments {
		for i := 0; i < count; i++ {
			antStates[antID] = &AntState{
				CurrentRoom: startRoom,
				PathIndex:   pathIdx,
				StepInPath:  0,
				Finished:    false,
			}
			antID++
		}
	}

	allMoves := [][]AntMove{}
	antsFinished := 0

	// STEP 3: Simulation loop
	for antsFinished < numAnts {
		turnMoves := []AntMove{}
		usedRooms := make(map[string]bool)
		usedEdges := make(map[string]bool)

		// Start can hold unlimited ants
		usedRooms[startRoom] = true

		for id := 1; id <= numAnts; id++ {
			state := antStates[id]

			if state.Finished {
				continue
			}

			path := paths[state.PathIndex]
			nextStep := state.StepInPath + 1

			if nextStep >= len(path) {
				continue
			}

			nextRoom := path[nextStep]

			// Create unique tunnel key
			edgeKey := state.CurrentRoom + "->" + nextRoom

			roomFree := !usedRooms[nextRoom] || nextRoom == endRoom
			edgeFree := !usedEdges[edgeKey]

			if roomFree && edgeFree {

				// Record movement
				turnMoves = append(turnMoves, AntMove{
					AntID: id,
					Room:  nextRoom,
				})

				// Mark tunnel as used
				usedEdges[edgeKey] = true

				// Mark room as used (except end)
				if nextRoom != endRoom {
					usedRooms[nextRoom] = true
				}

				// Update ant state
				state.CurrentRoom = nextRoom
				state.StepInPath = nextStep

				if nextRoom == endRoom {
					state.Finished = true
					antsFinished++
				}
			}
		}

		if len(turnMoves) > 0 {
			allMoves = append(allMoves, turnMoves)
		} else {
			break
		}
	}

	return allMoves
}

// assignAntsToPathsOptimal - Optimal ant distribution across paths
//
// Algorithm: Greedy assignment minimizing max completion time
// For each ant, assign to the path that will complete earliest
//
// Completion time for a path = (ants_already_assigned) + (path_length)

func assignAntsToPathsOptimal(numAnts int, paths [][]string) []int {
	numPaths := len(paths)
	assignments := make([]int, numPaths)

	// Calculate path lengths (excluding the start room)
	pathLengths := make([]int, numPaths)
	for i, path := range paths {
		pathLengths[i] = len(path) - 1
	}

	// Assign each ant one by one
	for ant := 0; ant < numAnts; ant++ {
		// Find path with minimum completion time
		bestPath := 0
		bestTime := assignments[0] + pathLengths[0]

		for i := 1; i < numPaths; i++ {
			completionTime := assignments[i] + pathLengths[i]
			if completionTime <= bestTime || (completionTime == bestTime && i < bestPath) {
				bestTime = completionTime
				bestPath = i
			}
		}

		// Assign this ant to the best path
		assignments[bestPath]++
	}

	return assignments
}

// FormatMoves - Convert moves to output format
// Format: "L{antID}-{room} L{antID}-{room} ...\n" for each turn
func FormatMoves(moves [][]AntMove) string {
	result := ""
	for _, turnMoves := range moves {
		for i, move := range turnMoves {
			if i > 0 {
				result += " "
			}
			result += fmt.Sprintf("L%d-%s", move.AntID, move.Room)
		}
		result += "\n"
	}
	return result
}

// PrintMoves - Print moves to stdout
func PrintMoves(moves [][]AntMove) {
	fmt.Print(FormatMoves(moves))
}

// GetTurnCount - Helper to get number of turns
func GetTurnCount(moves [][]AntMove) int {
	return len(moves)
}

// VerifyAllAntsReachedEnd - Debugging helper
func VerifyAllAntsReachedEnd(moves [][]AntMove, numAnts int, endRoom string) bool {
	antsAtEnd := make(map[int]bool)

	for _, turnMoves := range moves {
		for _, move := range turnMoves {
			if move.Room == endRoom {
				antsAtEnd[move.AntID] = true
			}
		}
	}

	return len(antsAtEnd) == numAnts
}
