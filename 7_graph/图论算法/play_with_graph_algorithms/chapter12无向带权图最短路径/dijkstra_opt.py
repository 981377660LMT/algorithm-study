from heapq import heappush
from heapq import heappop

from play_with_graph_algorithms.chapter11.weighted_graph import WeightedGraph


MAX_VALUE = 2 ** 31 - 1

class DijkstraOPT:

    def __init__(self, G, s):
        self._G = G
        self._G.validate_vertex(s)
        self._s = s

        self._dis = [MAX_VALUE - 1] * self._G.V
        self._dis[0] = 0
        # 是否已经确定了点是全局最小距离
        self._visited = [False] * self._G.V
        pq = []
        heappush(pq, (0, s))

        self._pre = [-1] * self._G.V
        self._pre[s] = s

        while pq:
            _, curr = heappop(pq)
            if self._visited[curr]:
                continue
            self._visited[curr] = True

            for w in self._G.adj(curr):
                if not self._visited[w]:
                    new_dis = self._dis[curr] + self._G.get_weight(curr, w)
                    if new_dis < self._dis[w]:
                        # dis[w]表示原点s到w的当前最短距离
                        self._dis[w] = new_dis
                        self._pre[w] = curr
                        heappush(pq, (self._dis[w], w))

    def is_connected_to(self, v):
        self._G.validate_vertex(v)
        return self._visited[v]

    def dist_to(self, v):
        self._G.validate_vertex(v)
        return self._dis[v]

    def path(self, t):
        res = []

        if not self.is_connected_to(t):
            return res

        curr = t
        while curr != self._s:
            res.append(curr)
            curr = self._pre[curr]
        res.append(self._s)
        return res[::-1]

if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter12/g.txt'
    g = WeightedGraph(filename)
    dijkstra_opt = DijkstraOPT(g, s=0)

    strings = []
    for v in range(g.V):
        strings.append(str(dijkstra_opt.dist_to(v)))
    print(' '.join(strings))

    print(' -> '.join(str(i) for i in dijkstra_opt.path(3)))
