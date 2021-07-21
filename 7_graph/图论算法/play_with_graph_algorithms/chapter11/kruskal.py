from play_with_graph_algorithms.chapter11.weighted_edge import WeightedEdge
from play_with_graph_algorithms.chapter11.weighted_graph import WeightedGraph
from play_with_graph_algorithms.chapter11.uf import UF
from play_with_graph_algorithms.chapter04.cc import CC


class Kruskal:

    def __init__(self, G):
        self._G = G
        self._mst = []

        cc = CC(G)
        if (cc.ccount > 1):
            return

        # Kruskal
        edges = []
        for v in range(G.V):
            for w in G.adj(v):
                if v < w:
                    edges.append(WeightedEdge(v, w, G.get_weight(v, w)))

        edges = sorted(edges, key=lambda x: x.weight)
        uf = UF(G.V)
        for edge in edges:
            v = edge.v
            w = edge.w
            weight = edge.weight
            if not uf.is_connected(v, w):
                self._mst.append(edge)
                uf.union_elements(v, w)

    def result(self):
        return self._mst


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter11/g.txt'
    g = WeightedGraph(filename)
    kruskal = Kruskal(g)
    print(kruskal.result())
