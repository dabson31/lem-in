package graphs

// this approach is only good if i am going to use BFS or DFS

type room struct {
	int
	name 			string
	x_coord			int
	Y_coord			int
	connections		[]*room
	occupied		bool
}

type ant struct {
	antID			int
	paths			[]*room
	intersection 	bool
	queuedAnts		int
	pathLength		int
}

type tunnel struct {
	
}

