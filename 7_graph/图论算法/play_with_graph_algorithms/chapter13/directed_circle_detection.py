from play_with_graph_algorithms.chapter13.graph import Graph


class DirectedCycleDetection:

    def __init__(self, G):
        if not G.is_directed():
            raise ValueError('CircleDetection only works in undirected graph')
        self._G = G
        self._visited = [False] * G.V
        self._on_path = [False] * G.V
        self._has_cycle = False

        # 遍历所有的点，相当于遍历图中所有可能存在的联通块
        for v in range(G.V):
            if not self._visited[v]:
                if self._dfs(v, v):
                    self._has_cycle = True
                    break

    def _dfs(self, v, parent):
        self._visited[v] = True
        self._on_path[v] = True
        for w in self._G.adj(v):
            if not self._visited[w]:
                if self._dfs(w, v):
                    return True
            # 说明此时找到一个环
            # 跟无向图不同，这里没有w != parent的判断
            # 因为有向图中，认为0-1和1-0是两条不同的路径(如果存在的话)，是合法的
            elif self._on_path[w]:
                return True
        self._on_path[v] = False
        return False

    def has_cycle(self):
        return self._has_cycle


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter13/ug.txt'
    g = Graph(filename, directed=True)
    directed_cycle_detection = DirectedCycleDetection(g)
    print(directed_cycle_detection.has_cycle())
