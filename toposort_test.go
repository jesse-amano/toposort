package toposort

import (
	"fmt"
	"testing"
)

func index(s []Interface, v string) int {
	for i, s := range s {
		if s.Name() == v {
			return i
		}
	}
	return -1
}

type Edge struct {
	From string
	To   string
}

func TestDuplicatedNode(t *testing.T) {
	graph := NewGraph(2)
	graph.AddNode("a")
	if err := graph.AddNode("a"); err != ErrNodeExists {
		t.Errorf("not raising duplicated node error")
	}
}

func TestRemoveNotExistEdge(t *testing.T) {
	graph := NewGraph(0)
	if err := graph.RemoveEdge("a", "b"); err != ErrNodeNotFound {
		t.Errorf("not raising not exist edge error")
	}
}

func TestWikipedia(t *testing.T) {
	graph := NewGraph(8)
	graph.AddNodes("2", "3", "5", "7", "8", "9", "10", "11")

	edges := []Edge{
		{"7", "8"},
		{"7", "11"},

		{"5", "11"},

		{"3", "8"},
		{"3", "10"},

		{"11", "2"},
		{"11", "9"},
		{"11", "10"},

		{"8", "9"},
	}

	for _, e := range edges {
		graph.AddEdge(e.From, e.To)
	}

	result, err := graph.Toposort()
	if err != nil {
		t.Errorf("error sorting valid DAG: %v", err)
	}

	for _, e := range edges {
		if i, j := index(result, e.From), index(result, e.To); i > j {
			t.Errorf("dependency failed: not satisfy %v(%v) > %v(%v)", e.From, i, e.To, j)
		}
	}
}

func TestCycle(t *testing.T) {
	graph := NewGraph(3)
	graph.AddNodes("1", "2", "3")

	graph.AddEdge("1", "2")
	graph.AddEdge("2", "3")
	graph.AddEdge("3", "1")

	_, err := graph.Toposort()
	if err == nil {
		t.Error("closed path not detected in closed pathed graph")
	}
}

func TestStructured(t *testing.T) {
	baskets := []basket{
		{count: 2, fruit: "bananas"},
		{count: 3, fruit: "cantaloupes"},
		{count: 5, fruit: "elderberries"},
		{count: 7, fruit: "grapes"},
		{count: 8, fruit: "honeydew melons"},
		{count: 9, fruit: "idared apples"},
		{count: 10, fruit: "jackfruit"},
		{count: 11, fruit: "kumquat"},
	}

	graph := NewGraph(len(baskets))
	for _, b := range baskets {
		graph.AddNode(b)
	}

	edges := []Edge{
		{"7 grapes", "8 honeydew melons"},
		{"7 grapes", "11 kumquat"},

		{"5 elderberries", "11 kumquat"},

		{"3 cantaloupes", "8 honeydew melons"},
		{"3 cantaloupes", "10 jackfruit"},

		{"11 kumquat", "2 bananas"},
		{"11 kumquat", "9 idared apples"},
		{"11 kumquat", "10 jackfruit"},

		{"8 honeydew melons", "9 idared apples"},
	}

	for _, e := range edges {
		graph.AddEdge(e.From, e.To)
	}

	result, err := graph.Toposort()
	if err != nil {
		t.Errorf("error sorting valid DAG: %v", err)
	}

	for _, e := range edges {
		if i, j := index(result, e.From), index(result, e.To); i > j {
			t.Errorf("dependency failed: not satisfy %v(%v) > %v(%v)", e.From, i, e.To, j)
		}
	}
}

type basket struct {
	count int
	fruit string
}

func (b basket) Name() string {
	return fmt.Sprintf("%d %s", b.count, b.fruit)
}
