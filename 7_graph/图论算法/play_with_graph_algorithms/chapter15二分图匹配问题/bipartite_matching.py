from play_with_graph_algorithms.chapter04.bi_partite_detection import BiPartiteDetection
from play_with_graph_algorithms.chapter13.weighted_graph import WeightedGraph
from play_with_graph_algorithms.chapter13.graph import Graph
from play_with_graph_algorithms.chapter14.max_flow import MaxFlow


class BipartiteMatching:

    def __init__(self, G):
        bd = BiPartiteDetection(G)
        if bd.is_bi_partite() is False:
            raise ValueError('BipartiteMatching only works for bipartite graph')
        self._G = G
        colors = bd.colors()
        network = WeightedGraph(empty_graph=True, directed=True, V=G.V + 2)

        for v in range(self._G.V):
            # 将源点和汇点分别和二分图中不同的部分相连
            if colors[v] == 0:
                network.add_edge(G.V, v, 1)
            else:
                network.add_edge(v, G.V + 1, 1)

            # 将二分图两个部分相连
            for w in self._G.adj(v):
                if v < w:
                    if colors[v] == 0:
                        network.add_edge(v, w, 1)
                    else:
                        network.add_edge(w, v, 1)

        maxflow = MaxFlow(
            network,
            self._G.V,
            self._G.V + 1,
            empty_graph=True,
            directed=True,
            V=G.V + 2,
        )
        self._max_matching = maxflow.result()

    def max_matching(self):
        return self._max_matching

    def is_perfect_matching(self):
        return self._max_matching * 2 == self._G.V


if __name__ == "__main__":
    filename = 'play_with_graph_algorithms/chapter15/g.txt'
    g = Graph(filename)
    bm = BipartiteMatching(g)
    print(bm.max_matching(), bm.is_perfect_matching())

    filename = 'play_with_graph_algorithms/chapter15/g2.txt'
    g = Graph(filename)
    bm = BipartiteMatching(g)
    print(bm.max_matching(), bm.is_perfect_matching())
