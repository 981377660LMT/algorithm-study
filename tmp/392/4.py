from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 n 个节点的带权无向图，节点编号为 0 到 n - 1 。

# 给你一个整数 n 和一个数组 edges ，其中 edges[i] = [ui, vi, wi] 表示节点 ui 和 vi 之间有一条权值为 wi 的无向边。

# 在图中，一趟旅途包含一系列节点和边。旅途开始和结束点都是图中的节点，且图中存在连接旅途中相邻节点的边。注意，一趟旅途可能访问同一条边或者同一个节点多次。

# 如果旅途开始于节点 u ，结束于节点 v ，我们定义这一趟旅途的 代价 是经过的边权按位与 AND 的结果。换句话说，如果经过的边对应的边权为 w0, w1, w2, ..., wk ，那么代价为w0 & w1 & w2 & ... & wk ，其中 & 表示按位与 AND 操作。

# 给你一个二维数组 query ，其中 query[i] = [si, ti] 。对于每一个查询，你需要找出从节点开始 si ，在节点 ti 处结束的旅途的最小代价。如果不存在这样的旅途，答案为 -1 。


# 返回数组 answer ，其中 answer[i] 表示对于查询 i 的 最小 旅途代价。


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


class Solution:
    def minimumCost(self, n: int, edges: List[List[int]], query: List[List[int]]) -> List[int]:
        res = [0] * len(query)
