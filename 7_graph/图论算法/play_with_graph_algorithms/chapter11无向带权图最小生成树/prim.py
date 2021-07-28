from heapq import heappush
from heapq import heappop

from play_with_graph_algorithms.chapter11.weighted_graph import WeightedGraph
from play_with_graph_algorithms.chapter11.weighted_edge import WeightedEdge
from play_with_graph_algorithms.chapter04.cc import CC


class Prim:

    def __init__(self, G):

        self._G = G
        self._mst = []

        cc = CC(G)
        if (cc.ccount > 1):
            return

        visited = [False] * self._G.V
        visited[0] = True
        pq = []

        for w in self._G.adj(0):
            heappush(pq, [self._G.get_weight(0, w), 0, w])

        while pq:
            weight, v, w = heappop(pq)
            if visited[v] and visited[w]:
                continue
            self._mst.append(WeightedEdge(v, w, weight))
            newv = w if visited[v] else v
            visited[newv] = True
            for w in self._G.adj(newv):
                if not visited[w]:
                    heappush(pq, [self._G.get_weight(newv, w), newv, w])

    def result(self):
        return self._mst


if __name__ == "__main__":
    filename = 'play_with_graph_algorithms/chapter11/g.txt'
    g = WeightedGraph(filename)
    prim = Prim(g)
    print(prim.result())