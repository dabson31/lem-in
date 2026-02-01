package algorithm

import (
	"fmt"
	"strings"
)

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

	// here we checkif 2 is connected to 0, for example
	// so if room 0's map does not contain 2, for example, it write it down
	// and within the map, it becomes true
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