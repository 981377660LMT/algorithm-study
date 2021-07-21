from play_with_graph_algorithms.chapter13.graph import Graph
from play_with_graph_algorithms.chapter13.graph_dfs import GraphDFS


class SCC:

    def __init__(self, G, recursive=True):
        self._G = G
        self._visited = [-1] * G.V
        self._sccount = 0

        if not self._G.is_directed:
            raise ValueError('CC only works in directed graph')

        dfs = GraphDFS(self._G.reverse_graph())
        order = dfs.post_order[::-1]

        # 遍历所有的点，相当于遍历图中所有可能存在的强联通块
        for v in order:
            if self._visited[v] == -1:
                if recursive:
                    self._dfs_recursive(v, self._sccount)
                else:
                    self._dfs_iteration(v, self._sccount)
                self._sccount += 1

    def _dfs_recursive(self, v, sccid):
        self._visited[v] = sccid
        for w in self._G.adj(v):
            if self._visited[w] == -1:
                self._dfs_recursive(w, sccid)

    def _dfs_iteration(self, v, sccid):
        """For preorder that's straight-forward by using one stack,
        but for postorder we need a augmented stack2
        """
        stack = [v]
        # within one func call, all visited nodes should be the same sccid
        self._visited[v] = sccid
        while stack:
            curr = stack.pop()
            for w in self._G.adj(curr):
                if self._visited[w] == -1:
                    stack.append(w)
                    # the same sccid
                    self._visited[w] = sccid

    def is_strongly_connected(self, v, w):
        self._G.validate_vertex(v)
        self._G.validate_vertex(w)
        return self._visited[v] == self._visited[w]

    @property
    def sccount(self):
        return self._sccount

    @property
    def components(self):
        res = [[] for _ in range(self._sccount)]
        for v in range(self._G.V):
            res[self._visited[v]].append(v)
        return res


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter13/ug4.txt'
    g = Graph(filename, directed=True)
    scc = SCC(g)

    print(scc.sccount)
    comp = scc.components
    for sccid in range(len(comp)):
        temp = ['{} : '.format(sccid)]
        for w in comp[sccid]:
            temp.append('{} '.format(w))
        print(''.join(temp))
