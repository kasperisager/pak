package asset

import (
	"net/url"
)

type edges struct {
	incoming map[Asset]bool
	outgoing map[Asset]bool
}

type Graph struct {
	nodes map[Asset]bool
	edges map[Asset]edges
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[Asset]bool),
		edges: make(map[Asset]edges),
	}
}

func (g *Graph) Size() int {
	return len(g.nodes)
}

func (g *Graph) Assets() []Asset {
	assets := make([]Asset, 0, len(g.nodes))

	for asset, _ := range g.nodes {
		assets = append(assets, asset)
	}

	return assets
}

func (g *Graph) Has(asset Asset) bool {
	return g.nodes[asset]
}

func (g *Graph) Add(asset Asset) bool {
	if g.nodes[asset] {
		return false
	}

	g.nodes[asset] = true
	g.edges[asset] = edges{
		incoming: make(map[Asset]bool),
		outgoing: make(map[Asset]bool),
	}

	return true
}

func (g *Graph) Remove(asset Asset) bool {
	if !g.nodes[asset] {
		return false
	}

	for ref, _ := range g.edges[asset].incoming {
		delete(g.edges[ref].outgoing, asset)
	}

	for ref, _ := range g.edges[asset].outgoing {
		delete(g.edges[ref].incoming, asset)
	}

	delete(g.nodes, asset)
	delete(g.edges, asset)

	return true
}

func (g *Graph) Reference(from Asset, to Asset) bool {
	if !g.nodes[from] || !g.nodes[to] {
		return false
	}

	g.edges[from].outgoing[to] = true
	g.edges[to].incoming[from] = true

	return true
}

func (g *Graph) Roots() []Asset {
	roots := make([]Asset, 0)

	for asset, _ := range g.nodes {
		if len(g.edges[asset].incoming) == 0 {
			roots = append(roots, asset)
		}
	}

	return roots
}

func (g *Graph) Leaves() []Asset {
	leaves := make([]Asset, 0)

	for asset, _ := range g.nodes {
		if len(g.edges[asset].outgoing) == 0 {
			leaves = append(leaves, asset)
		}
	}

	return leaves
}

func (g *Graph) Incoming(asset Asset) ([]Asset, bool) {
	if edges, ok := g.edges[asset]; ok {
		incoming := make([]Asset, 0, len(edges.incoming))

		for asset, _ := range edges.incoming {
			incoming = append(incoming, asset)
		}

		return incoming, true
	}

	return nil, false
}

func (g *Graph) Indegree(asset Asset) (int, bool) {
	if edges, ok := g.edges[asset]; ok {
		return len(edges.incoming), true
	}

	return 0, false
}

func (g *Graph) Outgoing(asset Asset) ([]Asset, bool) {
	if edges, ok := g.edges[asset]; ok {
		outgoing := make([]Asset, 0, len(edges.outgoing))

		for asset, _ := range edges.outgoing {
			outgoing = append(outgoing, asset)
		}

		return outgoing, true
	}

	return nil, false
}

func (g *Graph) Outdegree(asset Asset) (int, bool) {
	if edges, ok := g.edges[asset]; ok {
		return len(edges.outgoing), true
	}

	return 0, false
}

func (g *Graph) Lookup(u *url.URL) (Asset, bool) {
	for asset, _ := range g.nodes {
		v := asset.URL()

		if v.Scheme == u.Scheme && v.Host == v.Host && v.Path == u.Path {
			return asset, true
		}
	}

	return nil, false
}

func (g *Graph) Merge(target Asset, source Asset) bool {
	if !g.Has(target) || !g.Has(source) {
		return false
	}

	references, _ := g.Outgoing(source)

	for _, reference := range references {
		g.Reference(target, reference)
	}

	g.Remove(source)

	return true
}
