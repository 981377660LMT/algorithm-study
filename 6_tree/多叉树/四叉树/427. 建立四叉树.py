# Definition for a QuadTree node.
from typing import List


class Node:
    def __init__(self, val, isLeaf, topLeft, topRight, bottomLeft, bottomRight):
        self.val = val  # 储存叶子结点所代表的区域的值。1 对应 True，0 对应 False；
        self.isLeaf = isLeaf  # 当这个节点是一个叶子结点时为 True，如果它有 4 个子节点则为 False
        self.topLeft = topLeft
        self.topRight = topRight
        self.bottomLeft = bottomLeft
        self.bottomRight = bottomRight


class PreSumMatrix:
    """二维前缀和模板(矩阵不可变)"""

    def __init__(self, A: List[List[int]]):
        m, n = len(A), len(A[0])

        # 前缀和数组
        preSum = [[0] * (n + 1) for _ in range(m + 1)]
        for r in range(m):
            for c in range(n):
                preSum[r + 1][c + 1] = A[r][c] + preSum[r][c + 1] + preSum[r + 1][c] - preSum[r][c]
        self.preSum = preSum

    def sumRegion(self, r1: int, c1: int, r2: int, c2: int) -> int:
        """查询sum(A[r1:r2+1, c1:c2+1])的值::

        preSumMatrix.sumRegion(0, 0, 2, 2) # 左上角(0, 0)到右下角(2, 2)的值
        """
        return (
            self.preSum[r2 + 1][c2 + 1]
            - self.preSum[r2 + 1][c1]
            - self.preSum[r1][c2 + 1]
            + self.preSum[r1][c1]
        )


class Solution:
    def construct(self, grid: List[List[int]]) -> 'Node':
        def dfs(r1: int, c1: int, r2: int, c2: int) -> Node:
            """类似于克隆图"""
            val = preSumMatrix.sumRegion(r1, c1, r2 - 1, c2 - 1)

            # 叶子结点
            if val == 0:
                return Node(False, True, None, None, None, None)
            elif val == (r2 - r1) * (c2 - c1):
                return Node(True, True, None, None, None, None)

            # 非叶子结点
            node1, node2, node3, node4 = (
                dfs(r1, c1, (r1 + r2) // 2, (c1 + c2) // 2),
                dfs(r1, (c1 + c2) // 2, (r1 + r2) // 2, c2),
                dfs((r1 + r2) // 2, c1, r2, (c1 + c2) // 2),
                dfs((r1 + r2) // 2, (c1 + c2) // 2, r2, c2),
            )
            return Node(False, False, node1, node2, node3, node4)

        preSumMatrix = PreSumMatrix(grid)

        # 左闭右开,递归时处理边界坐标会方便一些
        return dfs(0, 0, len(grid), len(grid[0]))

