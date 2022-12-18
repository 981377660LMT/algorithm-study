# 给定一个n个顶点m条边的无向无权图(注意不一定连通)
# 求满足以下条件的结点组(u,v) 个数:
# 1. u,v间不存在边
# 2. u,v连接后,图是一个二分图

# 解:
# !1.如果不是二分图,那么就不存在这样的结点组
# !2.考虑不合法的情况,对每个连通块进行计算,不合法的对数为同色的点对数

from collections import deque
from typing import List, Tuple


def makeBipartite2(n: int, edges: List[Tuple[int, int]]) -> int:
    adjList = [[] for _ in range(n)]
    uf = UnionFindArray(n)
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)
        uf.union(u, v)
    colors, ok = isBipartite1(n, adjList)
    if not ok:
        return 0

    res = n * (n - 1) // 2 - len(edges)
    for group in uf.getGroups().values():
        zero, one = sum(colors[i] == 0 for i in group), sum(colors[i] == 1 for i in group)
        res -= zero * (zero - 1) // 2 + one * (one - 1) // 2
    return res


def isBipartite1(n: int, adjList: List[List[int]]) -> Tuple[List[int], bool]:
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
    print(makeBipartite2(n, edges))
