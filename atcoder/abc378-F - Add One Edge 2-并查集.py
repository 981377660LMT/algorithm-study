# F - Add One Edge 2
# https://atcoder.jp/contests/abc378/editorial/11304
# !给定一棵树，求加一条边的方案数，使得没有重边，且环上的所有点的度数为3。
#
# 等价于：
# !求点对(u,v)数量，满足：
# !u和v的度数为2, u到v路径上所有点(除开u、v)的度数为3。
#
# 将所有度数为3的点合并，然后统计相邻度数为2的点数量。


from collections import defaultdict


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
    N = int(input())
    U, V = [0] * (N - 1), [0] * (N - 1)
    deg = [0] * N
    for i in range(N - 1):
        U[i], V[i] = map(int, input().split())
        U[i] -= 1
        V[i] -= 1
        deg[U[i]] += 1
        deg[V[i]] += 1

    uf = UnionFindArraySimple(N)
    c2 = [0] * N  # 邻居度数为2的点数量
    for a, b in zip(U, V):
        if deg[a] == deg[b] == 3:
            uf.union(a, b)
        elif deg[a] == 3 and deg[b] == 2:
            c2[a] += 1
        elif deg[a] == 2 and deg[b] == 3:
            c2[b] += 1

    groups = defaultdict(list)
    for i in range(N):
        groups[uf.find(i)].append(i)

    res = 0
    for g in groups.values():
        c = 0
        for v in g:
            c += c2[v]
        res += c * (c - 1) // 2
    print(res)
