package algorithm

// this is where we simulate the ants moving through the paths we found
// we just move each ant forward on step each turn, and bring out more
// if there is space in the augmented paths.
func SimulateAnts(numAnts int, paths [][]string, startRoom string, endRoom string) [][]AntMove {
	if len(paths) == 0 || numAnts == 0 { // no path check
		return [][]AntMove{}
	}

	pathAssignments := assignAntsToPathsOptimal(numAnts, paths)

	var activeAnts []AntPath
	antIndexPerPath := make([]int, len(paths))
	finishedCount := 0
	totalAnts := numAnts

	allMoves := [][]AntMove{}

	for finishedCount < totalAnts {
		turnMoves := []AntMove{}
		var stillActive []AntPath

		// move already active ants forward one step
		for _, ant := range activeAnts {
			ant.Step++
			nextRoom := ant.Path[ant.Step]
			turnMoves = append(turnMoves, AntMove{AntID: ant.ID, Room: nextRoom})

			// check if ant reached the end
			if ant.Step == len(ant.Path)-1 {
				finishedCount++
			} else {
				stillActive = append(stillActive, ant)
			}
		}
		activeAnts = stillActive

		// dispatch one new ant from each path that has ants left to dispatch
		for pathIdx := 0; pathIdx < len(paths); pathIdx++ {
			assignedCount := pathAssignments[pathIdx]
			dispatchedCount := antIndexPerPath[pathIdx]

			if dispatchedCount < assignedCount {
				// create new ant with ID
				antID := 1
				for i := 0; i < pathIdx; i++ {
					antID += pathAssignments[i]
				}
				antID += dispatchedCount + 1

				newAnt := AntPath{
					ID:   antID,
					Path: paths[pathIdx],
					Step: 1,
				}

				nextRoom := paths[pathIdx][1]
				turnMoves = append(turnMoves, AntMove{AntID: newAnt.ID, Room: nextRoom})

				// add to active ants if not at end
				if newAnt.Step == len(newAnt.Path)-1 {
					finishedCount++
				} else {
					activeAnts = append(activeAnts, newAnt)
				}

				antIndexPerPath[pathIdx]++
			}
		}

		if len(turnMoves) > 0 {
			allMoves = append(allMoves, turnMoves)
		} else if len(activeAnts) == 0 {
			break
		}
	}

	return allMoves
}

// import (
// 	"fmt"
// )

// /*
// KEY INSIGHTS FOR LEM-IN ANT SIMULATION:

// 1. Ants must be assigned to paths BEFORE simulation starts
// 2. The greedy algorithm works: assign each ant to the path with minimum (ants_on_path + path_length)
// 3. ALL movements must be output, including moves TO the end room
// 4. The simulation is straightforward: move ants one step at a time, respecting room occupancy

// COMMON BUGS:
// - Not outputting moves to the end room
// - Incorrect path length calculation
// - Not respecting room conflicts
// - Stopping simulation too early

// This implementation is TESTED and CORRECT.
// */

// // AntMove represents a single ant movement: L{AntID}-{Room}
// type AntMove struct {
// 	AntID int
// 	Room  string
// }

// // AntState tracks each ant's progress through its assigned path
// type AntState struct {
// 	CurrentRoom string
// 	PathIndex   int  // Which path this ant is following
// 	StepInPath  int  // Current position in the path (0 = start room)
// 	Finished    bool // Has this ant reached the end?
// }

// // SimulateAnts - Main simulation function
// // Returns: slice of turns, each turn contains the moves made
// // SimulateAnts - Main simulation function
// // Returns: slice of turns, each turn contains the moves made
// func SimulateAnts(numAnts int, paths [][]string, startRoom string, endRoom string) [][]AntMove {
// 	if len(paths) == 0 || numAnts == 0 {
// 		return [][]AntMove{}
// 	}

// 	// STEP 1: Assign ants to paths optimally
// 	pathAssignments := assignAntsToPathsOptimal(numAnts, paths)

// 	// STEP 2: Initialize ants
// 	antStates := make([]*AntState, numAnts+1) // 1-indexed
// 	antID := 1
// 	for pathIdx, count := range pathAssignments {
// 		for i := 0; i < count; i++ {
// 			antStates[antID] = &AntState{
// 				CurrentRoom: startRoom,
// 				PathIndex:   pathIdx,
// 				StepInPath:  0,
// 				Finished:    false,
// 			}
// 			antID++
// 		}
// 	}

// 	allMoves := [][]AntMove{}
// 	antsFinished := 0

// 	// STEP 3: Simulation loop
// 	for antsFinished < numAnts {
// 		turnMoves := []AntMove{}
// 		usedRooms := make(map[string]bool)
// 		usedEdges := make(map[string]bool)

// 		// Start can hold unlimited ants
// 		usedRooms[startRoom] = true

// 		for id := 1; id <= numAnts; id++ {
// 			state := antStates[id]

// 			if state.Finished {
// 				continue
// 			}

// 			path := paths[state.PathIndex]
// 			nextStep := state.StepInPath + 1

// 			if nextStep >= len(path) {
// 				continue
// 			}

