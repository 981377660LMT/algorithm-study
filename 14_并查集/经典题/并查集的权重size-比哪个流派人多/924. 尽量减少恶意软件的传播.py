# 假设 M(initial) 是在恶意软件停止传播之后，整个网络中感染恶意软件的最终节点数。
# 如果从 initial 中移除某一节点能够最小化 M(initial)， 返回该节点。如果有多个节点满足条件，就返回索引最小的节点。
# 请注意，如果某个节点已从受感染节点的列表 initial 中删除，它以后仍有可能因恶意软件传播而受到感染。


from itertools import combinations
from typing import List
from collections import Counter


class Solution:
    def minMalwareSpread(self, graph: List[List[int]], initial: List[int]) -> int:
        n = len(graph)
        uf = UnionFindArray(n)
        for i, j in combinations(range(n), 2):
            if graph[i][j]:
                uf.union(i, j)

        virusGroup = Counter(uf.find(u) for u in initial)
        res, maxSave = min(initial), -1
        for node in initial:
            root = uf.find(node)
            if virusGroup[root] == 1:  # 初始感染者里这个流派只有一个人感染了,哪个流派人多就拯救哪个流派,否则会感染更多人
                if uf.rank[root] > maxSave or uf.rank[root] == maxSave and node < res:
                    res, maxSave = node, uf.rank[root]
        return res


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


print(Solution().minMalwareSpread(graph=[[1, 1, 0], [1, 1, 0], [0, 0, 1]], initial=[0, 1]))
