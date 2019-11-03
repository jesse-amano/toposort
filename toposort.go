package toposort

// Graph represents a directed graph. It is not safe
// for use by concurrent goroutines.
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
func (g *Graph) AddNode(element interface{}) error {
	g.initialize()
	if el, ok := element.(Interface); ok {
		return g.addNode(el)
	}
	if str, ok := element.(string); ok {
		return g.addNode(stringElement(str))
	}
	if str, ok := element.(stringer); ok {
		return g.addNode(stringElement(str.String()))
	}
	return ErrUnsupportedType
}

func (g *Graph) addNode(element Interface) error {
	name := element.Name()
	if _, ok := g.outputs[name]; ok {
		return ErrNodeExists
	}

	g.objects[name] = element
	g.nodes = append(g.nodes, name)

	g.outputs[name] = make(map[string]int)
	g.inputs[name] = 0
	return nil
}

// AddNodes is a convenience method to add multiple nodes at once.
func (g *Graph) AddNodes(elements ...interface{}) error {
	for _, e := range elements {
		if err := g.AddNode(e); err != nil {
			return err
		}
	}
	return nil
}

// AddEdge creates a directed edge from one node to another.
// The first edge will be required to appear before the second
// when the graph is traversed in topological order.
func (g *Graph) AddEdge(from, to string) error {
	g.initialize()
	m, ok := g.outputs[from]
	if !ok {
		return ErrNodeNotFound
	}
	_, ok = g.objects[to]
	if !ok {
		return ErrNodeNotFound
	}

	m[to] = len(m) + 1
	g.inputs[to]++

	return nil
}

func (g *Graph) unsafeRemoveEdge(from, to string) {
	delete(g.outputs[from], to)
	g.inputs[to]--
}

// RemoveEdge removes an edge from one node to another.
func (g *Graph) RemoveEdge(from, to string) error {
	g.initialize()
	if _, ok := g.objects[to]; !ok {
		return ErrNodeNotFound
	}
	if m, ok := g.outputs[from]; !ok {
		return ErrNodeNotFound
	} else if _, ok = m[to]; !ok {
		return ErrEdgeNotFound
	}
	g.unsafeRemoveEdge(from, to)
	return nil
}

// Toposort returns a slice representing a topological ordering
// of the nodes in the graph.
func (g *Graph) Toposort() ([]Interface, error) {
	g.initialize()
	return clone(g).DestructiveToposort()
}

// DestructiveToposort returns a slice representing a topological ordering
// of the nodes in the graph. It is significantly faster than Toposort but
// alters the structure of the underlying graph. Call Toposort instead if
// you want to reuse the graph structure.
func (g *Graph) DestructiveToposort() ([]Interface, error) {
	g.initialize()
	names, err := g.toposort()
	elements := make([]Interface, len(names))
	if err != nil {
		return elements, err
	}
	var ok bool
	for i := range names {
		elements[i], ok = g.objects[names[i]]
		if !ok {
			return elements, ErrNodeNotFound
		}
	}
	return elements, nil
}

func (g *Graph) toposort() ([]string, error) {
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
		return L, ErrCycle
	}

	return L, nil
}

func (g *Graph) initialize() {
	if g.objects != nil {
		return
	}
	g.nodes = make([]string, 0)
	g.inputs = make(map[string]int)
	g.outputs = make(map[string]map[string]int)
	g.objects = make(map[string]Interface)
}

func clone(g *Graph) *Graph {
	if g == nil {
		return nil
	}
	var h Graph
	h.nodes = make([]string, len(g.nodes))
	copy(h.nodes, g.nodes)
	h.inputs = make(map[string]int, len(g.inputs))
	for k, v := range g.inputs {
		h.inputs[k] = v
	}
	h.outputs = make(map[string]map[string]int, len(g.outputs))
	for k, v := range g.outputs {
		h.outputs[k] = make(map[string]int, len(v))
		for ik, iv := range v {
			h.outputs[k][ik] = iv
		}
	}
	h.objects = make(map[string]Interface, len(g.objects))
	for k, v := range g.objects {
		h.objects[k] = v
	}
	return &h
}
