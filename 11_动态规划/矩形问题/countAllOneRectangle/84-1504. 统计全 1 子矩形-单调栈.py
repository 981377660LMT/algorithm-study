# https://leetcode.cn/problems/count-submatrices-with-all-ones/
# !请你返回有多少个 子矩形 的元素全部都是 1 。
# 1 <= rows <= 150
# 1 <= columns <= 150
# 单调栈优化 O(n*m)

from typing import List


class Solution:
    def numSubmat(self, mat: List[List[int]]) -> int:
        """O(n*m)单调栈+dp消除重复计算"""
        n, m = len(mat), len(mat[0])
        rowDp = [[0] * m for _ in range(n)]
        for i in range(n):
            for j in range(m):
                if mat[i][j] == 1:
                    rowDp[i][j] = rowDp[i][j - 1] + 1 if j else 1

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


assert Solution().numSubmat(mat=[[1, 0, 1], [1, 1, 0], [1, 1, 0]]) == 13
assert Solution().numSubmat(mat=[[0, 1, 1, 0], [0, 1, 1, 1], [1, 1, 1, 0]]) == 24
# 输出：13
# 解释：
# 有 6 个 1x1 的矩形。
# 有 2 个 1x2 的矩形。
# 有 3 个 2x1 的矩形。
# 有 1 个 2x2 的矩形。
# 有 1 个 3x1 的矩形。
# 矩形数目总共 = 6 + 2 + 3 + 1 + 1 = 13 。
