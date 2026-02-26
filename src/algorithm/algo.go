package algorithm

import "fmt"

type AntMove struct {
	AntID int
	Room  string
}

type AntPath struct {
	ID   int
	Path []string
	Step int
}

// Bfs
// func ExtractShortestPaths(anthill map[string][]string, start, end string) [][]string {
// 	// the BFS we have here doesnt track distance/parents, instead it stores
// 	// them in the full path in the queue. basically every queue item is one possible
// 	// path from start to current node. this makes backtracking unnecessary because
// 	// a path is always preserved/complete in the queue.

// 	type queueItem struct {
// 		Path []string
// 	}

// 	queue := []queueItem{{Path: []string{start}}} // start queue with start room

// 	var paths [][]string

// 	for len(queue) > 0 {
// 		// pop first element (first come first out (FIFO))
// 		item := queue[0]
// 		queue = queue[1:]

// 		current := item.Path[len(item.Path)-1] // current is the last element in the stored path

// 		// if it reaches the end room, we copy the path
// 		// and store it to var paths to use later
// 		if current == end {
// 			newPath := make([]string, len(item.Path))
// 			copy(newPath, item.Path)
// 			paths = append(paths, newPath)
// 			continue
// 		}

// 		// otherwise, check neighbors of current and expand bfs tree
// 		// we only add neighbors that are not already in the path to avoid cycles
// 		// we copy the current path before modifying to make sure each path
// 		// in the queue is independent and complete, then we push extended path to the queue

// 		for _, neighbor := range anthill[current] {
// 			if !strings.contains(item.Path, neighbor) {
// 				newPath := append([]string{}, item.Path...)
// 				newPath = append(newPath, neighbor)
// 				queue = append(queue, queueItem{Path: newPath})
// 			}
// 		}
// 	}

// 	// this returns ALL paths from start to end, not filtered by shortest
// 	return paths
// }



// this is where we extract multiple seperate paths using max-flow style (Edmonds-Karp)
// we keep finding paths and updating flow until no more can be found.
// and paths have to stay within a reasonable length so we avoid using
// very long paths that would increase the total turns significantly.

// =====================================================================
// WE DO NOT NEED THE REVERSE RESIDUAL PATHS FOR LEM-IN BECAUSE WE ONLY
// CARE ABOUT THE DISJOINT PATHS, IM KEEPING IT FOR PROPER IMPLEMENTATION
// =====================================================================
func SelectDisjointPathsMaxFlow(anthill map[string][]string, start, end string) [][]string {
	if len(anthill) == 0 { // empty graph check
		return [][]string{}
	}

	// it gets shway complicated here, this is the flow
	// flow[u][v] = 1 means the edge is used in current path
	// flow[u][v] = -1 means the edge is used in reverse direction
	// flow[u][v] = 0 means the edge is not used
	// we only allow flow of 1 in either direction because we want seperated paths
	flow := make(map[string]map[string]int)

	var augmentedPaths [][]string // stores all valid paths
	minPathLen := -1              // stores shortest path found first

	// keep searching for augmenting paths until none can be found or they are too long (or bfs fails to find any path)
	for {
		// find one augmenting path using bfs in residual graph (anthill + flow)
		augmentingPath := findAugmentingPathBFSOnGraph(anthill, start, end, flow)
		if augmentingPath == nil { // path doesnt exist, max flow reached and we're done
			break
		}

		pathLen := len(augmentingPath) - 1 // (number of edges in path)

		// first path is the shortest baseline
		// that we compare other paths to later
		if minPathLen == -1 {
			minPathLen = pathLen
		}

		// if the length of the path is too long we reject it
		// in this case, length <= shortest + 5, this prevents
		// long inefficient paths that would unoptimize the program
		if pathLen <= minPathLen+5 {
			for i := 0; i < len(augmentingPath)-1; i++ {
				// we're updating flows here based on the paths we found
				u := augmentingPath[i]
				v := augmentingPath[i+1]

				if flow[u] == nil {
					flow[u] = make(map[string]int)
				}
				if flow[v] == nil {
					flow[v] = make(map[string]int)
				}

				flow[u][v]++ // 1 unit of flow in forward direction
				flow[v][u]--
			}

			augmentedPaths = append(augmentedPaths, augmentingPath)
		} else {
			break
		}
	}

	// returns all suitable disjoint paths (flow-based)
	return augmentedPaths
}

// this bfs is used to find one path in the residual graph (anthill + flow)
// we find ONE valid augmenting path from the start to the end node
// an augmenting path here is basically one where edges are not being used up by flow
// capacity is 1 for each node, meaniong we can only use an edge if flow is less than 1
func findAugmentingPathBFSOnGraph(anthill map[string][]string, start, end string, flow map[string]map[string]int) []string {
	queue := [][]string{{start}} // we store full path, not only the start node iirc
	visited := make(map[string]bool)
	visited[start] = true

	for len(queue) > 0 {
		currentPath := queue[0] // pop first path (FIFO)
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
			if !visited[neighbor] { // flow update
				flowValue := 0
				if flow[current] != nil {
					flowValue = flow[current][neighbor]
				}

				if flowValue < 1 {
					visited[neighbor] = true
					newPath := append([]string{}, currentPath...)
					newPath = append(newPath, neighbor)
					queue = append(queue, newPath)
				}
			}
		}
	}

	return nil
}


// assign each ant to the path that currently has the lowest cost
// each ant here chooses the path that currently gives the lowest cost
// which is calculated as path length + number of ants already assigned
// to said path, which balances shorter paths while avoiding overcrowding
func assignAntsToPathsOptimal(numAnts int, paths [][]string) []int {

	numPaths := len(paths)

	assignments := make([]int, numPaths)   // how many ants assigned to each path
	pathAntCounts := make([]int, numPaths) // track said ants

	// assign each ant to the path with lowest cost one by one
	for antID := 1; antID <= numAnts; antID++ {
		bestPathIdx := 0 // set first path as best by default
		minCost := len(paths[0]) + pathAntCounts[0]

		for i := 1; i < numPaths; i++ { // check all other paths to find better cost
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

func CountPrintedLines(moves [][]AntMove) int {
	return len(moves)
}

func PrintMoves(moves [][]AntMove) {
	fmt.Print(FormatMoves(moves))

	totalLines := CountPrintedLines(moves)
	fmt.Println("===============")
	fmt.Print("TOTAL MOVES:")
	fmt.Println(totalLines)
	fmt.Println("===============")

}

// unnecessary check, but just in case we miss something in the simulation we verify it here.
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
