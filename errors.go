package toposort

type GraphError int

const (
	ErrNodeExists GraphError = iota
	ErrNodeNotFound
	ErrUnsupportedType
	ErrCycle
	ErrUnknown
)

func (err GraphError) Error() string {
	switch err {
	case ErrNodeExists:
		return "toposort: node already exists"
	case ErrNodeNotFound:
		return "toposort: node not found"
	case ErrUnsupportedType:
		return "toposort: unsupported type"
	case ErrCycle:
		return "toposort: graph contains a cycle"
	case ErrUnknown:
		return "toposort: unknown error"
	default:
		return ErrUnknown.Error()
	}
}
