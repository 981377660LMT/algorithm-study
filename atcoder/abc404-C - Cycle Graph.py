# https://atcoder.jp/contests/abc404/tasks/abc404_c
# 判断一张图是否是环图
# 1. 连通分量为1
# 2. 每个点的度数为2

from collections import defaultdict
import sys
from typing import Callable, DefaultDict, List

input = lambda: sys.stdin.readline().rstrip("\r\n")


class UnionFindArray:
    """元素是0-n-1的并查集写法,不支持动态添加

    初始化的连通分量个数 为 n
    """

    __slots__ = ("n", "part", "_parent", "_rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self._parent = list(range(n))
        self._rank = [1] * n

    def find(self, x: int) -> int:
        while self._parent[x] != x:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        """按秩合并."""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self._rank[rootX] > self._rank[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        return True

    def unionTo(self, child: int, parent: int) -> bool:
        """定向合并.将child的父节点设置为parent."""
        rootX = self.find(child)
        rootY = self.find(parent)
        if rootX == rootY:
            return False
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        return True

    def unionWithCallback(self, x: int, y: int, f: Callable[[int, int], None]) -> bool:
        """
        f: 合并后的回调函数, 入参为 (big, small)
        """
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self._rank[rootX] > self._rank[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        f(rootY, rootX)
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(key) for key in self._parent))

    def getSize(self, x: int) -> int:
        return self._rank[self.find(x)]

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


if __name__ == "__main__":
    N, M = map(int, input().split())
    uf = UnionFindArray(N)
    deg = [0] * N
    for _ in range(M):
        u, v = map(int, input().split())
        u -= 1
        v -= 1
        uf.union(u, v)
        deg[u] += 1
        deg[v] += 1
    if uf.part != 1:
        print("No")
        exit(0)
    for d in deg:
        if d != 2:
            print("No")
            exit(0)
    print("Yes")
