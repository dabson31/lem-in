package main

import (
	"fmt"
	"lem-in/src/graphs"
)

func main() {
	test := &graphs.Graph{}
	test.AddVertex(2)

	for i := 0; i < 5; i++ {
		test.AddVertex(i)
		fmt.Println("slice containing pointers:        ", test.Vertices[i:])
		fmt.Println("address of pointer in slice:      ", &test.Vertices[i])
		fmt.Println("pointer value (vertex address):   ", test.Vertices[i])
		fmt.Println("vertex key:                       ", test.Vertices[i].Key)
		fmt.Println()

	}
	fmt.Println(test)
}
