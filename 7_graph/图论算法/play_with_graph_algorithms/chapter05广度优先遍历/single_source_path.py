from collections import deque
from play_with_graph_algorithms.chapter02.adj_set import AdjSet as Graph

class SingleSourcePath:

    def __init__(self, G, s):
        self._G = G
        self._visited = [False] * G.V
        self._s = s
        self._pre = [-1] * G.V

        self._bfs(s)

    def _bfs(self, s):
        queue = deque()
        queue.append(s)
        self._visited[s] = True
        self._pre[s] = s

        while queue:
            v = queue.popleft()
            for w in self._G.adj(v):
                if not self._visited[w]:
                    queue.append(w)
                    self._visited[w] = True
                    self._pre[w] = v

    def is_connected_to(self, t):
        self._G.validate_vertex(t)
        # if t is visited for current self._s
        # that implies t is in the path of self._s
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
    filename = 'play_with_graph_algorithms/chapter04/g1.txt'
    g = Graph(filename)
    sspath = SingleSourcePath(g, 0)
    print('0 -> 6 : {}'.format(sspath.path(6)))