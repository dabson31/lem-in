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

// SelectDisjointPathsMaxFlow uses max-flow algorithm but only accepts paths
// within a reasonable length threshold to avoid long bottleneck paths
func SelectDisjointPathsMaxFlow(anthill map[string][]string, start, end string) [][]string {
	if len(anthill) == 0 {
		return [][]string{}
	}

	// Build flow network
	flow := make(map[string]map[string]int)
	var augmentedPaths [][]string
	minPathLen := -1

	// Keep finding augmenting paths until none exist
	for {
		// Find one augmenting path using BFS
		augmentingPath := findAugmentingPathBFSOnGraph(anthill, start, end, flow)
		if augmentingPath == nil {
			break
		}

		pathLen := len(augmentingPath) - 1

		// On first path, set the baseline length
		if minPathLen == -1 {
			minPathLen = pathLen
		}

		// Only accept paths that are within +5 of the shortest found so far
		// This gives us more flexibility to find better combinations
		if pathLen <= minPathLen+5 {
			// Update flow on this path
			for i := 0; i < len(augmentingPath)-1; i++ {
				u := augmentingPath[i]
				v := augmentingPath[i+1]
				
				if flow[u] == nil {
					flow[u] = make(map[string]int)
				}
				if flow[v] == nil {
					flow[v] = make(map[string]int)
				}
				
				flow[u][v]++
				flow[v][u]--
			}

			augmentedPaths = append(augmentedPaths, augmentingPath)
		} else {
			// Path too long, stop searching
			break
		}
	}

	return augmentedPaths
}

// findAugmentingPathBFSOnGraph finds a path from start to end in the graph
func findAugmentingPathBFSOnGraph(anthill map[string][]string, start, end string, flow map[string]map[string]int) []string {
	queue := [][]string{{start}}
	visited := make(map[string]bool)
	visited[start] = true

	for len(queue) > 0 {
		currentPath := queue[0]
		queue = queue[1:]
		current := currentPath[len(currentPath)-1]

		if current == end {
			return currentPath
		}

		neighbors := anthill[current]
		if neighbors == nil {
			neighbors = []string{}
		}

		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				flowValue := 0
				if flow[current] != nil {
					flowValue = flow[current][neighbor]
				}

				// Can use edge if current flow < 1 (capacity is 1)
				if flowValue < 1 {
					visited[neighbor] = true
					newPath := append([]string{}, currentPath...)
					newPath = append(newPath, neighbor)
					queue = append(queue, newPath)
				}
			}
		}
	}

	return nil // No augmenting path found
}

type AntMove struct {
	AntID int
	Room  string
}

type AntPath struct {
	ID    int
	Path  []string
	Step  int
}

// Main simulation function - dispatch ants gradually (Ahmed's approach)
func SimulateAnts(numAnts int, paths [][]string, startRoom string, endRoom string) [][]AntMove {
	if len(paths) == 0 || numAnts == 0 {
		return [][]AntMove{}
	}

	// Assign ants to paths
	pathAssignments := assignAntsToPathsOptimal(numAnts, paths)

	var activeAnts []AntPath
	antIndexPerPath := make([]int, len(paths))
	finishedCount := 0
	totalAnts := numAnts

	allMoves := [][]AntMove{}

	for finishedCount < totalAnts {
		turnMoves := []AntMove{}
		var stillActive []AntPath

		// Move already active ants forward one step
		for _, ant := range activeAnts {
			ant.Step++
			nextRoom := ant.Path[ant.Step]
			turnMoves = append(turnMoves, AntMove{AntID: ant.ID, Room: nextRoom})

			// Check if ant reached the end
			if ant.Step == len(ant.Path)-1 {
				finishedCount++
			} else {
				stillActive = append(stillActive, ant)
			}
		}
		activeAnts = stillActive

		// Dispatch one new ant from each path that has ants left to dispatch
		for pathIdx := 0; pathIdx < len(paths); pathIdx++ {
			assignedCount := pathAssignments[pathIdx]
			dispatchedCount := antIndexPerPath[pathIdx]

			if dispatchedCount < assignedCount {
				// Create new ant with ID
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

				// Add to active ants if not at end
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

// Greedy ant assignment based on Ahmed's approach
// Assign each ant to the path that currently has the lowest cost
// Cost = path_length + current_ant_count_on_path
func assignAntsToPathsOptimal(numAnts int, paths [][]string) []int {
	numPaths := len(paths)
	assignments := make([]int, numPaths)

	// Initialize ant counts for each path
	pathAntCounts := make([]int, numPaths)

	// Assign each ant to the path with lowest cost
	for antID := 1; antID <= numAnts; antID++ {
		bestPathIdx := 0
		minCost := len(paths[0]) + pathAntCounts[0]

		for i := 1; i < numPaths; i++ {
			cost := len(paths[i]) + pathAntCounts[i]
			if cost < minCost {
				minCost = cost
				bestPathIdx = i
			}
		}

		assignments[bestPathIdx]++
		pathAntCounts[bestPathIdx]++
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
