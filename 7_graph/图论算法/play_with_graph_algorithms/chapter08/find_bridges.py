from collections import namedtuple

from play_with_graph_algorithms.chapter02.adj_set import AdjSet as Graph
from play_with_graph_algorithms.chapter03.graph_dfs import GraphDFS


class Edge(namedtuple('Edge', ['v', 'w'])):

    def __str__(self):
        return '{}-{}'.format(self.v, self.w)

    def __repr__(self):
        return self.__str__()


class FindBridges:

    def __init__(self, G):
        self._G = G
        self._visited = [False] * G.V
        # _ord记录每个点的访问顺序（根据_cnt的值，而且_cnt的值每访问完一个点自加1）
        self._ord = [-1] * G.V
        # _low记录每个点在当前所有已经访问过的点中（即_ord值比当前的点的_ord低）
        # 能够到达的最低的_ord值
        self._low = [2 ** 31 - 1] * G.V
        self._cnt = 0
        self._res = []
        # 遍历所有的点，相当于遍历图中所有可能存在的联通块
        for v in range(G.V):
            if not self._visited[v]:
                self._dfs(v, v)

    def _dfs(self, v, parent):
        self._visited[v] = True
        # 设置初始值
        self._ord[v] = self._cnt
        self._low[v] = self._ord[v]
        self._cnt += 1
        for w in self._G.adj(v):
            if not self._visited[w]:
                self._dfs(w, v)
                self._low[v] = min(self._low[v], self._low[w])
                # 判断v-w是不是桥
                if self._low[w] > self._ord[v]:
                    self._res.append(Edge(v=v, w=w))
            # v w是环上的一条边
            # 肯定不是环
            # 此时可能需要更新下v的low值
            elif w != parent:
                self._low[v] = min(self._low[v], self._low[w])

    @property
    def result(self):
        return self._res


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter08/g.txt'
    g = Graph(filename)
    find_bridegs = FindBridges(g)
    print(find_bridegs.result)

    filename = 'play_with_graph_algorithms/chapter08/g2.txt'
    g = Graph(filename)
    find_bridegs = FindBridges(g)
    print(find_bridegs.result)

    filename = 'play_with_graph_algorithms/chapter08/tree.txt'
    g = Graph(filename)
    find_bridegs = FindBridges(g)
    print(find_bridegs.result)