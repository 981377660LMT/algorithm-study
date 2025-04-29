# 685. 冗余连接 II
# https://leetcode.cn/problems/redundant-connection-ii/description/
# !有向基环树删边成树
#
# 我们不单要判断是否有环，还要看每个节点的入度是否不超过1（根节点的入度为0）。
# 如果所有节点的入度都为1，那么直接用并查集找出最后成环的边即可。
# 如果发现某节点的入度为2，那么最后的答案一定在这个节点的两条边中。
# !删去一条边，如果剩下的边不成环，那么这条边就是答案；否则另外一条边就是答案。

from typing import List


def findLoop(edges: List[List[int]], removed: List[int] = [-1, -1]) -> List[int]:
    n = len(edges)
    uf = UnionFindArray(n)
    for u, v in edges:
        if [u, v] == removed:
            continue
        if uf.isConnected(u, v):
            return [u, v]
        uf.union(u, v)
    return []


class Solution:
    def findRedundantDirectedConnection(self, edges: List[List[int]]) -> List[int]:
        edges = [[u - 1, v - 1] for u, v in edges]
        n = len(edges)
        rAdjList = [[] for _ in range(n)]
        indeg = [0] * n
        for u, v in edges:
            rAdjList[v].append(u)
            indeg[v] += 1

        if all(d == 1 for d in indeg):
            res = findLoop(edges)
            return [res[0] + 1, res[1] + 1]

        for v in range(n):
            if indeg[v] == 2:
                for u in rAdjList[v][::-1]:  # 从后往前删，因为后面的边可能是答案
                    if not findLoop(edges, removed=[u, v]):
                        return [u + 1, v + 1]

        raise Exception("No answer")


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
