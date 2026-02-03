package main

import (
	"fmt"
	"lem-in/src/algorithm"
	"lem-in/src/anthill"
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

	antsNum, startRoom, endRoom, tunnel_connections, err := algorithm.ParseInputs(string(text))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of ants: ", antsNum)
	fmt.Println("Start Room: ", startRoom)
	fmt.Println("End Room: ", endRoom)
	fmt.Println("The tunnels: ", tunnel_connections)

	paths, err := algorithm.FindAllShortestPaths(startRoom, endRoom, tunnel_connections)
	if err != nil {
		log.Fatal(err)
	}

	ants := anthill.AntsInPath(antsNum, paths)
	anthill.SimulateMovement(ants, paths, endRoom)

	fmt.Printf("time elapsed: %v\n", time.Since(startTimer))
}
