package asset

import (
	"net/url"
)

type references map[Asset]Reference

type edges struct {
	incoming references
	outgoing references
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
		incoming: make(references),
		outgoing: make(references),
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

func (g *Graph) Reference(from Asset, to Asset, reference Reference) bool {
	if !g.nodes[from] || !g.nodes[to] {
		return false
	}

	g.edges[from].outgoing[to] = reference
	g.edges[to].incoming[from] = reference

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

func (g *Graph) Incoming(asset Asset) ([]Asset, []Reference, bool) {
	if edges, ok := g.edges[asset]; ok {
		assets := make([]Asset, 0, len(edges.incoming))
		references := make([]Reference, 0, len(edges.incoming))

		for asset, reference := range edges.incoming {
			assets = append(assets, asset)
			references = append(references, reference)
		}

		return assets, references, true
	}

	return nil, nil, false
}

func (g *Graph) Indegree(asset Asset) (int, bool) {
	if edges, ok := g.edges[asset]; ok {
		return len(edges.incoming), true
	}

	return 0, false
}

func (g *Graph) Outgoing(asset Asset) ([]Asset, []Reference, bool) {
	if edges, ok := g.edges[asset]; ok {
		assets := make([]Asset, 0, len(edges.outgoing))
		references := make([]Reference, 0, len(edges.outgoing))

		for asset, reference := range edges.outgoing {
			assets = append(assets, asset)
			references = append(references, reference)
		}

		return assets, references, true
	}

	return nil, nil, false
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

	assets, references, _ := g.Outgoing(source)

	for i, source := range assets {
		g.Reference(target, source, references[i])
	}

	g.Remove(source)

	return true
}
