package algorithm

import (
	"fmt"
)

// main entry point
func FindAllShortestPaths(start string, end string, graph map[string][]string) ([][]string, error) {
	allPossiblePaths := [][]string{}
	graphCopy := copyGraph(graph)

	for {
		path, err := BFSSinglePath(start, end, graphCopy)
		if err != nil {
			// no more paths were found
			break
		}
		allPossiblePaths = append(allPossiblePaths, path)
		removePathEdges(graphCopy, path)
	}

	if len(allPossiblePaths) == 0 {
		return nil, fmt.Errorf("Error: No paths were found from %s to %s", start, end)
	}
	return allPossiblePaths, nil
}

// the modified version of BFS which returns only one path
func BFSSinglePath(start string, end string, graph map[string][]string) ([]string, error) {
	// in case start and end are the same, there is nothing to proces; early termination
	if start == end {
		return []string{start}, nil
	}

	visited := make(map[string]bool)
	parent := make(map[string]string)
	queue := []string{start}
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == end {
			// end condition
			break
		}

		for _, neighbours := range graph[current] {
			if !visited[neighbours] {
				visited[neighbours] = true
				// parent[C] = A  <- "C, i came from A"
				parent[neighbours] = current
				queue = append(queue, neighbours)
			}
		}
	}

	if !visited[end] {
		return nil, fmt.Errorf("Error: no paths were found from %s to %s", start, end)
	}

	// Backtrack from end to start using parent breadcrumbs
	path := []string{}
	currentNode := end
	for currentNode != "" {
		path = append([]string{currentNode}, path...)
		if currentNode == start {
			break
		}
		currentNode = parent[currentNode]
	}

	return path, nil
}

// remove path edges
func removePathEdges(graph map[string][]string, path []string) {
	for i := 0; i < len(path)-1; i++ {
		current := path[i]
		next := path[i+1]

		// remove the next from current's neighbours
		// and also remove current from next's neighbours
		// the reason is because they are bidirectional
		graph[current] = removeRoom(graph[current], next)
		graph[next] = removeRoom(graph[next], current)
		
	}
}

// remove a room from neighbors list
func removeRoom(rooms []string, target string) []string {
	result := []string{}
	for _, room := range rooms {
		if room != target {
			result = append(result, room)
		}
	}
	return result
}

// deep copy of the graph for mutation 
func copyGraph(graph map[string][]string) map[string][]string {
	graphCopy := make(map[string][]string)
	for room, neighbours := range graph {
		graphCopy[room] = append([]string{}, neighbours...)
	}
	return graphCopy
}
