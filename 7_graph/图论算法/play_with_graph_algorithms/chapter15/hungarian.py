from collections import deque

from play_with_graph_algorithms.chapter13.graph import Graph
from play_with_graph_algorithms.chapter04.bi_partite_detection import BiPartiteDetection


class Hungarian:

    def __init__(self, G, use_dfs=False):
        bd = BiPartiteDetection(G)
        if bd.is_bi_partite() is False:
            raise ValueError('Hungarian only works for bipartite graph')

        self._G = G
        self._maxmatching = 0
        self._matching = [-1] * self._G.V

        colors = bd.colors()
        for v in range(self._G.V):
            # 找到了左侧一个未匹配的点
            if colors[v] == 0 and self._matching[v] == -1:
                # bfs/dfs定义是从v出发，寻找是否(bool)有增广路径
                if use_dfs:
                    self._visited = [False] * self._G.V
                    find_aug_path = self._dfs(v)
                else:
                    find_aug_path = self._bfs(v)
                if find_aug_path:
                    self._maxmatching += 1

    def _dfs(self, v):
        self._visited[v] = True
        for u in self._G.adj(v):
            # 此时u是右半部分的点
            if self._visited[u] is False:
                self._visited[u] = True
                # self._matching[u]一定是左半部分的点
                # 返回值表示能不能找到增广路径
                if self._matching[u] == -1 or self._dfs(self._matching[u]):
                    self._matching[v] = u
                    self._matching[u] = v
                    return True
        return False

    def _bfs(self, v):
        q = deque()
        pre = [-1] * self._G.V

        q.append(v)
        pre[v] = v
        while q:
            # curr一定需要是二分图左侧的边
            cur = q.popleft()
            # 由于G是一个二分图，而且curr是在左侧
            # 所以next_一定是右侧的点
            for next_ in self._G.adj(cur):
                if pre[next_] == -1:
                    # self._matching[next_]是左侧的点
                    if self._matching[next_] != -1:
                        pre[next_] = cur
                        pre[self._matching[next_]] = next_
                        q.append(self._matching[next_])
                    else:
                        pre[next_] = cur
                        aug_path = self._get_aug_path(pre, v, next_)
                        for i in range(0, len(aug_path), 2):
                            self._matching[aug_path[i]] = aug_path[i + 1]
                            self._matching[aug_path[i + 1]] = aug_path[i]
                            i += 2
                        return True
        return False


    def _get_aug_path(self, pre, start, end):
        res = []
        cur = end
        while cur != start:
            res.append(cur)
            cur = pre[cur]
        res.append(start)
        return res

    def maxmatching(self):
        return self._maxmatching


if __name__ == "__main__":
    print('bfs:')
    filename = 'play_with_graph_algorithms/chapter15/g.txt'
    g = Graph(filename)
    hungarian = Hungarian(g)
    print(hungarian.maxmatching())

    filename = 'play_with_graph_algorithms/chapter15/g2.txt'
    g = Graph(filename)
    hungarian = Hungarian(g)
    print(hungarian.maxmatching())

    print('dfs:')
    filename = 'play_with_graph_algorithms/chapter15/g.txt'
    g = Graph(filename)
    hungarian = Hungarian(g, use_dfs=True)
    print(hungarian.maxmatching())

    filename = 'play_with_graph_algorithms/chapter15/g2.txt'
    g = Graph(filename)
    hungarian = Hungarian(g, use_dfs=True)
    print(hungarian.maxmatching())