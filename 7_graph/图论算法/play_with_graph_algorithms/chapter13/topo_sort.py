from collections import deque

from play_with_graph_algorithms.chapter13.graph import Graph


class TopoSort:

    def __init__(self, G):
        if not G.is_directed():
            raise ValueError('TopoSort only works in directed graph')

        self._G = G
        self._res = []
        self._has_cycle = False

        indegrees = [0] * self._G.V
        queue = deque()

        for v in range(self._G.V):
            indegrees[v] = self._G.indegree(v)
            if indegrees[v] == 0:
                queue.append(v)

        while queue:
            curr = queue.popleft()
            self._res.append(curr)
            for next_ in self._G.adj(curr):
                indegrees[next_] -= 1
                if indegrees[next_] == 0:
                    queue.append(next_)

        if len(self._res) != self._G.V:
            self._has_cycle = True
            self._res = []

    def has_cycle(self):
        return self._has_cycle
    
    def result(self):
        return self._res


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter13/ug.txt'
    g = Graph(filename, directed=True)

    topo_sort = TopoSort(g)
    print(topo_sort.result())