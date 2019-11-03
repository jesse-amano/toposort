package toposort

// Graph represents a directed graph.
type Graph struct {
	nodes   []string
	outputs map[string]map[string]int
	inputs  map[string]int
	objects map[string]Interface
}

// NewGraph returns a new graph with an initial capacity.
func NewGraph(cap int) *Graph {
	return &Graph{
		nodes:   make([]string, 0, cap),
		inputs:  make(map[string]int),
		outputs: make(map[string]map[string]int),
		objects: make(map[string]Interface, cap),
	}
}

// AddNode adds a single node to the graph containing an element.
// If element does not satisfy toposort.Interface but is already
// a string or stringer, it will be converted to a toposort.Interface
// value whose Name is equal to the string value of element.
func (g *Graph) AddNode(element interface{}) bool {
	if el, ok := element.(Interface); ok {
		return g.addNode(el)
	}
	if str, ok := element.(string); ok {
		return g.addNode(stringElement(str))
	}
	if str, ok := element.(stringer); ok {
		return g.addNode(stringElement(str.String()))
	}
	return false
}

func (g *Graph) addNode(element Interface) bool {
	name := element.Name()
	if _, ok := g.outputs[name]; ok {
		return false
	}

	g.objects[name] = element
	g.nodes = append(g.nodes, name)

	g.outputs[name] = make(map[string]int)
	g.inputs[name] = 0
	return true
}

// AddNodes is a convenience method to add multiple nodes at once.
func (g *Graph) AddNodes(elements ...interface{}) bool {
	for _, e := range elements {
		if ok := g.AddNode(e); !ok {
			return false
		}
	}
	return true
}

// AddEdge creates a directed edge from one node to another.
// The first edge will be required to appear before the second
// when the graph is traversed in topological order.
func (g *Graph) AddEdge(from, to string) bool {
	m, ok := g.outputs[from]
	if !ok {
		return false
	}

	m[to] = len(m) + 1
	g.inputs[to]++

	return true
}

func (g *Graph) unsafeRemoveEdge(from, to string) {
	delete(g.outputs[from], to)
	g.inputs[to]--
}

// RemoveEdge removes an edge from one node to another.
func (g *Graph) RemoveEdge(from, to string) bool {
	if _, ok := g.outputs[from]; !ok {
		return false
	}
	g.unsafeRemoveEdge(from, to)
	return true
}

// Toposort returns a slice representing a topological ordering
// of the nodes in the graph.
func (g *Graph) Toposort() ([]Interface, bool) {
	names, ok := g.toposort()
	elements := make([]Interface, len(names))
	if !ok {
		return elements, false
	}
	for i := range names {
		elements[i], ok = g.objects[names[i]]
		if !ok {
			return elements, false
		}
	}
	return elements, true
}

func (g *Graph) toposort() ([]string, bool) {
	L := make([]string, 0, len(g.nodes))
	S := make([]string, 0, len(g.nodes))

	for _, n := range g.nodes {
		if g.inputs[n] == 0 {
			S = append(S, n)
		}
	}

	for len(S) > 0 {
		var n string
		n, S = S[0], S[1:]
		L = append(L, n)

		ms := make([]string, len(g.outputs[n]))
		for m, i := range g.outputs[n] {
			ms[i-1] = m
		}

		for _, m := range ms {
			g.unsafeRemoveEdge(n, m)

			if g.inputs[m] == 0 {
				S = append(S, m)
			}
		}
	}

	N := 0
	for _, v := range g.inputs {
		N += v
	}

	if N > 0 {
		return L, false
	}

	return L, true
}
