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

	relations map[Relation]Asset
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

	for asset := range g.nodes {
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

func (g *Graph) Delete(asset Asset) bool {
	if !g.Has(asset) {
		return false
	}

	for _, edge := range g.edges[asset].incoming {
		for relation, found := range g.edges[edge].outgoing {
			if found == asset {
				delete(g.edges[edge].outgoing, relation)
			}
		}
	}

	for _, edge := range g.edges[asset].outgoing {
		for relation, found := range g.edges[edge].incoming {
			if found == asset {
				delete(g.edges[edge].incoming, relation)
			}
		}
	}

	delete(g.nodes, asset)
	delete(g.edges, asset)

	return true
}

func (g *Graph) Relate(from Asset, to Asset, relation Relation) bool {
	if !g.Has(from) || !g.Has(to) {
		return false
	}

	g.edges[from].outgoing[relation] = to
	g.edges[to].incoming[relation] = from

	return true
}

func (g *Graph) Relation(from Asset, to Asset) (Relation, bool) {
	if edges, ok := g.Outgoing(from); ok {
		for relation, asset := range edges {
			if to == asset {
				return relation, true
			}
		}
	}

	return nil, false
}

func (g *Graph) Roots() []Asset {
	roots := make([]Asset, 0)

	for asset := range g.nodes {
		indegree, _ := g.Indegree(asset)

		if indegree == 0 {
			roots = append(roots, asset)
		}
	}

	return roots
}

func (g *Graph) Leaves() []Asset {
	leaves := make([]Asset, 0)

	for asset := range g.nodes {
		outdegree, _ := g.Outdegree(asset)

		if outdegree == 0 {
			leaves = append(leaves, asset)
		}
	}

	return leaves
}

func (g *Graph) Incoming(asset Asset) (map[Relation]Asset, bool) {
	if edges, ok := g.edges[asset]; ok {
		result := make(map[Relation]Asset, len(edges.incoming))

		for relation, asset := range edges.incoming {
			result[relation] = asset
		}

		return result, true
	}

	return nil, false
}

func (g *Graph) Indegree(asset Asset) (int, bool) {
	if edges, ok := g.edges[asset]; ok {
		return len(edges.incoming), true
	}

	return 0, false
}

func (g *Graph) Outgoing(asset Asset) (map[Relation]Asset, bool) {
	if edges, ok := g.edges[asset]; ok {
		result := make(map[Relation]Asset, len(edges.outgoing))

		for relation, asset := range edges.outgoing {
			result[relation] = asset
		}

		return result, true
	}

	return nil, false
}

func (g *Graph) Outdegree(asset Asset) (int, bool) {
	if edges, ok := g.edges[asset]; ok {
		return len(edges.outgoing), true
	}

	return 0, false
}

func (g *Graph) Lookup(query Query) (Asset, bool) {
	for asset := range g.nodes {
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

	relation, ok := g.Relation(target, source)

	if !ok {
		return false
	}

	ok = target.Merge(source, relation)

	if !ok {
		return false
	}

	if edges, ok := g.Outgoing(source); ok {
		for relation, related := range edges {
			if related != target {
				switch relation := relation.(type) {
				case Reference:
					relation.Rewrite(
						rebase(
							relation.URL(),
							source.URL(),
							target.URL(),
						),
					)
				}

				g.Relate(target, related, relation)
			}
		}
	}

	if edges, ok := g.Incoming(source); ok {
		for relation, related := range edges {
			if related != target {
				switch relation := relation.(type) {
				case Reference:
					relation.Rewrite(
						rewrite(
							related.URL(),
							relation.URL(),
							target.URL(),
						),
					)
				}

				g.Relate(related, target, relation)
			}
		}
	}

	g.Delete(source)

	return true
}
