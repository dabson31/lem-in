package anthill

// this approach is only good if i am going to use BFS or DFS

type room struct {
	name        string
	x_coord     int
	Y_coord     int
	connections []*room
	occupied    bool
}

type Ant struct {
	antID           int
	currentRoom     string
	pathIndex       int
	currentPosition int
	reachedEnd      bool
}
