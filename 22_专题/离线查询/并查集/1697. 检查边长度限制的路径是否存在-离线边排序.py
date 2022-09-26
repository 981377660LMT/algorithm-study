"""
对于每个查询 queries[j] ，判断是否存在从 pj 到 qj 的路径，
且这条路径上的每一条边都 严格小于 limitj 。

按照边权从小到大排序,然后并查集合并
"""

from typing import List

from collections import defaultdict
from typing import DefaultDict, List


class UF:
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


class Solution:
    def distanceLimitedPathsExist(
        self, n: int, edgeList: List[List[int]], queries: List[List[int]]
    ) -> List[bool]:
        edgeList.sort(key=lambda x: x[2])
        queriesWithIndex = [(i, *q) for i, q in enumerate(queries)]
        queriesWithIndex.sort(key=lambda x: x[3])

        ei = 0
        res = [False] * len(queries)
        uf = UF(n)
        for qi, qu, qv, qw in queriesWithIndex:
            while ei < len(edgeList) and edgeList[ei][2] < qw:
                uf.union(edgeList[ei][0], edgeList[ei][1])
                ei += 1
            res[qi] = uf.isConnected(qu, qv)
        return res


print(
    Solution().distanceLimitedPathsExist(
        3, [[0, 1, 2], [1, 2, 4], [2, 0, 8], [1, 0, 16]], [[0, 1, 2], [0, 2, 5]]
    )
)
