from play_with_graph_algorithms.chapter11.weighted_graph import WeightedGraph


MAX_INT = 2 ** 31 - 1

class BellmanFord:

    def __init__(self, G, s):
        self._G = G
        self._G.validate_vertex(s)
        self._s = s
        self._dis = [MAX_INT] * self._G.V
        self._dis[0] = 0
        self._pre = [-1] * self._G.V
        self._has_negative_circle = False

        for _ in range(self._G.V):
            for v in range(self._G.V):
                for w in self._G.adj(v):
                    if self._dis[v] == MAX_INT:
                        continue
                    new_dis = self._dis[v] + self._G.get_weight(v, w)
                    if new_dis < self._dis[w]:
                        # 在更新距离的时候更新pre数组
                        self._pre[w] = v
                        self._dis[w] = new_dis

        # 检查是否有负权环
        for v in range(self._G.V):
            for w in self._G.adj(v):
                if self._dis[v] == MAX_INT:
                    continue
                new_dis = self._dis[v] + self._G.get_weight(v, w)
                if new_dis < self._dis[w]:
                    self._has_negative_circle = True
                    break

    def has_neg(self):
        return self._has_negative_circle

    def is_connected_to(self, v):
        self._G.validate_vertex(v)
        return self._dis[v] != MAX_INT

    def dist_to(self, v):
        self._G.validate_vertex(v)
        if self._has_negative_circle:
            raise ValueError('Exist negative circle.')
        return self._dis[v]

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
    filename = 'play_with_graph_algorithms/chapter12/g.txt'
    g = WeightedGraph(filename)
    bf = BellmanFord(g, s=0)

    if not bf.has_neg():
        strings = []
        for v in range(g.V):
            strings.append(bf.dist_to(v))
            print(bf.path(v))
        print(' '.join(str(i) for i in strings))
    else:
        print('exist negative circle')