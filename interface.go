package toposort

// A type that satisfies toposort.Interface can be contained in a graph node,
// and, therefore, organized and produced in topological ordering.
type Interface interface {
	Name() string
}

type Interfaces []Interface

func (e Interfaces) Names() []string {
	s := make([]string, len(e))
	for i := range e {
		s[i] = e[i].Name()
	}
	return s
}

type stringElement string

func (e stringElement) Name() string {
	return string(e)
}

type stringer interface {
	String() string
}
