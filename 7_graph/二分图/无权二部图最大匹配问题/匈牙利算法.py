# L,R<=1e5 M<=2e5
# 0<=ai<L 0<=bi<R
# 不存在重边

from typing import List, Tuple


class Hungarian:
    """
    軽量化Dinic法
    ref : https://snuke.hatenablog.com/entry/2019/05/07/013609
    """

    def __init__(self, row: int, col: int):
        """
        Args:
            row (int): 男孩的个数
            col (int): 女孩的个数
        """
        self._row = row
        self._col = col
        self._to = [[] for _ in range(row)]

    def addEdge(self, u: int, v: int) -> None:
        """男孩u和女孩v连边"""
        assert 0 <= u < self._row
        assert 0 <= v < self._col
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


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")

if __name__ == "__main__":
    L, R, M = map(int, input().split())  # L个左边的点，R个右边的点，M条边
    hungarian = Hungarian(L, R)
    for _ in range(M):
        u, v = map(int, input().split())
        hungarian.addEdge(u, v)
    matching = hungarian.work()
    print(len(matching))
    for u, v in matching:
        print(u, v)
