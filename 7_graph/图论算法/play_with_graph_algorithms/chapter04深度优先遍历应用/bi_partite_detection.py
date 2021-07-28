from play_with_graph_algorithms.chapter02.adj_set import AdjSet as Graph


class BiPartiteDetection:

    def __init__(self, G):
        self._G = G
        self._visited = [False] * G.V
        # -1 means not colored yet, should be 0 or 1
        self._colors = [-1] * G.V
        self._is_bi_partite = True
        for v in range(G.V):
            if not self._visited[v]:
                # if v is not visited
                # means we are entering a new tree!!
                # so it doesn't matter to color it as 0 or 1
                if self._dfs_recursive(v, 0) is False:
                    self._is_bi_partite = False
                    break

    # v should be colored using "color" value
    # the return value of _dfs_recursive represents whether
    # current tree is a bi-partite graph or not
    def _dfs_recursive(self, v, color):
        self._visited[v] = True
        self._colors[v] = color
        for w in self._G.adj(v):
            if not self._visited[w]:
                if not self._dfs_recursive(w, 1 - color):
                    return False
            elif self._colors[w] == self._colors[v]:
                return False
        return True

    def is_bi_partite(self):
        return self._is_bi_partite

    def colors(self):
        return self._colors


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter04/g3.txt'
    g = Graph(filename)
    bi_partite_detection = BiPartiteDetection(g)
    print(bi_partite_detection.is_bi_partite())

    filename = 'play_with_graph_algorithms/chapter04/g4.txt'
    g = Graph(filename)
    bi_partite_detection = BiPartiteDetection(g)
    print(bi_partite_detection.is_bi_partite())
