package algorithm

import "fmt"

// BFS-based shortest path extractor
func ExtractShortestPaths(anthill map[string][]string, start, end string) [][]string {
	type queueItem struct {
		Path []string
	}
	queue := []queueItem{{Path: []string{start}}}
	var paths [][]string

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]
		current := item.Path[len(item.Path)-1]

		if current == end {
			newPath := make([]string, len(item.Path))
			copy(newPath, item.Path)
			paths = append(paths, newPath)
			continue
		}

		for _, neighbor := range anthill[current] {
			if !contains(item.Path, neighbor) {
				newPath := append([]string{}, item.Path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, queueItem{Path: newPath})
			}
		}
	}
	return paths
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// Greedy selection of disjoint paths by length (shortest first)
func SelectDisjointPaths(paths [][]string) [][]string {
	disjoint := [][]string{}
	usedRooms := map[string]bool{}

	// Sort paths by length ascending
	for i := 0; i < len(paths)-1; i++ {
		for j := i + 1; j < len(paths); j++ {
			if len(paths[j]) < len(paths[i]) {
				paths[i], paths[j] = paths[j], paths[i]
			}
		}
	}

	for _, path := range paths {
		conflict := false
		for _, room := range path[1 : len(path)-1] { // exclude start/end
			if usedRooms[room] {
				conflict = true
				break
			}
		}
		if !conflict {
			disjoint = append(disjoint, path)
			for _, room := range path[1 : len(path)-1] {
				usedRooms[room] = true
			}
		}
	}
	return disjoint
}

type AntMove struct {
	AntID int
	Room  string
}

type AntState struct {
	CurrentRoom string
	PathIndex   int
	StepInPath  int
	Finished    bool
}

// Main simulation function
func SimulateAnts(numAnts int, paths [][]string, startRoom string, endRoom string) [][]AntMove {
	if len(paths) == 0 || numAnts == 0 {
		return [][]AntMove{}
	}

	pathAssignments := assignAntsToPathsOptimal(numAnts, paths)

	antStates := make([]*AntState, numAnts+1)
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

	for antsFinished < numAnts {
		turnMoves := []AntMove{}
		usedRooms := make(map[string]bool)
		usedEdges := make(map[string]bool)
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
			edgeKey := state.CurrentRoom + "->" + nextRoom
			roomFree := !usedRooms[nextRoom] || nextRoom == endRoom
			edgeFree := !usedEdges[edgeKey]

			if roomFree && edgeFree {
				turnMoves = append(turnMoves, AntMove{AntID: id, Room: nextRoom})
				usedEdges[edgeKey] = true
				if nextRoom != endRoom {
					usedRooms[nextRoom] = true
				}
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

// Optimal assignment using minimal turns formula
func assignAntsToPathsOptimal(numAnts int, paths [][]string) []int {
	numPaths := len(paths)
	assignments := make([]int, numPaths)

	lengths := make([]int, numPaths)
	for i := range paths {
		lengths[i] = len(paths[i]) - 1
	}

	T := findMinimalTurns(numAnts, lengths)

	totalAssigned := 0
	for i := 0; i < numPaths; i++ {
		ants := T - lengths[i] + 1
		if ants < 0 {
			ants = 0
		}
		assignments[i] = ants
		totalAssigned += ants
	}

	// Adjust if over-assigned
	for totalAssigned > numAnts {
		for i := numPaths - 1; i >= 0 && totalAssigned > numAnts; i-- {
			if assignments[i] > 0 {
				assignments[i]--
				totalAssigned--
			}
		}
	}
	return assignments
}

func findMinimalTurns(numAnts int, lengths []int) int {
	maxLen := 0
	for _, l := range lengths {
		if l > maxLen {
			maxLen = l
		}
	}
	for T := maxLen; ; T++ {
		sum := 0
		for _, l := range lengths {
			if T-l+1 > 0 {
				sum += T - l + 1
			}
		}
		if sum >= numAnts {
			return T
		}
	}
}

// Output formatting
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

func PrintMoves(moves [][]AntMove) {
	fmt.Print(FormatMoves(moves))
}

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
