package algorithm

import "fmt"

func BFS(start, end string, graph map[string][]string) ([][]string, error) {
	// if the start and end is the same why even run the algorithm
	if start == end {
		return [][]string{{start}}, nil
	}

	// we initialize the queue, distance and the parents map(backtrack)
	// dist is not being used to calculate distance using weightage (its an unweighted graph)
	// it is there to keep track of the turns taken to reach the end node from the start node
	// parents stores the parent nodes of each room within the shortest path/s which enables
	// backtracking from the end to the start so we generate all possible paths.
	// each room might have multiple parent nodes, which means multiple shortest paths exist to reach that room

	queue := []string{start}
	dist := make(map[string]int)
	dist[start] = 0
	parents := make(map[string][]string)

	for i := 0; i < len(queue); i++ {
		current := queue[i]
		neighbors := graph[current]

		// first if condition is for discovering and visiting new nodes
		// second if condition is if we find another path to the discovered nodes
		// with the same distance as first path, meaning more than one shortest path exists
		for j := 0; j < len(neighbors); j++ {
			neighbor := neighbors[j]
			_, exists := dist[neighbor]
			if !exists {
				dist[neighbor] = dist[current] + 1
				parents[neighbor] = []string{current}
				queue = append(queue, neighbor)
			} else if dist[neighbor] == dist[current]+1 {
				parents[neighbor] = append(parents[neighbor], current)
			}
		}
	}

	// unreachable check
	_, valid := dist[end]
	if !valid {
		return nil, fmt.Errorf("No path found from %s to %s", start, end)
	}

	// this backtracking is lacking, it only gives the shortest paths in terms of distance
	// but does not enumerate all possible path, need to change that to somehow get the shortest
	// 5 paths for example, not sure what ill do exactly yet
	paths := [][]string{}
	var backtrack func(node string, path []string)
	backtrack = func(node string, path []string) {
		path = append([]string{node}, path...)
		if node == start {
			tmp := make([]string, len(path))
			copy(tmp, path)
			paths = append(paths, tmp)
			return
		}
		for i := 0; i < len(parents[node]); i++ {
			backtrack(parents[node][i], path)
		}
	}

	backtrack(end, []string{})
	return paths, nil
}

// when we reach the start node we complete one path and append it to "paths"
// -----------------------------------------------------------------
// this is just a standard one-path BFS implementation for reference
// -----------------------------------------------------------------

// func BFS(start, end string, graph map[string][]string) ([]string, error) {

// 	queue := []string{start}
// 	visited := make(map[string]bool)
// 	parent := make(map[string]string)
// 	visited[start] = true

// 	for len(queue) > 0 {
// 		current := queue[0]
// 		queue = queue[1:]

// 		if current == end {
// 			break
// 		}

// 		neighbors := graph[current]
// 		for i := 0; i < len(neighbors); i++ {
// 			neighbor := neighbors[i]

// 			if !visited[neighbor] {
// 				visited[neighbor] = true
// 				parent[neighbor] = current
// 				queue = append(queue, neighbor)
// 			}
// 		}
// 	}

// 	if !visited[end] {
// 		return nil, fmt.Errorf("No path found from %s to %s", start, end)
// 	}

// 	path := []string{}
// 	for at := end; at != ""; {
// 		path = append([]string{at}, path...)
// 		if at == start {
// 			break
// 		}
// 		at = parent[at]
// 	}

// 	return path, nil
// }
