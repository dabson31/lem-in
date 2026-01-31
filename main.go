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

	antsNum, rooms, roomLinks, err := algorithm.ParseInputs(string(text))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of ants: ", antsNum)
	fmt.Println("Rooms: ", rooms)
	fmt.Println("Room Links: ", roomLinks)

}
