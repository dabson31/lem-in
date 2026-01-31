package algorithm

import (
	"fmt"
	"strings"
)

// here, well check for duplicate rooms
func doesRoomExists(rooms []string, duplicate string) bool {
	for _, room := range rooms {
		if room == duplicate {
			return true
		}
	}
	return false
}

func splitTunnels(lines string, antColony map[string][]string, seen map[string]map[string]bool) error {

	//? originally the implementation required looping through the array
	//? every time it checks, but using map for tracking allows for O(1) lookup

	splitRooms := strings.Split(lines, "-")

	if len(splitRooms) != 2 {
		return fmt.Errorf("Error: invalid tunnel format; %v", lines)
	}

	firstRoom, secondRoom := splitRooms[0], splitRooms[1]

	if seen[firstRoom] == nil {
		seen[firstRoom] = make(map[string]bool)
	}
	
	if seen[secondRoom] == nil {
		seen[secondRoom] = make(map[string]bool)
	}

	if !seen[firstRoom][secondRoom] {
		antColony[firstRoom] = append(antColony[firstRoom], secondRoom)
		seen[firstRoom][secondRoom] = true
	}

	if !seen[secondRoom][firstRoom] {
		antColony[secondRoom] = append(antColony[secondRoom], firstRoom)
		seen[secondRoom][firstRoom] = true
	}
	return nil
}