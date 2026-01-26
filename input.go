package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	ID      int
	X_coord int
	Y_coord int
}

func parseInput(fileInput string) (int, []Room, []string, error) {
	lines := strings.Split(fileInput, "\n")

	var ants int
	var rooms []Room
	var roomLinks []string

	antsDone := false
	dashLink := false

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "" || line[0] == '#' {
			continue
		}

		if !antsDone {
			n, err := strconv.Atoi(line)
			if err != nil {
				return 0, nil, nil, fmt.Errorf("invalid number of ants: %v", err)
			}
			ants = n
			antsDone = true
			continue
		}

		if strings.Contains(line, "-") {
			dashLink = true
		}

		if dashLink {
			roomLinks = append(roomLinks, line)
			continue
		}

		//if its not a link then its a room
		roomParts := strings.Fields(line)
		if len(roomParts) != 3 {
			return 0, nil, nil, fmt.Errorf("invalid room definition: %s", line)
		}

		ID, _ := strconv.Atoi(roomParts[0])
		X_coord, _ := strconv.Atoi(roomParts[1])
		Y_coord, _ := strconv.Atoi(roomParts[2])

		rooms = append(rooms, Room{ID: ID, X_coord: X_coord, Y_coord: Y_coord})
	}

	return ants, rooms, roomLinks, nil
}

func readText(file string) (string, error) {

	text, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	return string(text), nil

}
