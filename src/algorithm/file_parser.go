package algorithm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// type Room struct {
// 	ID      int
// 	X_coord int
// 	Y_coord int
// }

func ParseInputs(fileInput string) (numAnts int, startRoom string, endRoom string, tunnel_connections map[string][]string, err error) {

	
	var (
		startRooms string 
		endRooms string
		nextRoom string = ""
	)
	startFlagSwitch := false
	endFlagSwitch := false
	seenRoom := make(map[string]bool)
	anthill := make(map[string][]string)


	lines := strings.Split(fileInput, "\n")

	if len(fileInput) < 1 || strings.TrimSpace(fileInput) == "" {
		return 0, "", "", nil, fmt.Errorf("Error: the file being read from is empty.")
	}

	if strings.TrimSpace(lines[1]) != "##start" {
		return 0, "", "", nil, fmt.Errorf("Error: invalid formtting.")
	}

	trimmed := strings.TrimSpace(lines[1])
	firstLineChecker, err := strconv.Atoi(trimmed)

	if err != nil {
		return 0, "", "", nil, fmt.Errorf("Error: invalid formtting.\nFirst line could not be converted into an integer; %v.", firstLineChecker)
	}

	for i := 1; i < len(lines); i++ {

	}

	return 0, "", "", nil, nil
}

// func ParseInput(fileInput string) (int, []Room, []string, error) {
// 	lines := strings.Split(fileInput, "\n")

// 	var ants int
// 	var roomLinks []string
// 	var rooms []Room

// 	antsDone := false
// 	dashLink := false

// 	startRooms := false
// 	endRooms := false


// 	if len(fileInput) < 1 || strings.TrimSpace(fileInput) == "" {
// 		return 0, nil, nil, fmt.Errorf("Error: the file being read from is empty.")
// 	}

// 	// the second line

// 	for i := 0; i < len(lines); i++ {
// 		line := strings.TrimSpace(lines[i])

// 		if line == "" || line[0] == '#' {
// 			continue
// 		}

// 		if !antsDone {
// 			n, err := strconv.Atoi(line)
// 			if err != nil {
// 				return 0, nil, nil, fmt.Errorf("invalid number of ants: %v", err)
// 			}
// 			ants = n
// 			antsDone = true
// 			continue
// 		}

// 		if strings.Contains(line, "-") {
// 			dashLink = true
// 		}

// 		if dashLink {
// 			roomLinks = append(roomLinks, line)
// 			continue
// 		}

// 		//if its not a link then its a room
// 		roomParts := strings.Fields(line)
// 		if len(roomParts) != 3 {
// 			return 0, nil, nil, fmt.Errorf("invalid room definition: %s", line)
// 		}

// 		ID, err := strconv.Atoi(roomParts[0])
// 		if err != nil {
// 			return 0, nil, nil, fmt.Errorf("Error: Could not parse ID to an integer value: %v", err)
// 		}
// 		X_coord, err := strconv.Atoi(roomParts[1])
// 		if err != nil {
// 			return 0, nil, nil, fmt.Errorf("Error: Could not parse x_coord to an integer value: %v", err)
// 		}
// 		Y_coord, err := strconv.Atoi(roomParts[2])
// 		if err != nil {
// 			return 0, nil, nil, fmt.Errorf("Error: Could not parse y_coord to an integer value: %v", err)
// 		}
// 		rooms = append(rooms, Room{ID: ID, X_coord: X_coord, Y_coord: Y_coord})
// 	}

// 	return ants, rooms, roomLinks, nil
// }

//! I am not entirely sure whether this is necessary or not really, but it may be

func ReadText(file string) (string, error) {

	text, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("Error: could not read from file [%v]; %v", file, err)
	}

	return string(text), nil

}
