# 924. 尽量减少恶意软件的传播
# https://leetcode.cn/problems/minimize-malware-spread/description/
# 给定一个无向图，初始时有一些结点 initial 被病毒感染，每个病毒会感染整个联通块所有结点。
# 你可以使一个 initial 中的结点免疫，免疫后该结点不会被感染。
# 要使得最少的结点被感染，你需要使哪个结点免疫？
# 如果有多个解，请返回编号最小的解。

# !求只包含一个被感染节点的最大连通块


from typing import List
from collections import defaultdict


class Solution:
    def minMalwareSpread(self, graph: List[List[int]], initial: List[int]) -> int:
        if len(initial) == 0:
            return 0
        if len(initial) == 1:
            return initial[0]

        n = len(graph)
        bad = [False] * n
        for v in initial:
            bad[v] = True
        uf = UnionFindArraySimple(n)
        for i in range(n):
            for j in range(i + 1, n):
                if graph[i][j] == 1:
                    uf.union(i, j)

        belong = [uf.find(i) for i in range(n)]
        id = defaultdict(lambda: len(id))
        belong = [id[b] for b in belong]
        groupSize, groupBadCount = [0] * len(id), [0] * len(id)
        for i, b in enumerate(belong):
            groupSize[b] += 1
            groupBadCount[b] += bad[i]

        min_, argMin = n, -1
        for v in initial:
            b = belong[v]
            if groupBadCount[b] == 1:
                cand = n - groupSize[b]
                if cand < min_ or (cand == min_ and v < argMin):
                    min_, argMin = cand, v
        return argMin if min_ < n else min(initial)


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
