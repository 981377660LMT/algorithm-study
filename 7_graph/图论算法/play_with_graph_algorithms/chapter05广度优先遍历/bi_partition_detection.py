from collections import deque
from play_with_graph_algorithms.chapter02.adj_set import AdjSet as Graph


class BiPartitionDetection:

    def __init__(self, G):
        self._G = G
        self._colors = [-1] * G.V
        self._is_bi_partition_graph = True

        for v in range(G.V):
            if self._colors[v] == -1:
                if not self._bfs(v):
                    self._is_bi_partition_graph = False
                    break

    def _bfs(self, s):
        queue = deque()
        queue.append(s)
        self._colors[s] = 0

        while queue:
            v = queue.popleft()
            for w in self._G.adj(v):
                if self._colors[w] == -1:
                    queue.append(w)
                    self._colors[w] = 1 - self._colors[v]
                # 如果下一个点w的颜色是当前处理点的颜色
                # 说明不能对其染色了
                # 即说明当前的图不是二分图
                elif self._colors[w] == self._colors[v]:
                    return False
        return True

    def is_bi_partition_graph(self):
        return self._is_bi_partition_graph


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter05/gg.txt'
    g = Graph(filename)
    bi_partition_detection = BiPartitionDetection(g)
    print('Is this a bi-partition graph? : {}'.format(
        bi_partition_detection.is_bi_partition_graph()),
    )

    filename = 'play_with_graph_algorithms/chapter05/gg2.txt'
    g = Graph(filename)
    bi_partition_detection = BiPartitionDetection(g)
    print('Is this a bi-partition graph? : {}'.format(
        bi_partition_detection.is_bi_partition_graph()),
    )

    filename = 'play_with_graph_algorithms/chapter05/gg3.txt'
    g = Graph(filename)
    bi_partition_detection = BiPartitionDetection(g)
    print('Is this a bi-partition graph? : {}'.format(
        bi_partition_detection.is_bi_partition_graph()),
    )
