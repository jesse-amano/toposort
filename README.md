go-toposort
==
[![GoDoc](https://godoc.org/github.com/jesse-amano/toposort?status.svg)](https://godoc.org/github.com/jesse-amano/toposort)

Deterministic topological sort for Go. Forked from [philopon/go-toposort](https://github.com/philopon/go-toposort).

This fork of the package is extended to accommodate custom elements. If you only need to sort strings or other simple data types, check out other forks. Custom element support is achieved with safe type assertions, so performance is not maximized, but for small dependency graphs it should still be okay.

License
--
MIT

Example
--

```.go
package main

import (
	"fmt"

	toposort "github.com/jesse-amano/toposort"
)

func main() {
	graph := toposort.NewGraph(8)
	graph.AddNodes("2", "3", "5", "7", "8", "9", "10", "11")

	graph.AddEdge("7", "8")
	graph.AddEdge("7", "11")

	graph.AddEdge("5", "11")

	graph.AddEdge("3", "8")
	graph.AddEdge("3", "10")

	graph.AddEdge("11", "2")
	graph.AddEdge("11", "9")
	graph.AddEdge("11", "10")

	graph.AddEdge("8", "9")

	result, ok := graph.Toposort()
	if !ok {
		panic("cycle detected")
	}

	fmt.Println(result)
}
```

```
[3 5 7 8 11 2 9 10]
```

