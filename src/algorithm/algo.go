package algorithm

import (
	"fmt"
)

func FindAllShortestPaths(start string, end string, graph map[string][]string) ([][]string, error) {
	allPossiblePaths := [][]string{}
	graphCopy := copyGraph(graph)

	for {
		path, err := BFSSinglePath(start, end, graphCopy)
		if err != nil {
			break
		}
		allPossiblePaths = append(allPossiblePaths, path)
		removePathEdges(graphCopy, path, start, end)
	}

	if len(allPossiblePaths) == 0 {
		return nil, fmt.Errorf("Error: No paths were found from %s to %s", start, end)
	}

	if len(allPossiblePaths) >= 5 {
		selectedPaths := selectNodeDisjointPaths(allPossiblePaths, start, end)
		if len(selectedPaths) >= 3 {
			return selectedPaths, nil
		}
	}

	return allPossiblePaths, nil
}

func selectNodeDisjointPaths(paths [][]string, start string, end string) [][]string {
	selected := [][]string{}
	usedNodes := make(map[string]bool)
	usedNodes[start] = true
	usedNodes[end] = true

	sortedPaths := append([][]string{}, paths...)
	for i := 0; i < len(sortedPaths); i++ {
		for j := i + 1; j < len(sortedPaths); j++ {
			if len(sortedPaths[j]) < len(sortedPaths[i]) {
				sortedPaths[i], sortedPaths[j] = sortedPaths[j], sortedPaths[i]
			}
		}
	}

	for _, path := range sortedPaths {
		canUse := true
		for i := 1; i < len(path)-1; i++ {
			if usedNodes[path[i]] {
				canUse = false
				break
			}
		}

		if canUse {
			selected = append(selected, path)
			for i := 1; i < len(path)-1; i++ {
				usedNodes[path[i]] = true
			}
		}
	}

	return selected
}

func selectBestPathCombination(paths [][]string, start string, end string, count int) [][]string {
	if len(paths) <= count {
		return paths
	}

	bestCombination := paths[:count]
	bestOverlapScore := calculateOverlapScore(bestCombination, start, end)

	for i := 0; i < len(paths)-count; i++ {
		for j := 0; j < count; j++ {
			candidate := make([][]string, count)
			copy(candidate, bestCombination)
			candidate[j] = paths[count+i]

			overlapScore := calculateOverlapScore(candidate, start, end)
			if overlapScore < bestOverlapScore {
				bestOverlapScore = overlapScore
				bestCombination = candidate
			}
		}
	}

	return bestCombination
}

func calculateOverlapScore(paths [][]string, start string, end string) int {
	nodeCount := make(map[string]int)

	for _, path := range paths {
		seenInPath := make(map[string]bool)
		for i := 1; i < len(path)-1; i++ {
			node := path[i]
			if !seenInPath[node] {
				nodeCount[node]++
				seenInPath[node] = true
			}
		}
	}

	overlapScore := 0
	for _, count := range nodeCount {
		if count > 1 {
			overlapScore += count - 1
		}
	}

	return overlapScore
}

func BFSSinglePath(start string, end string, graph map[string][]string) ([]string, error) {
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
			break
		}

		for _, neighbours := range graph[current] {
			if !visited[neighbours] {
				visited[neighbours] = true
				parent[neighbours] = current
				queue = append(queue, neighbours)
			}
		}
	}

	if !visited[end] {
		return nil, fmt.Errorf("Error: no paths were found from %s to %s", start, end)
	}

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

func removePathEdges(graph map[string][]string, path []string, start string, end string) {
	for i := 0; i < len(path)-1; i++ {
		current := path[i]
		next := path[i+1]

		if current == start && next == end {
			graph[current] = removeRoom(graph[current], next)
			graph[next] = removeRoom(graph[next], current)
			continue
		}

		if current == start && next != end {
			continue
		}

		if next == end && current != start {
			continue
		}

		graph[current] = removeRoom(graph[current], next)
		graph[next] = removeRoom(graph[next], current)
	}
}

func removeRoom(rooms []string, target string) []string {
	result := []string{}
	for _, room := range rooms {
		if room != target {
			result = append(result, room)
		}
	}
	return result
}

func copyGraph(graph map[string][]string) map[string][]string {
	graphCopy := make(map[string][]string)
	for room, neighbours := range graph {
		graphCopy[room] = append([]string{}, neighbours...)
	}
	return graphCopy
}
