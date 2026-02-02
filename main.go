package main

import (
	"fmt"
	"lem-in/src/algorithm"
	"log"
	"os"
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

	antsNum, startRoom, endRoom, tunnel_connections, err := algorithm.ParseInputs(string(text))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of ants: ", antsNum)
	fmt.Println("Start Room: ", startRoom)
	fmt.Println("End Room: ", endRoom)
	fmt.Println("The tunnels: ", tunnel_connections)

	paths, err := algorithm.BFS(startRoom, endRoom, tunnel_connections)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(paths)

}
