from copy import copy

from play_with_graph_algorithms.chapter02.adj_set import AdjSet as Graph
from play_with_graph_algorithms.chapter04.cc import CC


class EulerLoop:
    def __init__(self, G):
        self._G = G

    def has_euler_loop(self):
        cc = CC(self._G)
        if cc.ccount > 1:
            return False

        for v in range(self._G.V()):
            if self._G.degree(v) % 2 != 0:
                return False

        return True

    def result(self):
        res = []
        if not self.has_euler_loop:
            return res
        g = copy(self._G)

        stack = []
        currv = 0
        stack.append(currv)

        while stack:
            if g.degree(currv) != 0:
                stack.append(currv)
                # 模拟一个iterator
                w = next(self._iter_next_adj(g.adj(currv)))
                g.remove_edge(currv, w)
                currv = w
            else:
                # 此时说明找到了一个环
                res.append(currv)
                currv = stack.pop()

        return res

    def _iter_next_adj(self, adj):
        yield from sorted(adj)


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter10/g.txt'
    g = Graph(filename)
    eluer_loop = EulerLoop(g)
    print(eluer_loop.result())

    filename = 'play_with_graph_algorithms/chapter10/g2.txt'
    g = Graph(filename)
    eluer_loop = EulerLoop(g)
    print(eluer_loop.result())