package graphs

type Graph struct {
	// it needs a list to hold (point at) the Vertices
	Vertices []*Vertex
}

type Vertex struct {
	Key int
	adjacent []*Vertex
}

func (g *Graph) AddVertex (k int) {
	g.Vertices = append(g.Vertices, &Vertex{Key: k})
}

