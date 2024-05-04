from typing import List


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


INF = int(1e18)


class Solution:
    # 2492. 两个城市间路径的最小分数
    # https://leetcode.cn/problems/minimum-score-of-a-path-between-two-cities/
    # 求无向图中1到n中所有路径的最小边权
    # 一条路径可以 多次 包含同一条道路，你也可以沿着路径多次到达城市 1 和城市 n 。
    # 测试数据保证城市 1 和城市n 之间 至少 有一条路径。
    def minScore(self, n: int, roads: List[List[int]]) -> int:
        uf = UnionFindArraySimple(n + 1)
        for u, v, _ in roads:
            uf.union(u, v)
        res = INF
        root1 = uf.find(1)
        for u, v, w in roads:
            if uf.find(u) == root1:
                res = min(res, w)
        return res

    # 3108. 带权图里旅途的最小代价
    # https://leetcode.cn/problems/minimum-cost-walk-in-weighted-graph/description/
    # 一趟旅途的 代价 是经过的边权按位与 AND 的结果。
    # 对于每一个查询，你需要找出从节点开始 si ，在节点 ti 处结束的旅途的最小代价。如果不存在这样的旅途，答案为 -1 。
    # !求每个联通分量所有节点的按位与.
    def minimumCost(self, n: int, edges: List[List[int]], query: List[List[int]]) -> List[int]:
        uf = UnionFindArraySimple(n)
        for u, v, _ in edges:
            uf.union(u, v)
        groupValue = [-1] * n
        for u, v, w in edges:
            groupValue[uf.find(u)] &= w
        res = []
        for u, v in query:
            if uf.find(u) != uf.find(v):
                res.append(-1)
            else:
                res.append(groupValue[uf.find(u)])
        return res
