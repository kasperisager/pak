package asset

type (
	Graph struct {
		nodes map[Asset]bool
		edges map[Asset]edges
	}

	edges struct {
		incoming relations
		outgoing relations
	}

	relations map[Asset]Relation
)

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
	if g.Has(asset) {
		return false
	}

	g.nodes[asset] = true
	g.edges[asset] = edges{
		incoming: make(relations),
		outgoing: make(relations),
	}

	return true
}

func (g *Graph) Remove(asset Asset) bool {
	if !g.Has(asset) {
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

func (g *Graph) Relation(from Asset, to Asset, relation Relation) bool {
	if !g.Has(from) || !g.Has(to) {
		return false
	}

	g.edges[from].outgoing[to] = relation
	g.edges[to].incoming[from] = relation

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

func (g *Graph) Incoming(asset Asset) ([]Asset, []Relation, bool) {
	if edges, ok := g.edges[asset]; ok {
		assets := make([]Asset, 0, len(edges.incoming))
		relations := make([]Relation, 0, len(edges.incoming))

		for asset, relation := range edges.incoming {
			assets = append(assets, asset)
			relations = append(relations, relation)
		}

		return assets, relations, true
	}

	return nil, nil, false
}

func (g *Graph) Indegree(asset Asset) (int, bool) {
	if edges, ok := g.edges[asset]; ok {
		return len(edges.incoming), true
	}

	return 0, false
}

func (g *Graph) Outgoing(asset Asset) ([]Asset, []Relation, bool) {
	if edges, ok := g.edges[asset]; ok {
		assets := make([]Asset, 0, len(edges.outgoing))
		relations := make([]Relation, 0, len(edges.outgoing))

		for asset, relation := range edges.outgoing {
			assets = append(assets, asset)
			relations = append(relations, relation)
		}

		return assets, relations, true
	}

	return nil, nil, false
}

func (g *Graph) Outdegree(asset Asset) (int, bool) {
	if edges, ok := g.edges[asset]; ok {
		return len(edges.outgoing), true
	}

	return 0, false
}

func (g *Graph) Lookup(query Query) (Asset, bool) {
	for asset, _ := range g.nodes {
		if query(asset) {
			return asset, true
		}
	}

	return nil, false
}

func (g *Graph) Merge(target Asset, source Asset) bool {
	if !g.Has(target) || !g.Has(source) {
		return false
	}

	assets, relations, _ := g.Outgoing(source)

	for i, source := range assets {
		g.Relation(target, source, relations[i])
	}

	g.Remove(source)

	return true
}
