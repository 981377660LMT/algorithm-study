from collections import deque
from play_with_graph_algorithms.chapter02.adj_set import AdjSet as Graph
from play_with_graph_algorithms.chapter05.single_source_path import SingleSourcePath


class AllPairsPath:

    def __init__(self, G):
        self._G = G
        self._paths = []

        for v in range(G.V):
            self._paths.append(SingleSourcePath(G, v))

    def is_connected_to(self, s, t):
        self._G.validate_vertex(s)
        return self._paths[s].is_connected_to(t)

    def path(self, s, t):
        self._G.validate_vertex(s)
        return self._paths[s].path(t)


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter04/g1.txt'
    g = Graph(filename)
    appath = AllPairsPath(g)
    print(appath.path(0, 6))