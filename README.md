# Lem-in: Ant Farm Path Optimization

## Description

A graph-based optimization problem where n ants must traverse from a start room to an end room through an interconnected colony in the minimum number of moves.

## Problem Statement

- Route n ants from `##start` to `##end` through rooms connected by tunnels
- Each room (except end) holds only one ant per turn
- Each tunnel is usable once per turn
- Minimize total traversal time by distributing ants across multiple disjoint paths
- The graphs are weightless and undirected; the tunnels are bidirectional.

## Input Format

```
number_of_ants
##start
room_name x y
room_name x y
##end
room_name x y
room1-room2
room1-room3
...
```

## Output Format

```
Lx-y Lz-w Lr-o ...
Lx-y Lz-w ...
```

Where `Lx-y` represents ant x moving to room y. One line per turn.

## Key Concepts

- **Path Disjointness:** Routes should share no edges to maximize parallelization
- **Room Occupancy:** Only one ant per room per turn (except start/end)
- **Turnwise Movement:** All eligible ants move simultaneously each turn
- **Path Optimization:** Multiple paths reduce total move count through parallelization

## Running

```bash
go run . testing/examples/example00.txt
```

## Alternative way to run the test files
```bash
cd miscellaneous/scripts/
sh script.sh
```

## Validation Checklist

- All ants start in `##start` and end in `##end`
- No ant moves twice per turn
- No two ants occupy same room per turn (except start/end)
- No self-loops or duplicate rooms
- No duplicate connections
- Room names cannot start with `L` or `#`
- All movements follow `Lx-y` format
- Detailed error messages for invalid input

## Credits

- Ridha Hasan (rihasan)
- Hasan Ali (hasanali4)
- Ali Watheqi (awatheqi)
