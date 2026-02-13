package main

import (
	"fmt"
	"lem-in/src/algorithm"
	"log"
	"os"
	"time"
)

func main() {

	args := os.Args

	if len(args) != 2 {
		fmt.Println("Usage: go run . examples/[*.txt]")
		return
	}

	text, err := os.ReadFile(args[1])
	if err != nil {
		log.Fatal(err)
	}

	startTimer := time.Now()

	fmt.Println("=== PARSING INPUT ===")
	antsNum, startRoom, endRoom, tunnel_connections, err := algorithm.ParseInputs(string(text))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of ants:", antsNum)
	fmt.Println("Start Room:", startRoom)
	fmt.Println("End Room:", endRoom)
	fmt.Println("Graph connections:")
	for room, neighbors := range tunnel_connections {
		fmt.Printf("  %s -> %v\n", room, neighbors)
	}
	fmt.Println()

	// DEBUG: Check if start and end rooms exist in the graph
	if _, exists := tunnel_connections[startRoom]; !exists {
		log.Fatal("ERROR: Start room not found in graph!")
	}
	if _, exists := tunnel_connections[endRoom]; !exists {
		log.Fatal("ERROR: End room not found in graph!")
	}

	fmt.Println("=== FINDING PATHS ===")
	paths, err := algorithm.FindAllShortestPaths(startRoom, endRoom, tunnel_connections)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d path(s):\n", len(paths))
	for i, path := range paths {
		fmt.Printf("  Path %d (%d steps): %v\n", i+1, len(path)-1, path)
	}
	fmt.Println()

	fmt.Println("=== SIMULATING ANT MOVEMENT ===")
	moves := algorithm.SimulateAnts(antsNum, paths, startRoom, endRoom)

	fmt.Printf("Total turns: %d\n", len(moves))
	fmt.Println()

	fmt.Println("=== ANT MOVEMENTS ===")
	algorithm.PrintMoves(moves)

	fmt.Printf("\n=== COMPLETED IN: %v ===\n", time.Since(startTimer))
}
