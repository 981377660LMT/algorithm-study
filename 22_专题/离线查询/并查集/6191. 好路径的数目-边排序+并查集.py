"""
一条 好路径 需要满足以下条件：

开始节点和结束节点的值 相同 。
开始节点和结束节点中间的所有节点值都小于等于开始节点的值
（也就是说开始节点的值应该是路径上所有节点的最大值）。
请你返回不同好路径的数目。

注意，一条路径和它反向的路径算作 同一 路径。
比方说， 0 -> 1 与 1 -> 0 视为同一条路径。单个节点也视为一条合法路径。

n<=3e4
"""

# !按顺序并入，然后再查 —— 并查集(离线查询排序的技巧:边排序+反向并查集)
# !把所有边按两个端点的最大值排序分组，然后按这个最大值从小到大依次合并，
# 每合并完一组查看当前组对应的点的连通情况并统计答案，最后别忘了加上只有单个点的路径

from collections import defaultdict
from typing import List


class UnionFind:
    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        """union后x所在的root的parent指向y所在的root"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False

        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootY]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


class Solution:
    def numberOfGoodPaths(self, vals: List[int], edges: List[List[int]]) -> int:
        n = len(vals)
        nodeGroup = defaultdict(list)
        edgeGroup = defaultdict(list)  # 按端点最大值分组

        for i, num in enumerate(vals):
            nodeGroup[num].append(i)
        for u, v in edges:
            max_ = max(vals[u], vals[v])
            edgeGroup[max_].append((u, v))

        res = 0
        uf = UnionFind(n)
        keys = sorted(edgeGroup)  # 按顺序合并边
        for curMax in keys:
            for u, v in edgeGroup[curMax]:
                uf.union(u, v)
            groupCounter = defaultdict(int)  # !统计当前最大值在哪些连通分量中 同一个连通分量可以连接路径
            for node in nodeGroup[curMax]:
                groupCounter[uf.find(node)] += 1
            res += sum([v * (v - 1) // 2 for v in groupCounter.values()])  # !comb(v, 2)

        return res + n  # 单个点的路径