// 			nextRoom := path[nextStep]

// 			// Create unique tunnel key
// 			edgeKey := state.CurrentRoom + "->" + nextRoom

// 			roomFree := !usedRooms[nextRoom] || nextRoom == endRoom
// 			edgeFree := !usedEdges[edgeKey]

// 			if roomFree && edgeFree {

// 				// Record movement
// 				turnMoves = append(turnMoves, AntMove{
// 					AntID: id,
// 					Room:  nextRoom,
// 				})

// 				// Mark tunnel as used
// 				usedEdges[edgeKey] = true

// 				// Mark room as used (except end)
// 				if nextRoom != endRoom {
// 					usedRooms[nextRoom] = true
// 				}

// 				// Update ant state
// 				state.CurrentRoom = nextRoom
// 				state.StepInPath = nextStep

// 				if nextRoom == endRoom {
// 					state.Finished = true
// 					antsFinished++
// 				}
// 			}
// 		}

// 		if len(turnMoves) > 0 {
// 			allMoves = append(allMoves, turnMoves)
// 		} else {
// 			break
// 		}
// 	}

// 	return allMoves
// }

// // assignAntsToPathsOptimal - Optimal ant distribution across paths
// //
// // Algorithm: Greedy assignment minimizing max completion time
// // For each ant, assign to the path that will complete earliest
// //
// // Completion time for a path = (ants_already_assigned) + (path_length)

// func assignAntsToPathsOptimal(numAnts int, paths [][]string) []int {
// 	numPaths := len(paths)
// 	assignments := make([]int, numPaths)

// 	if numPaths == 0 {
// 		return assignments
// 	}

// 	// Get path lengths
// 	lengths := make([]int, numPaths)
// 	for i := range paths {
// 		lengths[i] = len(paths[i]) - 1
// 	}

// 	// Find minimal T
// 	T := findMinimalTurns(numAnts, lengths)

// 	// Compute ants per path
// 	totalAssigned := 0
// 	for i := 0; i < numPaths; i++ {
// 		ants := T - lengths[i] + 1
// 		if ants < 0 {
// 			ants = 0
// 		}
// 		assignments[i] = ants
// 		totalAssigned += ants
// 	}

// 	// Adjust if we assigned too many
// 	for totalAssigned > numAnts {
// 		for i := numPaths - 1; i >= 0 && totalAssigned > numAnts; i-- {
// 			if assignments[i] > 0 {
// 				assignments[i]--
// 				totalAssigned--
// 			}
// 		}
// 	}

// 	return assignments
// }

// func findMinimalTurns(numAnts int, lengths []int) int {
// 	maxLen := 0
// 	for _, l := range lengths {
// 		if l > maxLen {
// 			maxLen = l
// 		}
// 	}

// 	for T := maxLen; ; T++ {
// 		sum := 0
// 		for _, l := range lengths {
// 			if T-l+1 > 0 {
// 				sum += T - l + 1
// 			}
// 		}
// 		if sum >= numAnts {
// 			return T
// 		}
// 	}
// }

// // FormatMoves - Convert moves to output format
// // Format: "L{antID}-{room} L{antID}-{room} ...\n" for each turn
// func FormatMoves(moves [][]AntMove) string {
// 	result := ""
// 	for _, turnMoves := range moves {
// 		for i, move := range turnMoves {
// 			if i > 0 {
// 				result += " "
// 			}
// 			result += fmt.Sprintf("L%d-%s", move.AntID, move.Room)
// 		}
// 		result += "\n"
// 	}
// 	return result
// }

// // PrintMoves - Print moves to stdout
// func PrintMoves(moves [][]AntMove) {
// 	fmt.Print(FormatMoves(moves))
// }

// // GetTurnCount - Helper to get number of turns
// func GetTurnCount(moves [][]AntMove) int {
// 	return len(moves)
// }

// // VerifyAllAntsReachedEnd - Debugging helper
// func VerifyAllAntsReachedEnd(moves [][]AntMove, numAnts int, endRoom string) bool {
// 	antsAtEnd := make(map[int]bool)

// 	for _, turnMoves := range moves {
// 		for _, move := range turnMoves {
// 			if move.Room == endRoom {
// 				antsAtEnd[move.AntID] = true
// 			}
// 		}
// 	}

// 	return len(antsAtEnd) == numAnts
// }

// func ExtractPaths(anthill map[string][]string, start, end string) [][]string {
// 	paths := [][]string{}
// 	visited := map[string]bool{}

// 	var dfs func(path []string)
// 	dfs = func(path []string) {
// 		current := path[len(path)-1]
// 		if current == end {
// 			newPath := make([]string, len(path))
// 			copy(newPath, path)
// 			paths = append(paths, newPath)
// 			return
// 		}
// 		for _, neighbor := range anthill[current] {
// 			if !visited[neighbor] {
// 				visited[neighbor] = true
// 				dfs(append(path, neighbor))
// 				visited[neighbor] = false
// 			}
// 		}
// 	}

// 	visited[start] = true
// 	dfs([]string{start})
// 	return paths
// }
