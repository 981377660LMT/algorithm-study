from play_with_graph_algorithms.chapter02.adj_set import AdjSet as Graph


class CycleDetection:

    def __init__(self, G):
        self._G = G
        self._visited = [False] * G.V
        self._has_cycle = False

        # 遍历所有的点，相当于遍历图中所有可能存在的联通块
        for v in range(G.V):
            if not self._visited[v]:
                self._dfs_recursive(v, v)

    def _dfs_recursive(self, v, parent):
        self._visited[v] = True
        for w in self._G.adj(v):
            if not self._visited[w]:
                if self._dfs_recursive(w, v):
                    return True
            # 说明此时找到一个环
            elif w != parent:
                self._has_cycle = True
                return True
        return False

    def has_cycle(self):
        return self._has_cycle


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter04/g2.txt'
    g = Graph(filename)
    cycle_detection = CycleDetection(g)
    print(cycle_detection.has_cycle())

    filename = 'play_with_graph_algorithms/chapter04/g3.txt'
    g = Graph(filename)
    cycle_detection = CycleDetection(g)
    print(cycle_detection.has_cycle())