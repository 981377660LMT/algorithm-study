from collections import defaultdict
from math import floor, log2
from typing import DefaultDict, List, Set


class LCAManager:
    def __init__(self, n: int, adjMap: DefaultDict[int, Set[int]], root=0) -> None:
        """查询 LCA

        `nlogn` 预处理
        `logn`查询两点的LCA

        Args:
            n (int): 树节点编号 默认 0 ~ n-1
            adjMap (DefaultDict[int, Set[int]]): 树
            root (int, optional): 根节点 默认 0
        """
        self.depth = defaultdict(lambda: -1)
        self.parent = defaultdict(lambda: -1)
        self._BITLEN = floor(log2(n)) + 1
        self._MAX = n
        self._adjMap = adjMap
        self._dfs(root, -1, 0)
        self._fa = self._makeDp(self.parent)

    def queryLCA(self, root1: int, root2: int) -> int:
        """ `logn` 查询 """
        if self.depth[root1] < self.depth[root2]:
            root1, root2 = root2, root1

        for i in range(self._BITLEN - 1, -1, -1):
            if self.depth[self._fa[root1][i]] >= self.depth[root2]:
                root1 = self._fa[root1][i]

        if root1 == root2:
            return root1

        for i in range(self._BITLEN - 1, -1, -1):
            if self._fa[root1][i] != self._fa[root2][i]:
                root1 = self._fa[root1][i]
                root2 = self._fa[root2][i]

        return self._fa[root1][0]

    def queryDist(self, root1: int, root2: int) -> int:
        """查询树节点两点间距离"""
        return self.depth[root1] + self.depth[root2] - 2 * self.depth[self.queryLCA(root1, root2)]

    def _dfs(self, cur: int, pre: int, dep: int) -> None:
        """处理高度、父节点信息"""
        self.depth[cur], self.parent[cur] = dep, pre
        for next in self._adjMap[cur]:
            if next == pre:
                continue
            self._dfs(next, cur, dep + 1)

    def _makeDp(self, parent: DefaultDict[int, int]) -> List[List[int]]:
        """nlogn预处理"""
        dp = [[0] * self._BITLEN for _ in range(self._MAX)]
        for i in range(self._MAX):
            dp[i][0] = parent[i]
        for j in range(self._BITLEN - 1):
            for i in range(self._MAX):
                if dp[i][j] == -1:
                    dp[i][j + 1] = -1
                else:
                    dp[i][j + 1] = dp[dp[i][j]][j]
        return dp

