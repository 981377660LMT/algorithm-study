# E - Reachable Set
# https://atcoder.jp/contests/abc401/tasks/abc401_e
# 给定一张无向图.
# 对k=0,1,...,n-1，能否选出一张联通子图，只包含0到k的点，如果可以，输出需要删除的最少点数.

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


class UnionFindArraySimple:
    __slots__ = ("part", "n", "_data")

    def __init__(self, n: int):
        self.part = n
        self.n = n
        self._data = [-1] * n

    def union(self, key1: int, key2: int) -> bool:
        root1, root2 = self.find(key1), self.find(key2)
        if root1 == root2:
            return False
        if self._data[root1] > self._data[root2]:
            root1, root2 = root2, root1
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        if self._data[key] < 0:
            return key
        self._data[key] = self.find(self._data[key])
        return self._data[key]

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]


if __name__ == "__main__":
    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        u -= 1
        v -= 1
        adjList[u].append(v)
        adjList[v].append(u)

    res = 0
    uf = UnionFindArraySimple(n)
    marked = [False] * n
    for i in range(n):
        if marked[i]:
            marked[i] = False
            res -= 1
        for j in adjList[i]:
            if i > j:  # !只考虑前i+1个顶点之间的连通性
                uf.union(i, j)
            else:
                if not marked[j]:
                    marked[j] = True
                    res += 1
        print(res if uf.getSize(i) == i + 1 else -1)
