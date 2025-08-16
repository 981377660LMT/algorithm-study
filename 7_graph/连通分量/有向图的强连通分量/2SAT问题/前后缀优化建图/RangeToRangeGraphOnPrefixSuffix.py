from typing import Callable


class RangeToRangeGraphOnPrefixSuffix1:
    """前后缀优化建图1 - 点向区间连边."""

    __slot__ = ("_n", "_size")

    def __init__(self, n: int):
        self._n = n
        self._size = 3 * n  # [0,n):原始点，[n,2n):后缀点，[2n,3n):前缀点

    def size(self) -> int:
        """新图的结点数, 前n个节点为原图的节点."""
        return self._size

    def init(self, f: Callable[[int, int], None]):
        """初始化辅助边."""
        n1, n2 = self._n, self._n * 2
        for i in range(self._n):
            f(i + n1, i)
            f(i + n2, i)
            if i > 0:
                f(i + n1 - 1, i + n1)
                f(i + n2, i + n2 - 1)

    def add(self, fromNode: int, toNode: int, f: Callable[[int, int], None]):
        """添加有向边 from -> to."""
        f(fromNode, toNode)

    def addPointToPrefix(self, fromNode: int, prefixEnd: int, f: Callable[[int, int], None]):
        """添加从单点到前缀区间的边：from -> [0, prefixEnd)."""
        if prefixEnd <= 0:
            return
        f(fromNode, self._n * 2 + prefixEnd - 1)

    def addPointToSuffix(self, fromNode: int, suffixStart: int, f: Callable[[int, int], None]):
        """添加从单点到后缀区间的边：from -> [suffixStart, n)."""
        if suffixStart >= self._n:
            return
        f(fromNode, self._n + suffixStart)


class RangeToRangeGraphOnPrefixSuffix2:
    """前后缀优化建图2 - 区间向点连边."""

    __slot__ = ("_n", "_size")

    def __init__(self, n: int):
        self._n = n
        self._size = 3 * n  # [0,n):原始点，[n,2n):前缀点，[2n,3n):后缀点

    def size(self) -> int:
        """新图的结点数, 前n个节点为原图的节点."""
        return self._size

    def init(self, f: Callable[[int, int], None]):
        """初始化辅助边."""
        n1, n2 = self._n, self._n * 2
        for i in range(self._n):
            f(i, i + n1)
            f(i, i + n2)
            if i > 0:
                f(i + n1 - 1, i + n1)
                f(i + n2, i + n2 - 1)

    def add(self, fromNode: int, toNode: int, f: Callable[[int, int], None]):
        """添加有向边 from -> to."""
        f(fromNode, toNode)

    def addPrefixToPoint(self, prefixEnd: int, toNode: int, f: Callable[[int, int], None]):
        """添加从前缀区间到单点的边：[0, prefixEnd) -> to."""
        if prefixEnd <= 0:
            return
        f(self._n + prefixEnd - 1, toNode)

    def addSuffixToPoint(self, suffixStart: int, toNode: int, f: Callable[[int, int], None]):
        """添加从后缀区间到单点的边：[suffixStart, n) -> to."""
        if suffixStart >= self._n:
            return
        f(self._n * 2 + suffixStart, toNode)


