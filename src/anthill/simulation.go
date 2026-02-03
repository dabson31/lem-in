package anthill

import (
	"fmt"
	"strings"
)

func AntsInPath(numAnts int, paths [][]string) []Ant {
    ants := []Ant{}
    
    for antNo := 1; antNo <= numAnts; antNo++ {
        pathIndex := (antNo - 1) % len(paths)
        
        ant := Ant{
            antID: antNo,
            currentRoom: "start",
            pathIndex: pathIndex,
            currentPosition: 0,
            reachedEnd: false,
        }
        ants = append(ants, ant)
    }
    
    return ants
}

func SimulateMovement(ants []Ant, paths [][]string, endRoom string) {
    for {
        allFinished := true
        for _, ant := range ants {
            if !ant.reachedEnd {
                // set a flag switch to check if all ants reached to the end 
                allFinished = false
                break
            }
        }
        if allFinished {
            break
        }
        
        // 
        occupancy := make(map[string]bool)
        // technically doesnt matter if its false since its already set as false 
        // but its better for my readability so idgaf nga
        occupancy["start"] = false  
        occupancy[endRoom] = false
        
        // this here stores "L1-room1", "L2-room2" which we append stuff to it later
        moves := []string{}  
        
        // loop through all ants by index and skip if theyre already at end
        for i := range ants {
            if ants[i].reachedEnd {
                continue
            }
            // you +1 because of how indexing works
            currentPath := paths[ants[i].pathIndex]
            nextPos := ants[i].currentPosition + 1
            
            if nextPos >= len(currentPath) {
                continue  
            }
            
            nextRoom := currentPath[nextPos]
            
            // check if next room is even available
            if nextRoom == endRoom || !occupancy[nextRoom] {
                // set currentRoom to the new room
                // and increment currentPosition to the new position in the path
                ants[i].currentRoom = nextRoom
                ants[i].currentPosition = nextPos
                
                // if no, mark the room as occupied; other ants cant move there this turn
                if nextRoom == endRoom {
                    ants[i].reachedEnd = true
                } else {
                    occupancy[nextRoom] = true
                }
                
                moves = append(moves, fmt.Sprintf("L%d-%s", ants[i].antID, nextRoom))
            }
        }
        
        // print this turn's moves
        if len(moves) > 0 {
            fmt.Println(strings.Join(moves, " "))
        }
    }
}

func PrintOutput(numAnts int, fileContent string, paths [][]string, endRoom string) {
    fmt.Println(numAnts)
    fmt.Println(fileContent)
    fmt.Println()
    SimulateMovement(AntsInPath(numAnts, paths), paths, endRoom)
}