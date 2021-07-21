from collections import deque

from play_with_graph_algorithms.chapter13.weighted_graph import WeightedGraph


class MaxFlow:

    def __init__(self, network, s, t, empty_graph=False, directed=True, V=None):
        if not network.is_directed():
            raise ValueError('MaxFlow only works in directed graph!')

        if network.V < 2:
            raise ValueError('The network should have at least two points')

        network.validate_vertex(s)
        network.validate_vertex(t)

        if s == t:
            raise ValueError('s and t should be different')

        self._network = network
        self._s = s
        self._t = t
        if not empty_graph:
            self._rG = self._network.generate_redisual_graph()
        else:
            # 生成残量图
            temp = []
            for v in range(self._network.V):
                for w in self._network.adj(v):
                    temp.append([w, v, 0])
            for each in temp:
                self._network.add_edge(*each)
            self._rG = self._network

        self._max_flow = 0

        while True:
            aug_path = self._get_augmenting_path()
            if not aug_path:
                break
            else:
                f = 2 ** 31 - 1
                for i in range(1, len(aug_path)):
                    v = aug_path[i - 1]
                    w = aug_path[i]
                    f = min(f, self._rG.get_weight(v, w))

                self._max_flow += f
                for i in range(1, len(aug_path)):
                    v = aug_path[i - 1]
                    w = aug_path[i]

                    # 不管正向边还是反向边，更新的方式都是一样的
                    # if self._network.has_edge(v, w):
                    #     self._rG.set_weight(v, w, self._rG.get_weight(v, w) - f)
                    #     self._rG.set_weight(w, v, self._rG.get_weight(w, v) + f)
                    # else:
                    #     self._rG.set_weight(w, v, self._rG.get_weight(w, v) + f)
                    #     self._rG.set_weight(v, w, self._rG.get_weight(v, w) - f)
                    self._rG.set_weight(w, v, self._rG.get_weight(w, v) + f)
                    self._rG.set_weight(v, w, self._rG.get_weight(v, w) - f)

    def result(self):
        return self._max_flow

    def flow(self, v, w):
        if not self._network.has_edge(v, w):
            raise ValueError('No edge {}-{}'.format(v, w))
        return self._rG.get_weight(w, v)

    def _get_augmenting_path(self):
        q = deque()
        pre = [-1] * self._network.V

        q.append(self._s)
        pre[self._s] = self._s

        while q:
            curr = q.popleft()
            for next_ in self._rG.adj(curr):
                if pre[next_] == -1 and self._rG.get_weight(curr, next_) > 0:
                    pre[next_] = curr
                    q.append(next_)

        res = []
        if pre[self._t] == -1:
            return res

        curr = self._t
        while curr != self._s:
            res.append(curr)
            curr = pre[curr]
        res.append(self._s)
        return res[::-1]


if __name__ == "__main__":
    filename = 'play_with_graph_algorithms/chapter14/network.txt'
    network = WeightedGraph(filename, directed=True)

    maxflow = MaxFlow(network, 0, 3)
    print(maxflow.result())
    for v in range(network.V):
        for w in network.adj(v):
            print('{}-{} : {} / {}'.format(v, w, maxflow.flow(v, w), network.get_weight(v, w)))

    print('=' * 30)

    filename = 'play_with_graph_algorithms/chapter14/network2.txt'
    network = WeightedGraph(filename, directed=True)

    maxflow = MaxFlow(network, 0, 3)
    print(maxflow.result())
    for v in range(network.V):
        for w in network.adj(v):
            print('{}-{} : {} / {}'.format(v, w, maxflow.flow(v, w), network.get_weight(v, w)))

    filename = 'play_with_graph_algorithms/chapter14/baseball.txt'
    network = WeightedGraph(filename, directed=True)
    maxflow = MaxFlow(network, 0, 10)
    print(maxflow.result())
