package main

import (
	"fmt"
	"log"
	"os"
	//"lem-in/src/graphs"
)

func main() {

	args := os.Args

	if len(args) != 2 {
		fmt.Println("Usage: go run . input_file")
		return
	}

	// if len(args) == 1 {
	// 	fmt.Println("No file provided")
	// 	return
	// } else if len(args) > 2 {
	// 	fmt.Println("Too many arguments provided")
	// 	return
	// }

	text, err := os.ReadFile(args[1])
	if err != nil {
		log.Fatal(err)
	}

	antsNum, rooms, roomLinks, err := parseInput(string(text))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of ants: ", antsNum)
	fmt.Println("Rooms: ", rooms)
	fmt.Println("Room Links: ", roomLinks)

	// test := &graphs.Graph{}
	// test.AddVertex(2)

	// for i := 0; i < 5; i++ {
	// 	test.AddVertex(i)
	// 	fmt.Println("slice containing pointers:        ", test.Vertices[i:])
	// 	fmt.Println("address of pointer in slice:      ", &test.Vertices[i])
	// 	fmt.Println("pointer value (vertex address):   ", test.Vertices[i])
	// 	fmt.Println("vertex key:                       ", test.Vertices[i].Key)
	// 	fmt.Println()

	// }
	// fmt.Println(test)

}
