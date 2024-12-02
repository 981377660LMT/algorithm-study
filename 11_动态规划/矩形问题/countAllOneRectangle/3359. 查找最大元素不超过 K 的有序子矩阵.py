# 3359. 查找最大元素不超过 K 的有序子矩阵
# https://leetcode.cn/problems/find-sorted-submatrices-with-maximum-element-at-most-k/solutions/2997274/dan-diao-zhan-jie-fa-by-fzldq-8067/
# 给定一个大小为 m x n 的二维矩阵 grid。同时给定一个 非负整数 k。
#
# 返回满足下列条件的 grid 的子矩阵数量：
#
# 子矩阵中最大的元素 小于等于 k。
# 子矩阵的每一行都以 非递增 顺序排序。
# 矩阵的子矩阵 (x1, y1, x2, y2) 是通过选择所有满足 x1 <= x <= x2 且 y1 <= y <= y2 的 grid[x][y] 元素组成的矩阵。
#
# !参考 https://leetcode.cn/problems/count-submatrices-with-all-ones/description/

from typing import List


class Solution:
    def countSubmatrices(self, grid: List[List[int]], k: int) -> int:
        n, m = len(grid), len(grid[0])
        rowDp = [[0] * m for _ in range(n)]
        for i in range(n):
            for j in range(m):
                if grid[i][j] <= k:  # check
                    rowDp[i][j] = (
                        rowDp[i][j - 1] + 1 if j and grid[i][j - 1] >= grid[i][j] else 1
                    )  # checkRow

        res = 0
        for j in range(m):
            stack = []
            total = 0
            for i in range(n):
                height = 1
                while stack and stack[-1][0] > rowDp[i][j]:
                    # 弹出的时候要减去多于的答案
                    total -= stack[-1][1] * (stack[-1][0] - rowDp[i][j])
                    height += stack[-1][1]
                    stack.pop()
                total += rowDp[i][j]
                res += total
                stack.append((rowDp[i][j], height))
        return res
