# 928. 尽量减少恶意软件的传播 II
# https://leetcode.cn/problems/minimize-malware-spread-ii/description/
# 给定一个无向图，初始时有一些结点 initial 被病毒感染，每个病毒会感染整个联通块所有结点。
# !你可以删除 initial 中的一个结点.
# 要使得最少的结点被感染，你需要删除哪个结点？
# 如果有多个解，请返回编号最小的解。


# 并查集本身只适合用作集合的合并，并不适合用作集合的拆分.
# 当碰到`拆分` 并查集的题干，应该想到 逆向思维 地利用并查集.
# !将问题由 "删除一个感染节点，能减少多个节点受到感染" 转换成 "添加一个感染节点会增加多少个被感染节点"，
# 即添加一个感染节点，使得该节点能感染的节点最多.


from typing import List
from collections import Counter, defaultdict


class Solution:
    def minMalwareSpread(self, adjMatrix: List[List[int]], initialVirus: List[int]) -> int:
        n = len(adjMatrix)
        uf = UnionFindArraySimple(n)
        bad = [False] * n
        for v in initialVirus:
            bad[v] = True

        # 忽略所有感染节点，只考虑正常节点
        for i in range(n):
            if not bad[i]:
                for j in range(i + 1, n):
                    if not bad[j] and adjMatrix[i][j] == 1:
                        uf.union(i, j)

        belong = [uf.find(i) for i in range(n)]
        groupSize = [uf.getSize(i) for i in range(n)]

        # 每个感染源感染了哪些组
        infect = defaultdict(set)
        for u in initialVirus:
            for v in range(n):
                if not bad[v] and adjMatrix[u][v] == 1:
                    infect[u].add(belong[v])

        # 统计每个组的感染次数(被几个感染源感染)
        freq = Counter()
        for v in infect.values():
            freq.update(v)

        res = min(initialVirus)
        best = -1
        for u in initialVirus:
            validCount = 0
            for v in infect[u]:
                # v组只被u感染, 有效地感染了多少个节点
                if freq[v] == 1:
                    validCount += groupSize[v]
            if validCount > best or validCount == best and u < res:
                res, best = u, validCount
        return res


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
