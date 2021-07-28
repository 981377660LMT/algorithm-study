from play_with_graph_algorithms.chapter02.adj_set import AdjSet as Graph


class SingleSourcePath:

    def __init__(self, G, s, recursive=True):
        G.validate_vertex(s)
        self._G = G
        self._s = s
        self._visited = [False] * G.V
        self._pre = [-1] * G.V

        if recursive:
            self._dfs_recursive(s, s)
        else:
            self._dfs_iteration(s, s)

    def _dfs_recursive(self, v, parent):
        self._visited[v] = True
        self._pre[v] = parent
        for w in self._G.adj(v):
            if not self._visited[w]:
                self._dfs_recursive(w, v)

    def _dfs_iteration(self, v, parent):
        stack = [v]
        self._visited[v] = True
        self._pre[v] = parent
        while stack:
            curr = stack.pop()
            for w in self._G.adj(curr):
                if not self._visited[w]:
                    stack.append(w)
                    self._visited[w] = True
                    self._pre[w] = curr

    def is_connected_to(self, t):
        """this funtion is called during the dfs
        so if current node t is visited (self._visited[t] == True)
        this means the current t is connected to the source node
        """
        self._G.validate_vertex(t)
        return self._visited[t]

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
    filename = 'play_with_graph_algorithms/chapter04/g2.txt'
    g = Graph(filename)

    single_source_path = SingleSourcePath(g, 0)
    print('0 -> 6: ' + str(single_source_path.path(6)))

    single_source_path = SingleSourcePath(g, 0)
    print('0 -> 5: ' + str(single_source_path.path(5)))