class RangeToRangeGraphOnPrefixSuffix3:
    """前后缀优化建图3 - 区间向区间连边."""

    __slots__ = ("_n", "maxSize", "allocPtr")

    def __init__(self, n: int, rangeToRangeOpCount: int):
        """
        新建一个区间图，n 为原图的节点数，rangeToRangeOpCount 为区间到区间的最大操作次数.

        节点分布:
        [0,n): 原始点
        [n,2n): 前缀1
        [2n,3n): 后缀1
        [3n,4n): 后缀2
        [4n,5n): 前缀2
        [5n,5n+rangeToRangeOpCount): 新建立的点
        """
        self._n = n
        self.maxSize = 5 * n + rangeToRangeOpCount
        self.allocPtr = 5 * n

    def size(self) -> int:
        """新图的结点数, 前n个节点为原图的节点."""
        return self.maxSize

    def init(self, f: Callable[[int, int], None]):
        """初始化辅助边."""
        n1, n2, n3, n4 = self._n, self._n * 2, self._n * 3, self._n * 4
        for i in range(self._n):
            f(i, i + n1)  # 原始点 -> 前缀1
            f(i, i + n2)  # 原始点 -> 后缀1
            f(i + n3, i)  # 后缀2 -> 原始点
            f(i + n4, i)  # 前缀2 -> 原始点
            if i > 0:
                f(i + n1 - 1, i + n1)  # 前缀1链
                f(i + n3 - 1, i + n3)  # 后缀2链
                f(i + n2, i + n2 - 1)  # 后缀1链
                f(i + n4, i + n4 - 1)  # 前缀2链

    def add(self, fromNode: int, toNode: int, f: Callable[[int, int], None]):
        """添加有向边 from -> to."""
        f(fromNode, toNode)

    def addPrefixToPoint(self, prefixEnd: int, toNode: int, f: Callable[[int, int], None]):
        """添加从前缀区间到单点的边：[0, prefixEnd) -> to."""
        if prefixEnd <= 0:
            return
        f(self._n + prefixEnd - 1, toNode)

    def addSuffixToPoint(self, suffixStart: int, toNode: int, f: Callable[[int, int], None]):
        """添加从后缀区间到单点的边：[suffixStart, n) -> to."""
        if suffixStart >= self._n:
            return
        f(self._n * 2 + suffixStart, toNode)

    def addPointToPrefix(self, fromNode: int, prefixEnd: int, f: Callable[[int, int], None]):
        """添加从单点到前缀区间的边：from -> [0, prefixEnd)."""
        if prefixEnd <= 0:
            return
        f(fromNode, self._n * 4 + prefixEnd - 1)

    def addPointToSuffix(self, fromNode: int, suffixStart: int, f: Callable[[int, int], None]):
        """添加从单点到后缀区间的边：from -> [suffixStart, n)."""
        if suffixStart >= self._n:
            return
        f(fromNode, self._n * 3 + suffixStart)

    def addPrefixToSuffix(self, prefixEnd: int, suffixStart: int, f: Callable[[int, int], None]):
        """添加从前缀区间到后缀区间的边：[0, prefixEnd) -> [suffixStart, n)."""
        newNode = self.allocPtr
        self.allocPtr += 1
        self.addPrefixToPoint(prefixEnd, newNode, f)
        self.addPointToSuffix(newNode, suffixStart, f)

    def addSuffixToPrefix(self, suffixStart: int, prefixEnd: int, f: Callable[[int, int], None]):
        """添加从后缀区间到前缀区间的边：[suffixStart, n) -> [0, prefixEnd)."""
        newNode = self.allocPtr
        self.allocPtr += 1
        self.addSuffixToPoint(suffixStart, newNode, f)
        self.addPointToPrefix(newNode, prefixEnd, f)


if __name__ == "__main__":
    from typing import List

    class Solution:
        # 3651. 带传送的最小路径成本
        # https://leetcode.cn/problems/minimum-cost-path-with-teleportations/description/
        # 给你一个 m x n 的二维整数数组 grid 和一个整数 k。你从左上角的单元格 (0, 0) 出发，目标是到达右下角的单元格 (m - 1, n - 1)。
        # 有两种移动方式可用：
        # 普通移动：你可以从当前单元格 (i, j) 向右或向下移动，即移动到 (i, j + 1)（右）或 (i + 1, j)（下）。成本为目标单元格的值。
        # 传送：你可以从任意单元格 (i, j) 传送到任意满足 grid[x][y] <= grid[i][j] 的单元格 (x, y)；此移动的成本为 0。你最多可以传送 k 次。
        # 返回从 (0, 0) 到达单元格 (m - 1, n - 1) 的 最小 总成本。
        def minCost(self, grid: List[List[int]], k: int) -> int: ...
