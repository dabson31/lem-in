package main

import (
	"fmt"
	"os"

	"lem-in/src/algorithm"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: usage: go run . <input_file>")
		return
	}

	filename := os.Args[1]
	text, err := algorithm.ReadText(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	numAnts, startRoom, endRoom, anthill, err := algorithm.ParseInputs(text)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	// Extract all shortest paths using BFS
	paths := algorithm.ExtractShortestPaths(anthill, startRoom, endRoom)
	if len(paths) == 0 {
		fmt.Println("ERROR: no paths found from start to end")
		return
	}

	// Choose disjoint paths (greedy by length)
	disjointPaths := algorithm.SelectDisjointPaths(paths)

	// Simulate ants on these paths
	allMoves := algorithm.SimulateAnts(numAnts, disjointPaths, startRoom, endRoom)

	// Print moves in required format
	algorithm.PrintMoves(allMoves)

	// Optional: verify all ants reached the end
	if !algorithm.VerifyAllAntsReachedEnd(allMoves, numAnts, endRoom) {
		fmt.Println("ERROR: some ants did not reach the end")
	}
}
