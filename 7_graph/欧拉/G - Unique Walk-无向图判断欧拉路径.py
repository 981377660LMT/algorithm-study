"""
部分边的欧拉路径

walk:步道/路径

给定一张无向图，以及一个重要边集合。
问是否存在一条路径，其经过了每条重要边仅此一次。

如果所有边都是重要边，那题意就是求是否存在一条欧拉通路。
而如果有一些边不重要，即它可以不经过，或经过若干次，
!那该边所连接的两个点就可以相互到达无限次，我们可以将这两个点缩成一个点。
!即将原图里所有不重要的边进行缩边，最后得到一张仅由重要边组成的图，判一下是否存在欧拉通路即可。
无向连通图判断欧拉通路的条件是：仅存在0个或2个度数为奇数的点。

n<=2e5
"""

from typing import List, Tuple


def uniqueWalk(n: int, edges: List[Tuple[int, int]], eulerEdges: List[int]) -> bool:
    uf = UnionFindArray(n)
    isEulerEdge = [False] * len(edges)
    for i in eulerEdges:
        isEulerEdge[i] = True
    for i, (u, v) in enumerate(edges):
        if not isEulerEdge[i]:
            uf.union(u, v)

    deg = [0] * n
    for i, (u, v) in enumerate(edges):
        if isEulerEdge[i]:
            u, v = uf.find(u), uf.find(v)
            deg[u] += 1
            deg[v] += 1

    odd = sum(1 for x in deg if x % 2 == 1)
    return odd == 0 or odd == 2


from collections import defaultdict
from typing import DefaultDict, List


class UnionFindArray:

    __slots__ = ("n", "part", "parent", "rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        while x != self.parent[x]:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.part -= 1
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
        return list(set(self.find(key) for key in self.parent))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v))
    k = int(input())
    eulerEdges = [int(x) - 1 for x in input().split()]

    res = uniqueWalk(n, edges, eulerEdges)
    print("Yes" if res else "No")
