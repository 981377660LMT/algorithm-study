from play_with_graph_algorithms.chapter11.weighted_graph import WeightedGraph


MAX_VALUE = 2 ** 31 - 1

class Dijkstra:

    def __init__(self, G, s):
        self._G = G
        self._G.validate_vertex(s)
        self._s = s

        self._dis = [MAX_VALUE - 1] * self._G.V
        self._dis[0] = 0
        # 是否已经确定了点是全局最小距离
        self._visited = [False] * self._G.V

        while True:
            currdis = MAX_VALUE - 1
            curr = -1
            # 第一段逻辑:找最小值
            # 第一次循环找到的点肯定是s
            for v in range(self._G.V):
                if not self._visited[v] and self._dis[v] < currdis:
                    currdis = self._dis[v]
                    curr = v
            if curr == -1:
                break
            # 第二段逻辑：确定curr点被访问过
            self._visited[curr] = True

            # 第三段逻辑：利用curr点来更新其相邻节点w与原点s的距离
            for w in self._G.adj(curr):
                if not self._visited[w]:
                    if self._dis[curr] + self._G.get_weight(curr, w) < self._dis[w]:
                        # dis[w]表示原点s到w的当前最短距离
                        self._dis[w] = self._dis[curr] + self._G.get_weight(curr, w)

    def is_connected_to(self, v):
        self._G.validate_vertex(v)
        return self._visited[v]

    def dist_to(self, v):
        self._G.validate_vertex(v)
        return self._dis[v]


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter12/g.txt'
    g = WeightedGraph(filename)
    dijkstra = Dijkstra(g, s=0)

    strings = []
    for v in range(g.V):
        strings.append(str(dijkstra.dist_to(v)))

    print(' '.join(strings))