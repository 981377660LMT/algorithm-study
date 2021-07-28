from play_with_graph_algorithms.chapter02.adj_set import AdjSet as Graph


class HamiltonPath:

    def __init__(self, G, s):
        self._G = G
        self._s = s
        self._visited = [False] * G.V
        self._pre = [0] * G.V
        self._end = -1
        self._dfs(s, s, G.V)

    def _dfs(self, v, parent, left):
        self._visited[v] = True
        self._pre[v] = parent
        left -= 1
        if left == 0:
            self._end = v
            return True

        for w in self._G.adj(v):
            if not self._visited[w]:
                if self._dfs(w, v, left):
                    return True
        
        self._visited[v] = False

        return False

    def result(self):
        res = []
        if self._end == -1:
            return res
        
        curr = self._end
        while curr != self._s:
            res.append(curr)
            curr = self._pre[curr]
        res.append(self._s)

        return res[::-1]


class HamiltonPathV2:

    def __init__(self, G, s):
        self._G = G
        self._s = s
        self._pre = [0] * G.V
        self._end = -1
        self._dfs(0, s, s, G.V)

    def _dfs(self, visited, v, parent, left):
        visited += (1 << v)
        self._pre[v] = parent
        left -= 1
        if left == 0:
            self._end = v
            return True

        for w in self._G.adj(v):
            if (visited & (1 << w)) == 0:
                if self._dfs(visited, w, v, left):
                    return True

        return False

    def result(self):
        res = []
        if self._end == -1:
            return res
        
        curr = self._end
        while curr != self._s:
            res.append(curr)
            curr = self._pre[curr]
        res.append(self._s)

        return res[::-1]


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter09/g.txt'
    graph = Graph(filename)
    hamilton_path = HamiltonPath(graph, 0)
    print(hamilton_path.result())

    filename = 'play_with_graph_algorithms/chapter09/g2.txt'
    graph = Graph(filename)
    hamilton_path = HamiltonPath(graph, 1)
    print(hamilton_path.result())

    filename = 'play_with_graph_algorithms/chapter09/g.txt'
    graph = Graph(filename)
    hamilton_path_v2 = HamiltonPathV2(graph, 0)
    print(hamilton_path_v2.result())

    filename = 'play_with_graph_algorithms/chapter09/g2.txt'
    graph = Graph(filename)
    hamilton_path_v2 = HamiltonPathV2(graph, 1)
    print(hamilton_path_v2.result())