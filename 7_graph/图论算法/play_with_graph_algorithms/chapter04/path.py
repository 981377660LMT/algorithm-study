from play_with_graph_algorithms.chapter02.adj_set import AdjSet as Graph


class Path:

    def __init__(self, G, s, t, recursive=True):
        G.validate_vertex(s)
        G.validate_vertex(t)
        self._G = G
        self._s = s
        self._t = t
        self._visited = [False] * G.V
        self._pre = [-1] * G.V

        if recursive:
            self._dfs_recursive(s, s)
        else:
            self._dfs_iteration(s, s)
        print(self._visited)

    def _dfs_recursive(self, v, parent):
        self._visited[v] = True
        self._pre[v] = parent
        if v == self._t:
            return True
        for w in self._G.adj(v):
            if not self._visited[w]:
                if self._dfs_recursive(w, v):
                    return True
        return False

    def _dfs_iteration(self, v, parent):
        stack = [v]
        self._visited[v] = True
        self._pre[v] = parent

        while stack:
            curr = stack.pop()
            # 在pop时候check比较合适
            # 因为确保了visited和pre都被赋过值了
            if curr == self._t:
                return True
            for w in self._G.adj(curr):
                if not self._visited[w]:
                    stack.append(w)
                    self._visited[w] = True
                    self._pre[w] = curr

        return False

    def is_connected(self):
        return self._visited[self._t]

    def path(self):
        res = []
        if not self.is_connected():
            return res
        curr = self._t
        while curr != self._s:
            res.append(curr)
            curr = self._pre[curr]
        res.append(self._s)
        return res[::-1]


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter04/g2.txt'
    g = Graph(filename)

    path = Path(g, 0, 6)
    print('0 -> 6: ' + str(path.path()))

    path = Path(g, 0, 5)
    print('0 -> 5: ' + str(path.path()))