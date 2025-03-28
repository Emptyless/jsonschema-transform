package parse

import (
	"github.com/Emptyless/jsonschema-transform/domain"
	"github.com/kaptinlin/jsonschema"
)

// DepthMap w.r.t. the provided jsonschema.Schema roots that are directly parsed from the glob patterns. Each directly
// referenced file (e.g. some/file/schema.json) must have a depth of 0. For each not directly referenced file the distance
// is the shortest distance from any root to that file using the classical Dijkstra's shortest-path algorithm.
func DepthMap(roots []*jsonschema.Schema, classes []*domain.Class, relations []*domain.Relation) map[*domain.Class]int {
	graph := nodes{}
	for _, class := range classes {
		graph[class] = &node{distance: unvisited, class: class, edges: nil}
	}
	graph.SetEdges(relations)

	depthMap := map[*domain.Class]int{}
	for _, root := range roots {
		// find class that belongs to root
		var class *domain.Class
		for _, c := range classes {
			if c.Schema == root {
				class = c
				break
			}
		}

		if class == nil {
			panic("class belonging to schema not found which is not possible")
		}

		n := graph[class]
		distances := graph.Distance(n)
		for k, v := range distances {
			if _, ok := depthMap[k]; !ok {
				depthMap[k] = v
				continue
			}

			if depthMap[k] > v {
				depthMap[k] = v
			}
		}

		graph.Reset()
	}

	return depthMap
}

// unvisited nodes have the MAX_INT distance value
const unvisited = int(^uint(0) >> 1)

// nodes is the underlying graph datastructure backed by a map
type nodes map[*domain.Class]*node

// SetEdges on the graph using the domain.Relation(s)
func (graph nodes) SetEdges(relations []*domain.Relation) {
	for class, n := range graph {
		for _, relation := range relations {
			if relation.From == class {
				n.edges = append(n.edges, graph[relation.To])
			}

			if relation.To == class {
				n.edges = append(n.edges, graph[relation.From])
			}
		}
	}
}

// Distance from node to domain.Class'es it can reach via edges
func (graph nodes) Distance(from *node) map[*domain.Class]int {
	from.distance = 0
	from.Distance()

	res := map[*domain.Class]int{}
	for k, v := range graph {
		res[k] = v.distance
	}

	return res
}

// Reset all distances on graph to unvisited
func (graph nodes) Reset() {
	if len(graph) == 0 {
		return
	}

	for _, v := range graph {
		v.distance = unvisited
	}

	return
}

// node structure to build graph and track distances
type node struct {
	distance int
	class    *domain.Class
	edges    []*node
}

// Distance from some node to all nodes it can find in the graph. Each edge traversed has a distance of 1
func (n *node) Distance() {
	for _, edge := range n.edges {
		if n.distance+1 > edge.distance {
			continue
		}

		edge.distance = n.distance + 1
		edge.Distance()
	}
}
