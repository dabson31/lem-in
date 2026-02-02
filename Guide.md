# Lem-in: Complete Implementation Guide

## A University-Level Tutorial on Graph-Based Path Optimization

---

## Table of Contents

1. [Understanding the Problem](#understanding-the-problem)
2. [Core Concepts](#core-concepts)
3. [Input Format Breakdown](#input-format-breakdown)
4. [Data Structures You'll Build](#data-structures-youll-build)
5. [Implementation Architecture](#implementation-architecture)
6. [Step-by-Step Implementation](#step-by-step-implementation)
7. [Algorithm Explanations](#algorithm-explanations)
8. [Common Pitfalls & Solutions](#common-pitfalls--solutions)
9. [Testing Strategy](#testing-strategy)

---

## Understanding the Problem

### The Big Picture

Imagine an ant colony represented as a graph where:

- **Nodes** = Rooms
- **Edges** = Tunnels connecting rooms
- **Start** = Where all ants begin
- **End** = Where all ants must reach

Your job: Get N ants from start to end in the **minimum number of turns**.

### The Challenge

Each ant takes multiple turns to traverse paths. If you use only ONE path, ants must go sequentially - taking N turns minimum for N ants. But if you find MULTIPLE PATHS, ants can move in parallel, reducing total turns significantly.

### Real-World Analogy

Think of a highway network:

- **Single path**: Traffic bottleneck, everyone waits
- **Multiple paths**: Distribute traffic across highways, everything moves faster

### Constraints (The Rules of the Game)

1. **Room Capacity**: Each room (except start/end) holds **exactly 1 ant per turn**
   - Start and end can hold unlimited ants
   - If a room is occupied, another ant cannot enter that turn

2. **Tunnel Usage**: Each tunnel used **once per turn**
   - If ant A uses tunnel X→Y in turn 1, ant B cannot use X→Y in turn 1
   - But ant B CAN use it in turn 2

3. **Simultaneous Movement**: All ants that can move DO move in the same turn
   - Turn-based, not sequential

4. **Valid Paths Only**: Ants must only move through tunnels (connected rooms)
   - No teleporting between non-connected rooms

---

## Core Concepts

### What is a Room?

A room is a node in your graph with three attributes:

```go
type Room struct {
    name        string
    x, y        int         // Coordinates (ignored for pathfinding)
    connections []*Room     // Adjacent rooms via tunnels
    occupied    bool        // Tracking during simulation
}
```

**Key insight**: Room coordinates don't affect pathfinding. They're just metadata. Two rooms are connected if there's a tunnel between them.

**Usage example**:

```go
// Create a room
room0 := &Room{name: "0", x: 0, y: 3, connections: []*Room{}}

// Access room properties
fmt.Println(room0.name)  // Output: "0"
fmt.Println(room0.x, room0.y)  // Output: 0 3
```

### What is a Tunnel?

A tunnel is a bidirectional edge between two rooms.

```
room1-room2  // Creates connections: room1 ↔ room2
```

This means:

- An ant in room1 can move to room2
- An ant in room2 can move to room1
- The tunnel itself can only be used once per turn per ant direction

### What is a Path?

A path is a sequence of rooms from start to end with NO REPEATED ROOMS.

```
Path example: [start] → [0] → [2] → [3] → [end]
Length: 5 rooms (or 4 moves)
```

**Why no repeated rooms?**

- Ants never need to go backward
- Backward movement creates infinite loops
- The shortest path is always forward

### What is Path Disjointness?

Two paths are **edge-disjoint** if they share NO tunnels (edges).

```
Path 1: start → A → B → end    (uses tunnels: start-A, A-B, B-end)
Path 2: start → C → D → end    (uses tunnels: start-C, C-D, D-end)
✓ Edge-disjoint: No shared tunnels

Path 1: start → A → B → end    (uses tunnels: start-A, A-B, B-end)
Path 2: start → A → C → end    (uses tunnels: start-A, A-C, C-end)
✗ NOT edge-disjoint: Both use start-A
```

**Why does this matter?**

- Edge-disjoint paths prevent tunnel bottlenecks
- Ants never compete for the same tunnel
- Maximum parallelization possible

---

## Input Format Breakdown

### Visual Example

```
4                    ← Number of ants to move
##start              ← Marker: next room definition is the start
0 0 3                ← Room "0" at coordinates (0, 3)
2 2 5                ← Room "2" at coordinates (2, 5)
3 4 0                ← Room "3" at coordinates (4, 0)
##end                ← Marker: next room definition is the end
1 8 3                ← Room "1" at coordinates (8, 3)
0-2                  ← Tunnel connecting rooms "0" and "2"
2-3                  ← Tunnel connecting rooms "2" and "3"
3-1                  ← Tunnel connecting rooms "3" and "1"
```

### Parsing Logic

**Step 1**: Read first line → **Number of ants** (e.g., 4)

**Step 2**: Find `##start` marker → Next room is **start room**

**Step 3**: Read room definitions (format: `name x y`)

```
name = string identifier (unique, cannot start with 'L' or '#')
x, y = integer coordinates
```

**Step 4**: Find `##end` marker → Next room is **end room**

**Step 5**: Read remaining room definitions

**Step 6**: Read tunnel definitions (format: `room1-room2`)

```
Creates bidirectional connection: room1 ↔ room2
```

### Comments and Whitespace

- Lines starting with `#` (but NOT `##`) are comments → **IGNORE**
- Empty lines → **IGNORE**
- Only `##start` and `##end` are special markers
- Everything else is either a room or tunnel definition

### Real Example: example00.txt

```
4                    ← 4 ants to move
##start
0 0 3                ← Room "0" (start)
2 2 5                ← Room "2"
3 4 0                ← Room "3"
##end
1 8 3                ← Room "1" (end)
0-2                  ← 0 ↔ 2
2-3                  ← 2 ↔ 3
3-1                  ← 3 ↔ 1

Resulting Graph:
    0 ------- 2 ------- 3 ------- 1
  (start)                          (end)
```

Only one path exists: `0 → 2 → 3 → 1`

---

## Data Structures You'll Build

### The Room Structure

```go
type Room struct {
    name        string
    x, y        int
    connections []*Room
    occupied    bool
}

// Factory function to create rooms
func NewRoom(name string, x, y int) *Room {
    return &Room{
        name: name,
        x: x,
        y: y,
        connections: []*Room{},
        occupied: false,
    }
}
```

**Usage**:

```go
room0 := NewRoom("0", 0, 3)
room1 := NewRoom("1", 8, 3)

// Add connection (bidirectional)
room0.connections = append(room0.connections, room1)
room1.connections = append(room1.connections, room0)
```

### The Graph (Colony)

```go
type Colony struct {
    rooms      map[string]*Room   // Map for O(1) lookup
    start      *Room              // Starting room
    end        *Room              // Ending room
    numAnts    int                // Number of ants
}

// Factory function
func NewColony() *Colony {
    return &Colony{
        rooms: make(map[string]*Room),
    }
}
```

**Usage**:

```go
colony := NewColony()

// Add rooms
room0 := NewRoom("0", 0, 3)
room1 := NewRoom("1", 8, 3)

colony.rooms["0"] = room0
colony.rooms["1"] = room1
colony.start = room0
colony.end = room1
colony.numAnts = 4

// Access any room by name in O(1)
r := colony.rooms["0"]  // Get room "0"
```

### The Ant

```go
type Ant struct {
    id       int        // Unique identifier (1, 2, 3, ...)
    path     []*Room    // Assigned path to follow
    position int        // Current index in path
    arrived  bool       // Reached the end?
}

// Factory function
func NewAnt(id int, path []*Room) *Ant {
    return &Ant{
        id: id,
        path: path,
        position: 0,
        arrived: false,
    }
}

// Get current room
func (a *Ant) CurrentRoom() *Room {
    if a.position >= len(a.path) {
        return nil
    }
    return a.path[a.position]
}

// Try to move to next room
func (a *Ant) CanMove() bool {
    return a.position < len(a.path)-1
}

func (a *Ant) Move() {
    if a.CanMove() {
        a.position++
        if a.path[a.position] == a.path[len(a.path)-1] {
            a.arrived = true
        }
    }
}
```

**Usage**:

```go
path := []*Room{room0, room2, room3, room1}
ant1 := NewAnt(1, path)

fmt.Println(ant1.CurrentRoom().name)  // Output: "0" (start)
ant1.Move()
fmt.Println(ant1.CurrentRoom().name)  // Output: "2"
```

### The Path

```go
type Path struct {
    rooms  []*Room
}

// Get path length (number of moves)
func (p *Path) Length() int {
    return len(p.rooms) - 1
}

// Get room at specific index
func (p *Path) Room(i int) *Room {
    if i >= 0 && i < len(p.rooms) {
        return p.rooms[i]
    }
    return nil
}

// Get middle rooms (excluding start and end)
func (p *Path) MiddleRooms() []*Room {
    if len(p.rooms) <= 2 {
        return []*Room{}
    }
    return p.rooms[1 : len(p.rooms)-1]
}
```

**Usage**:

```go
path := &Path{rooms: []*Room{room0, room2, room3, room1}}

fmt.Println(path.Length())  // Output: 3 (4 rooms - 1)
fmt.Println(path.Room(1).name)  // Output: "2"
```

---

## Implementation Architecture

### The Three-Phase Process

```
┌──────────────────────────────────────────────────────┐
│ PHASE 1: PARSE & BUILD                               │
│ Read file → Extract ants, rooms, tunnels             │
│ Build graph data structure                           │
└─────────────┬──────────────────────────────────────────┘
              │
┌─────────────▼──────────────────────────────────────────┐
│ PHASE 2: FIND PATHS                                  │
│ Find all possible paths from start to end            │
│ Select optimal set of paths                          │
└─────────────┬──────────────────────────────────────────┘
              │
┌─────────────▼──────────────────────────────────────────┐
│ PHASE 3: SIMULATE & OUTPUT                           │
│ Assign ants to paths                                 │
│ Move ants turn-by-turn                               │
│ Output movements in required format                  │
└──────────────────────────────────────────────────────────┘
```

### Module Organization (Recommended)

```
main.go              ← Entry point, orchestration
├── parser.go        ← Read file, extract data
├── validator.go     ← Check for invalid input
├── pathfinder.go    ← Find all paths from start to end
├── optimizer.go     ← Select best set of paths
├── simulator.go     ← Move ants, track positions
└── formatter.go     ← Format output
```

Each module handles one concern. Makes testing and debugging easier.

---

## Step-by-Step Implementation

### Phase 1: Parsing the Input

#### Step 1A: Read the File

```
Input: filename
Output: []string of lines

1. Open file
2. Read all content
3. Split into lines (handle both \n and \r\n)
4. Trim whitespace
```

**Go code**:

```go
import (
    "os"
    "strings"
)

func ReadFile(filename string) ([]string, error) {
    content, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    // Handle both \n and \r\n
    text := strings.ReplaceAll(string(content), "\r\n", "\n")

    // Split into lines and trim
    lines := strings.Split(strings.TrimSpace(text), "\n")

    return lines, nil
}
```

**Usage**:

```go
lines, err := ReadFile("example.txt")
if err != nil {
    fmt.Println("ERROR: invalid data format")
    return
}
// lines = ["4", "##start", "0 0 3", ...]
```

#### Step 1B: Extract Number of Ants

```
First line = number of ants
Validate: Must be > 0

If invalid → print "ERROR" and exit
```

**Go code**:

```go
import "strconv"

func ExtractAntCount(lines []string) (int, error) {
    if len(lines) == 0 {
        return 0, fmt.Errorf("empty input")
    }

    numAnts, err := strconv.Atoi(strings.TrimSpace(lines[0]))
    if err != nil || numAnts <= 0 {
        return 0, fmt.Errorf("invalid ant count")
    }

    return numAnts, nil
}
```

**Usage**:

```go
numAnts, err := ExtractAntCount(lines)
if err != nil {
    fmt.Println("ERROR: invalid data format")
    return
}
fmt.Println(numAnts)  // Output: 4
```

#### Step 1C: Find Start and End Rooms

```
Iterate through lines:
  If line == "##start":
    Next room definition is START
  If line == "##end":
    Next room definition is END

Validate: Exactly one start and one end
```

**Go code**:

```go
func FindStartEnd(lines []string) (startIdx, endIdx int, err error) {
    startIdx, endIdx = -1, -1

    for i, line := range lines {
        trimmed := strings.TrimSpace(line)

        if trimmed == "##start" {
            if startIdx != -1 {
                return -1, -1, fmt.Errorf("multiple ##start markers")
            }
            startIdx = i
        } else if trimmed == "##end" {
            if endIdx != -1 {
                return -1, -1, fmt.Errorf("multiple ##end markers")
            }
            endIdx = i
        }
    }

    if startIdx == -1 || endIdx == -1 {
        return -1, -1, fmt.Errorf("missing ##start or ##end")
    }

    return startIdx, endIdx, nil
}
```

**Usage**:

```go
startIdx, endIdx, err := FindStartEnd(lines)
if err != nil {
    fmt.Println("ERROR: invalid data format")
    return
}
// startIdx = 1, endIdx = 4 (from example)
```

#### Step 1D: Parse Rooms

```
For each line with format "name x y":
  1. Split by spaces
  2. Extract: name (string), x (int), y (int)
  3. Create Room object
  4. Store in map: rooms[name] = &Room

Validate:
  - No duplicate room names
  - Coordinates are valid integers
  - Names don't start with 'L' or '#'
```

**Go code**:

```go
func ParseRooms(lines []string, startIdx, endIdx int) (map[string]*Room, string, string, error) {
    rooms := make(map[string]*Room)
    var startRoom, endRoom string
    nextType := ""  // "start" or "end" or ""

    for i := 1; i < len(lines); i++ {
        line := strings.TrimSpace(lines[i])

        // Skip empty lines and comments
        if line == "" || (strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##")) {
            continue
        }

        // Skip connection lines (they'll be parsed separately)
        if strings.Contains(line, "-") {
            continue
        }

        // Check for markers
        if line == "##start" {
            nextType = "start"
            continue
        }
        if line == "##end" {
            nextType = "end"
            continue
        }

        // Parse room definition: "name x y"
        fields := strings.Fields(line)
        if len(fields) < 3 {
            continue
        }

        name := fields[0]
        x, err1 := strconv.Atoi(fields[1])
        y, err2 := strconv.Atoi(fields[2])

        if err1 != nil || err2 != nil {
            return nil, "", "", fmt.Errorf("invalid room coordinates")
        }

        // Validate name
        if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
            return nil, "", "", fmt.Errorf("invalid room name: %s", name)
        }

        // Check for duplicates
        if _, exists := rooms[name]; exists {
            return nil, "", "", fmt.Errorf("duplicate room: %s", name)
        }

        // Create room
        room := NewRoom(name, x, y)
        rooms[name] = room

        // Track start and end
        if nextType == "start" {
            startRoom = name
            nextType = ""
        } else if nextType == "end" {
            endRoom = name
            nextType = ""
        }
    }

    return rooms, startRoom, endRoom, nil
}
```

**Usage**:

```go
rooms, startName, endName, err := ParseRooms(lines, startIdx, endIdx)
if err != nil {
    fmt.Println("ERROR: invalid data format")
    return
}
// rooms = map[string]*Room{"0": room0, "2": room2, ...}
// startName = "0"
// endName = "1"
```

#### Step 1E: Parse Tunnels

```
For each line with format "room1-room2":
  1. Split by "-"
  2. Get: room1 = rooms[name1], room2 = rooms[name2]
  3. Add connection: room1.connections.append(room2)
  4. Add reverse: room2.connections.append(room1)

Validate:
  - Both rooms exist
  - Not a self-loop (room1 != room2)
  - No duplicate connections
```

**Go code**:

```go
func ParseTunnels(lines []string, rooms map[string]*Room) error {
    seen := make(map[string]bool)  // Track connections

    for _, line := range lines {
        line = strings.TrimSpace(line)

        // Only process lines with hyphens (connections)
        if !strings.Contains(line, "-") {
            continue
        }

        parts := strings.Split(line, "-")
        if len(parts) != 2 {
            return fmt.Errorf("invalid connection format: %s", line)
        }

        room1Name := strings.TrimSpace(parts[0])
        room2Name := strings.TrimSpace(parts[1])

        // Validate rooms exist
        room1, exists1 := rooms[room1Name]
        room2, exists2 := rooms[room2Name]

        if !exists1 || !exists2 {
            return fmt.Errorf("connection references non-existent room: %s-%s", room1Name, room2Name)
        }

        // Check for self-loop
        if room1Name == room2Name {
            return fmt.Errorf("self-loop detected: %s", room1Name)
        }

        // Check for duplicate (normalize: smaller name first)
        key := room1Name + "-" + room2Name
        if room1Name > room2Name {
            key = room2Name + "-" + room1Name
        }

        if seen[key] {
            return fmt.Errorf("duplicate connection: %s", line)
        }
        seen[key] = true

        // Add bidirectional connection
        room1.connections = append(room1.connections, room2)
        room2.connections = append(room2.connections, room1)
    }

    return nil
}
```

**Usage**:

```go
err := ParseTunnels(lines, rooms)
if err != nil {
    fmt.Println("ERROR: invalid data format")
    return
}
// rooms["0"].connections now includes rooms["2"]
// rooms["2"].connections now includes rooms["0"]
```

#### Step 1F: Build Colony Structure

```go
type Colony struct {
    rooms map[string]*Room
    start *Room
    end   *Room
    numAnts int
}

colony.start = rooms["start_name"]
colony.end = rooms["end_name"]
colony.numAnts = numAnts
```

**Complete parsing code**:

```go
func ParseInput(lines []string) (*Colony, error) {
    // Step 1: Get ant count
    numAnts, err := ExtractAntCount(lines)
    if err != nil {
        return nil, err
    }

    // Step 2: Find start and end indices
    startIdx, endIdx, err := FindStartEnd(lines)
    if err != nil {
        return nil, err
    }

    // Step 3: Parse rooms
    rooms, startName, endName, err := ParseRooms(lines, startIdx, endIdx)
    if err != nil {
        return nil, err
    }

    // Step 4: Parse tunnels
    err = ParseTunnels(lines, rooms)
    if err != nil {
        return nil, err
    }

    // Step 5: Build colony
    colony := &Colony{
        rooms: rooms,
        start: rooms[startName],
        end: rooms[endName],
        numAnts: numAnts,
    }

    return colony, nil
}
```

**Usage in main**:

```go
func main() {
    lines, err := ReadFile(os.Args[1])
    if err != nil {
        fmt.Println("ERROR: invalid data format")
        return
    }

    colony, err := ParseInput(lines)
    if err != nil {
        fmt.Println("ERROR: invalid data format")
        return
    }

    // Now you have a fully built graph!
    fmt.Printf("Loaded %d rooms, %d ants\n", len(colony.rooms), colony.numAnts)
}
```

---

### Phase 2: Finding and Selecting Paths

#### Step 2A: Find ALL Possible Paths

This is crucial: You need to find EVERY valid path from start to end.

**Algorithm: Breadth-First Search (BFS)**

```
Pseudocode:
queue = [ [start] ]         // Start with path containing just start room
completedPaths = []

while queue is not empty:
    currentPath = queue.pop_front()
    lastRoom = currentPath.last()

    if lastRoom == end:
        completedPaths.append(currentPath)
        continue    // Don't explore further from completed path

    for each neighbor in lastRoom.connections:
        if neighbor NOT in currentPath:  // Avoid cycles
            newPath = currentPath + [neighbor]
            queue.append(newPath)

return completedPaths
```

**Go implementation**:

```go
func FindAllPaths(colony *Colony) [][]string {
    var completedPaths [][]string
    var queue [][][]*Room  // Queue of paths (each path is []*Room)

    // Initialize with starting path
    startPath := []*Room{colony.start}
    queue = append(queue, startPath)

    for len(queue) > 0 {
        // Pop first path from queue (BFS)
        currentPath := queue[0]
        queue = queue[1:]

        lastRoom := currentPath[len(currentPath)-1]

        // Check if we reached the end
        if lastRoom == colony.end {
            // Convert path to string names for output
            pathNames := make([]string, len(currentPath))
            for i, room := range currentPath {
                pathNames[i] = room.name
            }
            completedPaths = append(completedPaths, pathNames)
            continue  // Don't explore further from completed path
        }

        // Explore neighbors
        for _, neighbor := range lastRoom.connections {
            // Check if neighbor is already in path (avoid cycles)
            if !roomInPath(neighbor, currentPath) {
                // Create new path with neighbor
                newPath := make([]*Room, len(currentPath))
                copy(newPath, currentPath)
                newPath = append(newPath, neighbor)
                queue = append(queue, newPath)
            }
        }
    }

    return completedPaths
}

// Helper: Check if room is already in path
func roomInPath(room *Room, path []*Room) bool {
    for _, r := range path {
        if r == room {
            return true
        }
    }
    return false
}
```

**Usage**:

```go
allPaths := FindAllPaths(colony)

for i, path := range allPaths {
    fmt.Printf("Path %d: %v (length %d)\n", i+1, path, len(path)-1)
}
// Output:
// Path 1: [0 4 1] (length 2)
// Path 2: [0 2 3 1] (length 3)
```

**Why BFS?**

- Finds paths in order of length (shorter first)
- Complete: finds all paths
- Systematic: explores all possibilities

**Example Walkthrough**:

```
Graph: 0-2-3-1, 0-4-1

Initial: queue = [[0]]

Turn 1:
  Pop [0]
  Neighbors: [2, 4]
  Add [0,2] and [0,4] to queue

Turn 2:
  Pop [0,2]
  Neighbors: [0, 3] → Skip 0 (in path)
  Add [0,2,3] to queue

Turn 3:
  Pop [0,4]
  Neighbors: [0, 1] → Skip 0 (in path)
  1 is END → Add [0,4,1] to completedPaths

Turn 4:
  Pop [0,2,3]
  Neighbors: [2, 1] → Skip 2 (in path)
  1 is END → Add [0,2,3,1] to completedPaths

Result:
  Path 1: [0,4,1] length 2
  Path 2: [0,2,3,1] length 3
```

---

#### Step 2B: Select Optimal Paths

Found all paths. Now select which ones to use.

**Key Principle**: Select paths that maximize disjointness and efficiency.

**Algorithm: Greedy Path Selection**

```
1. Sort paths by length (shortest first)
2. usedRooms = empty set
3. selectedPaths = empty list

for each path in sorted_paths:
    if no middle rooms of path are in usedRooms:
        Add path to selectedPaths
        Mark all middle rooms as used

return selectedPaths
```

**Go implementation**:

```go
func SelectOptimalPaths(pathStrings [][]string, colony *Colony) [][]string {
    // Sort paths by length (shortest first)
    sort.Slice(pathStrings, func(i, j int) bool {
        return len(pathStrings[i]) < len(pathStrings[j])
    })

    var selectedPaths [][]string
    usedRooms := make(map[string]bool)

    for _, path := range pathStrings {
        // Check if path conflicts with already selected paths
        canSelect := true

        // Only check middle rooms (index 1 to len-2)
        // Start and end can be shared
        for i := 1; i < len(path)-1; i++ {
            if usedRooms[path[i]] {
                canSelect = false
                break
            }
        }

        if canSelect {
            // Add this path to selection
            selectedPaths = append(selectedPaths, path)

            // Mark middle rooms as used
            for i := 1; i < len(path)-1; i++ {
                usedRooms[path[i]] = true
            }
        }
    }

    return selectedPaths
}
```

**Usage**:

```go
allPaths := FindAllPaths(colony)
selectedPaths := SelectOptimalPaths(allPaths, colony)

fmt.Printf("Selected %d paths out of %d\n", len(selectedPaths), len(allPaths))
for i, path := range selectedPaths {
    fmt.Printf("Path %d: %v\n", i+1, path)
}
// Output:
// Selected 2 paths out of 2
// Path 1: [0 4 1]
// Path 2: [0 2 3 1]
```

**Why "middle rooms"?**

- Start and end can be shared (all ants start/end there)
- Only intermediate rooms create bottlenecks

Step 3: Process [0,2,3,1]
Middle rooms: [2,3]
Not used yet → SELECT
Mark 2,3 as used

Step 4: Process [0,4,2,3,1]
Middle rooms: [4,2,3]
4 already used → REJECT

Final selection: [[0,4,1], [0,2,3,1]]

```

This gives you TWO disjoint paths that ants can use in parallel!

---

### Phase 3: Simulating Ant Movement

#### Step 3A: Assign Ants to Paths

```

Ant 1 → Path 1: [0,4,1]
Ant 2 → Path 2: [0,2,3,1]
Ant 3 → Path 1: [0,4,1]
Ant 4 → Path 2: [0,2,3,1]
...

Assignment: Round-robin through paths

````go

**Go implementation**:

```
func AssignAntsToPaths(numAnts int, selectedPaths [][]string, colony *Colony) []*Ant {
    var ants []*Ant

    if len(selectedPaths) == 0 {
        return ants
    }

    // Convert string paths to Room paths
    roomPaths := make([][]*Room, len(selectedPaths))
    for i, pathNames := range selectedPaths {
        path := make([]*Room, len(pathNames))
        for j, name := range pathNames {
            path[j] = colony.rooms[name]
        }
        roomPaths[i] = path
    }

    // Assign ants round-robin through paths
    for i := 1; i <= numAnts; i++ {
        pathIndex := (i - 1) % len(roomPaths)
        ant := NewAnt(i, roomPaths[pathIndex])
        ants = append(ants, ant)
    }

    return ants
}
````

**Usage**:

```go
ants := AssignAntsToPaths(colony.numAnts, selectedPaths, colony)

fmt.Printf("Assigned %d ants to %d paths\n", len(ants), len(selectedPaths))
for i, ant := range ants {
    fmt.Printf("Ant %d: path of length %d\n", ant.id, len(ant.path)-1)
}
```

#### Step 3B: Simulate Turn-by-Turn Movement

**Go implementation**:

```go
func SimulateAnts(ants []*Ant, colony *Colony) {
    turn := 0

    for {
        // Check if all ants have arrived
        allArrived := true
        for _, ant := range ants {
            if !ant.arrived {
                allArrived = false
                break
            }
        }

        if allArrived {
            break
        }

        turn++
        var turnOutput []string

        // Try to move each ant
        for _, ant := range ants {
            // Skip if already arrived
            if ant.arrived {
                continue
            }

            // Check if ant can move
            if !ant.CanMove() {
                ant.arrived = true
                continue
            }

            // Get next room
            nextRoom := ant.path[ant.position+1]

            // Check if next room is occupied
            if nextRoom.occupied && nextRoom != colony.end {
                // Can't move, stay in current room
                continue
            }

            // Move ant
            ant.Move()
            nextRoomName := ant.CurrentRoom().name

            // Record movement
            turnOutput = append(turnOutput, fmt.Sprintf("L%d-%s", ant.id, nextRoomName))

            // Mark room as occupied (but not if it's the end room)
            if ant.CurrentRoom() != colony.end {
                ant.CurrentRoom().occupied = true
            }
        }

        // Clear occupancy for next turn (except end room and start room)
        for _, room := range colony.rooms {
            if room != colony.start && room != colony.end {
                room.occupied = false
            }
        }

        // Print turn output
        if len(turnOutput) > 0 {
            fmt.Println(strings.Join(turnOutput, " "))
        }
    }
}
```

**Critical Point**: Occupy rooms during a turn, then clear them (except end).

#### Step 3C: Output Format

For each turn, print space-separated ant movements:

```
L1-4 L2-2
L1-1 L2-3
L3-4
L3-1 L4-3
L4-1
```

Format: `L<ant_id>-<destination_room>`

The output is automatically printed in the SimulateAnts function above with: `fmt.Println(strings.Join(turnOutput, " "))`

**Complete main function**:

```go
func main() {
    if len(os.Args) != 2 {
        fmt.Println("ERROR: invalid data format")
        return
    }

    lines, err := ReadFile(os.Args[1])
    if err != nil {
        fmt.Println("ERROR: invalid data format")
        return
    }

    colony, err := ParseInput(lines)
    if err != nil {
        fmt.Println("ERROR: invalid data format")
        return
    }

    allPaths := FindAllPaths(colony)
    if len(allPaths) == 0 {
        fmt.Println("ERROR: invalid data format")
        return
    }

    selectedPaths := SelectOptimalPaths(allPaths, colony)
    ants := AssignAntsToPaths(colony.numAnts, selectedPaths, colony)

    SimulateAnts(ants, colony)
}
```

---

## Algorithm Explanations

### Algorithm 1: Breadth-First Search (BFS) for Path Finding

**Goal**: Find ALL possible paths from start to end

**Time Complexity**: O(V^V) in worst case (exponential)

- V = number of rooms
- Small graphs: fine
- Large graphs: slow

**Space Complexity**: O(V^P)

- P = number of paths
- Can explode with many alternative routes

**When to use**:

- Small graphs (< 20 nodes)
- Need all possible paths
- Path length matters

**How it works**:

1. Keep a queue of "partial paths"
2. Each iteration, expand the shortest path
3. Add new rooms only if not already in path (prevents cycles)
4. Complete paths added to result

---

### Algorithm 2: Greedy Path Selection

**Goal**: Choose set of paths that maximizes parallelization

**Time Complexity**: O(P \* R)

- P = number of paths
- R = average path length

**Correctness**: Near-optimal (not guaranteed optimal for all cases)

**Intuition**:

- Choose shortest paths first (minimize total moves)
- Avoid room conflicts (maximize disjointness)
- Room occupancy is the bottleneck

**When to use**:

- Most practical scenarios
- Sufficient performance for project requirements
- Simple to understand and implement

---

### Algorithm 3: Suurballe's Algorithm (Advanced)

**Goal**: Guaranteed optimal edge-disjoint paths

**Time Complexity**: O(V \* E) per path iteration

**Correctness**: Proven optimal by Suurballe

**How it works**:

1. Find shortest path using Bellman-Ford
2. Reverse edges in that path
3. Find next shortest path (accounting for reversed edges)
4. Merge paths where they diverge
5. Repeat until no improvement

**When to use**:

- Need proven optimal solution
- Large complex graphs
- Production-grade implementation

**Note**: Sample #1 & #3 implementations use this.

---

## Common Pitfalls & Solutions

### Pitfall 1: Not Clearing Room Occupancy Between Turns

**Problem**:

```
Turn 1: Ant A moves to Room X → Mark X occupied
Turn 2: No other ant can move to X (still thinks occupied)
```

**Solution - Code**:

```go
// WRONG - Occupancy never cleared
func SimulateAnts_Wrong(ants []*Ant, colony *Colony) {
    for {
        for _, ant := range ants {
            // ... move ant ...
            ant.CurrentRoom().occupied = true  // Never cleared!
        }
    }
}

// CORRECT - Occupancy cleared each turn
func SimulateAnts(ants []*Ant, colony *Colony) {
    for {
        var turnOutput []string

        // Move ants
        for _, ant := range ants {
            if canMove {
                ant.Move()
                turnOutput = append(turnOutput, fmt.Sprintf("L%d-%s", ant.id, ant.CurrentRoom().name))
                if ant.CurrentRoom() != colony.end {
                    ant.CurrentRoom().occupied = true
                }
            }
        }

        // CRITICAL: Clear occupancy after turn
        for _, room := range colony.rooms {
            if room != colony.start && room != colony.end {
                room.occupied = false
            }
        }

        fmt.Println(strings.Join(turnOutput, " "))
    }
}
```

### Pitfall 2: Allowing Ants to Revisit Rooms

**Problem**:

```
Path: start → A → B → A → C → end  ← CYCLE!
```

**Solution - Code**:

```go
// WRONG - Can create cycles
for _, neighbor := range room.connections {
    newPath = append(newPath, neighbor)  // No check if already visited!
    queue = append(queue, newPath)
}

// CORRECT - Check before adding
for _, neighbor := range lastRoom.connections {
    if !roomInPath(neighbor, currentPath) {  // Check if already in path
        newPath := make([]*Room, len(currentPath))
        copy(newPath, currentPath)
        newPath = append(newPath, neighbor)
        queue = append(queue, newPath)
    }
}

// Helper function
func roomInPath(room *Room, path []*Room) bool {
    for _, r := range path {
        if r == room {
            return true
        }
    }
    return false
}
```

### Pitfall 3: Not Handling the Start Room Specially

**Problem**:

```
Ant moves from start → Room X
Next turn, start room is marked occupied (blocking next ant!)
```

**Solution - Code**:

```go
// WRONG - Start room gets marked occupied
for _, room := range colony.rooms {
    room.occupied = false  // Clears start room too!
    // Next ant can't start from start room
}

// CORRECT - Skip start room
for _, room := range colony.rooms {
    if room != colony.start && room != colony.end {
        room.occupied = false
    }
}

// Also in movement:
if ant.CurrentRoom() != colony.end {
    ant.CurrentRoom().occupied = true  // But NOT if at start
}
```

```

**Solution**:

```

Never mark start room as occupied:
for each room except start and end:
room.occupied = false

```

### Pitfall 4: Output Format

**Wrong**:

```

L1-room1 L2-room2
L1-room2 L3-room3 ← Ant 1 moved twice in one turn!

```

**Correct**:

```

L1-room1 L2-room2
L1-room2 L3-room3 ← Different ant IDs or same ant in DIFFERENT turns

```

Each ant moves (or doesn't move) once per turn.

### Pitfall 5: Misunderstanding Tunnel Bidirectionality

**Wrong**:

```

0-1 ← Means only 0→1 connection

```

**Correct**:

```

0-1 ← Means 0↔1 (bidirectional)
Ant in 0 can move to 1
Ant in 1 can move to 0

```

### Pitfall 6: Not Validating Input

**Must Check**:

- Number of ants > 0
- Exactly one start and one end
- No duplicate rooms
- No duplicate connections
- No self-loops
- Both rooms in tunnel exist
- Room names don't start with 'L' or '#'

---

## Testing Strategy

### Test Case 1: Minimal Example

```

2
##start
0 0 0
##end
1 1 1
0-1

```

Expected output:

```

L1-1 L2-1

```

(Both ants move on first turn)

---

### Test Case 2: Sequential Path

```

2
##start
0 0 0
1 1 1
##end
2 2 2
0-1
1-2

```

Expected output:

```

L1-1 L2-1
L1-2
L2-2

```

---

### Test Case 3: Multiple Disjoint Paths

```

4
##start
0 0 0
2 2 0
3 3 0
##end
1 1 0
0-2
0-3
2-1
3-1

```

Expected output:

```

L1-2 L2-3 L3-2 L4-3
L1-1 L2-1 L3-3 L4-2
L3-1
L4-1

````

---

### Validation Checklist

For each test, verify:

- [ ] All ants start in start room (turn 1)
- [ ] All ants end in end room (final turn)
- [ ] No ant appears twice in one turn
- [ ] No two ants in same room (except start/end)
- [ ] All movements follow `Lx-y` format
- [ ] Ants only move to connected rooms
- [ ] Output matches expected turn count

---

## Putting It All Together: Main Function

```go
func main() {
    // Step 1: Parse input
    content := readFile(filename)
    lines := parseLines(content)

    numAnts, startRoom, endRoom, colony := parseInput(lines)

    // Step 2: Validate
    if !validateInput(numAnts, startRoom, endRoom, colony) {
        fmt.Println("ERROR")
        return
    }

    // Step 3: Find paths
    allPaths := findAllPaths(colony, startRoom, endRoom)
    if len(allPaths) == 0 {
        fmt.Println("ERROR")
        return
    }

    // Step 4: Select optimal paths
    selectedPaths := selectOptimalPaths(allPaths)

    // Step 5: Simulate
    simulate(numAnts, selectedPaths, colony)
}
````

---

## Key Takeaways

1. **The Problem**: Route N ants through a graph in minimum turns
2. **The Solution**: Find multiple non-interfering paths and parallelize
3. **The Approach**:
   - Parse input into graph structure
   - Find all possible paths (BFS)
   - Select best path combination (Greedy)
   - Simulate ant movement turn-by-turn
   - Output in required format

4. **The Constraints**:
   - One ant per room per turn (except start/end)
   - No cycles in paths
   - One tunnel use per turn
   - Simultaneous movement

5. **The Challenge**: Optimizing path selection to minimize total traversal time

---

## Final Advice

**Start Simple**:

1. Get file parsing working
2. Build graph data structure
3. Find one path (DFS/BFS)
4. Simulate ants on that path
5. Gradually add multiple path selection
6. Optimize path selection algorithm


Hasan's corner (main source)
https://www.geeksforgeeks.org/dsa/print-all-shortest-paths-between-given-source-and-destination-in-an-undirected-graph/
https://www.geeksforgeeks.org/dsa/breadth-first-search-or-bfs-for-a-graph/
https://github.com/sohamkamani/golang-graph-traversal/blob/master/bfs.go