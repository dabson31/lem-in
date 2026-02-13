package algorithm

import "errors"

func splitTunnels(line string, anthill map[string][]string, seenTunnels map[string]map[string]bool) error {
	// Split the tunnel definition (format: "room1-room2")
	rooms := splitByDash(line)
	if len(rooms) != 2 {
		return errors.New("Error: invalid tunnel format")
	}

	firstRoom := rooms[0]
	secondRoom := rooms[1]

	// Prevent self-loops (room linking to itself)
	if firstRoom == secondRoom {
		return errors.New("Error: room cannot link to itself")
	}

	// Check for duplicate tunnels
	if seenTunnels[firstRoom] == nil {
		seenTunnels[firstRoom] = make(map[string]bool)
	}
	if seenTunnels[secondRoom] == nil {
		seenTunnels[secondRoom] = make(map[string]bool)
	}

	// Check if this tunnel already exists (bidirectional check)
	if seenTunnels[firstRoom][secondRoom] || seenTunnels[secondRoom][firstRoom] {
		return errors.New("Error: duplicate tunnel found")
	}

	seenTunnels[firstRoom][secondRoom] = true
	seenTunnels[secondRoom][firstRoom] = true

	anthill[firstRoom] = append(anthill[firstRoom], secondRoom)
	anthill[secondRoom] = append(anthill[secondRoom], firstRoom)

	return nil
}

func splitByDash(s string) []string {
	result := []string{}
	current := ""
	for _, char := range s {
		if char == '-' {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	result = append(result, current)
	return result
}
