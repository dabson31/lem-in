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

	roomType := ""
	startFlagSwitch := false
	endFlagSwitch := false
	seenRoom := make(map[string]bool)
	anthill := make(map[string][]string)
	// loop variables
	seenTunnels := make(map[string]map[string]bool)
	lines := strings.Split(fileInput, "\n")

	// the following two if statements are somewhat "restricted", but it allows for better control and error checking
	if len(fileInput) < 1 || strings.TrimSpace(fileInput) == "" {
		return 0, "", "", nil, fmt.Errorf("Error: the file being read from is empty.")
	}
	if strings.TrimSpace(lines[1]) != "##start" {
		return 0, "", "", nil, fmt.Errorf("Error: invalid formtting.")
	}

	trimmed := strings.TrimSpace(lines[0])
	firstLineChecker, err := strconv.Atoi(trimmed)

	if err != nil {
		return 0, "", "", nil, fmt.Errorf("Error: invalid formtting.\nFirst line could not be converted into an integer; %v.", firstLineChecker)
	}
	numAnts = firstLineChecker

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// once we find the first two commands, we set their bool to true respectively
		// if there are any duplicate commands we return an error just to make extra sure
		// that there will be no errors within the parsing logic
		// we ignore comments that start with a single #
		if line == "##start" {
			if startFlagSwitch {
				return 0, "", "", nil, fmt.Errorf("Error: duplicate ##start found")
			}
			// next toom is the start room; contains values pertaining to that
			startFlagSwitch = true
			roomType = "start"
			continue
		}
		if line == "##end" {
			if endFlagSwitch {
				return 0, "", "", nil, fmt.Errorf("Error: duplicate ##end found")
			}
			endFlagSwitch = true
			roomType = "end"
			continue
		}
		// comments ignoring
		if strings.HasPrefix(line, "#") /*|| strings.HasSuffix(line, "#") i dont think we need to handle suffixes too */ {
			continue
		}

		if strings.Contains(line, "-") {
			rooms := strings.Split(line, "-")

			if len(rooms) != 2 {
				return 0, "", "", nil, fmt.Errorf("Error: invalid tunnel format; %v", line)
			}

			firstRoom, secondRoom := rooms[0], rooms[1]

			// in case of trying to link rooms that do not exist
			if seenRoom[firstRoom] == false || seenRoom[secondRoom] == false {
				return 0, "", "", nil, fmt.Errorf("Error: tunnel references unknown room; %v", line)
			}

			err := splitTunnels(line, anthill, seenTunnels)
			if err != nil {
				return 0, "", "", nil, err
			}
			continue
		}

		// we're just parsing rooms here in case it's not a tunnel, comment or a command
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return 0, "", "", nil, fmt.Errorf("Error: invalid room definition. [USAGE]: <name> <x coord> <y coord>")
		}

		roomName := fields[0]

		if seenRoom[roomName] {
			return 0, "", "", nil, fmt.Errorf("Error: duplicate room found.")
		}
		_, err = strconv.Atoi(fields[1])
		if err != nil {
			return 0, "", "", nil, fmt.Errorf("Error: invalid x coordinate.")
		}
		_, err = strconv.Atoi(fields[2])
		if err != nil {
			return 0, "", "", nil, fmt.Errorf("Error: invalid y coordinate.")
		}

		seenRoom[roomName] = true
		anthill[roomName] = []string{}

		if roomType == "start" {
			startRoom = roomName
			roomType = ""
		}

		if roomType == "end" {
			endRoom = roomName
			roomType = ""
		}
	}
	// here is the main parsing after checking everything

	if startRoom == "" || endRoom == "" {
		return 0, "", "", nil, fmt.Errorf("Error: missing start or end room")
	}

	return numAnts, startRoom, endRoom, anthill, nil
}

//! ignore this:
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
