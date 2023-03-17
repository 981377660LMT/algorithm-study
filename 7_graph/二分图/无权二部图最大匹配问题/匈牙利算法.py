# L,R<=1e5 M<=2e5
# 0<=ai<L 0<=bi<R
# 可能存在重边
# 匈牙利算法

from collections import deque
from typing import List, Optional, Tuple


class Hungarian:
    """
    軽量化Dinic法
    ref : https://snuke.hatenablog.com/entry/2019/05/07/013609
    """

    __slots__ = ("_row", "_col", "_to")

    def __init__(self, graph: Optional[List[List[int]]] = None):
        self._row = 0
        self._col = 0
        self._to = [[]]
        if graph is not None:
            colors, ok = isBipartite(len(graph), graph)
            if not ok:
                raise ValueError("graph is not bipartite")
            for u, vs in enumerate(graph):
                if colors[u] == 0:
                    for v in vs:
                        if colors[v] == 1:
                            self.addEdge(u, v)

    def addEdge(self, u: int, v: int) -> None:
        """男孩u和女孩v连边"""
        if self._col <= v:
            self._col = v + 1
        if self._row <= u:
            self._row = u + 1
            while len(self._to) <= u:
                self._to.append([])
        self._to[u].append(v)

    def work(self) -> List[Tuple[int, int]]:
        """返回最大匹配"""
        n, m, to = self._row, self._col, self._to
        pre = [-1] * n
        root = [-1] * n
        p = [-1] * n
        q = [-1] * m
        upd = True
        while upd:
            upd = False
            s = []
            s_front = 0
            for i in range(n):
                if p[i] == -1:
                    root[i] = i
                    s.append(i)
            while s_front < len(s):
                v = s[s_front]
                s_front += 1
                if p[root[v]] != -1:
                    continue
                for u in to[v]:
                    if q[u] == -1:
                        while u != -1:
                            q[u] = v
                            p[v], u = u, p[v]
                            v = pre[v]
                        upd = True
                        break
                    u = q[u]
                    if pre[u] != -1:
                        continue
                    pre[u] = v
                    root[u] = root[v]
                    s.append(u)
            if upd:
                for i in range(n):
                    pre[i] = -1
                    root[i] = -1
        return [(v, p[v]) for v in range(n) if p[v] != -1]


def isBipartite(n: int, adjList: List[List[int]]) -> Tuple[List[int], bool]:
    """二分图检测 bfs染色"""

    def bfs(start: int) -> bool:
        colors[start] = 0
        queue = deque([start])
        while queue:
            cur = queue.popleft()
            for next in adjList[cur]:
                if colors[next] == -1:
                    colors[next] = colors[cur] ^ 1
                    queue.append(next)
                elif colors[next] == colors[cur]:
                    return False
        return True

    colors = [-1] * n
    for i in range(n):
        if colors[i] == -1 and not bfs(i):
            return [], False
    return colors, True


if __name__ == "__main__":
    L, R, M = map(int, input().split())
    hungarian = Hungarian()
    for _ in range(M):
        u, v = map(int, input().split())
        hungarian.addEdge(u, v)
    matching = hungarian.work()
    print(len(matching))
    for u, v in matching:
        print(u, v)
