package graphs

// this approach is only good if i am going to use BFS or DFS

type room struct {
	parent		[](*room) // for incoming connections
	children	[](*room) // for outgoing connections
	name		string
	occupied	bool
	queue		int
}

type ant struct {
	antID			int
	paths			[]*room
	intersection 	bool
	queuedAnts		int
	pathLength		int
}

var ( 
	firstRoom *room
	lastRoom *room
	//! i stil need to add more variables later
)

