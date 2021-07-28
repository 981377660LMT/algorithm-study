from play_with_graph_algorithms.chapter02.adj_set import AdjSet as Graph


class CC:

    def __init__(self, G, recursive=True):
        self._G = G
        self._visited = [-1] * G.V
        self._ccount = 0

        if self._G.is_directed:
            raise ValueError('CC only works in undirected graph')

        # 遍历所有的点，相当于遍历图中所有可能存在的联通块
        for v in range(G.V):
            if self._visited[v] == -1:
                if recursive:
                    self._dfs_recursive(v, self._ccount)
                else:
                    self._dfs_iteration(v, self._ccount)
                self._ccount += 1

    def _dfs_recursive(self, v, ccid):
        self._visited[v] = ccid
        for w in self._G.adj(v):
            if self._visited[w] == -1:
                self._dfs_recursive(w, ccid)

    def _dfs_iteration(self, v, ccid):
        """For preorder that's straight-forward by using one stack,
        but for postorder we need a augmented stack2
        """
        stack = [v]
        # within one func call, all visited nodes should be the same ccid
        self._visited[v] = ccid
        while stack:
            curr = stack.pop()
            for w in self._G.adj(curr):
                if self._visited[w] == -1:
                    stack.append(w)
                    # the same ccid
                    self._visited[w] = ccid

    def is_connected(self, v, w):
        self._G.validate_vertex(v)
        self._G.validate_vertex(w)
        return self._visited[v] == self._visited[w]

    @property
    def ccount(self):
        return self._ccount

    @property
    def groups(self):
        res = [[] for _ in range(self._ccount)]
        for v in range(self._G.V):
            res[self._visited[v]].append(v)
        return res


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter04/g1.txt'
    g = Graph(filename)
    cc = CC(g)
    print(cc.ccount)
    print(cc.groups)

    filename = 'play_with_graph_algorithms/chapter04/g2.txt'
    g = Graph(filename)
    cc = CC(g)
    print(cc.ccount)
    print(cc.groups)

    filename = 'play_with_graph_algorithms/chapter04/g1.txt'
    g = Graph(filename)
    cc = CC(g, recursive=False)
    print(cc.ccount)
    print(cc.groups)

    filename = 'play_with_graph_algorithms/chapter04/g2.txt'
    g = Graph(filename)
    cc = CC(g, recursive=False)
    print(cc.ccount)
    print(cc.groups)

    print(cc.is_connected(0, 6))
    print(cc.is_connected(0, 5))