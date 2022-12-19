# https://leetcode.cn/problems/count-submatrices-with-all-ones/

from typing import List


# !请你返回有多少个 子矩形 的元素全部都是 1 。
# 1 <= rows <= 150
# 1 <= columns <= 150
# 单调栈优化 O(n*m)


class Solution:
    def numSubmat(self, mat: List[List[int]]) -> int:
        """O(n*m)单调栈+dp消除重复计算"""
        ROW, COL = len(mat), len(mat[0])
        up, left = [[0] * COL for _ in range(ROW)], [[0] * COL for _ in range(ROW)]
        for row in range(ROW):
            for col in range(COL):
                if mat[row][col] == 1:
                    left[row][col] = left[row][col - 1] + 1 if col > 0 else 1
                    up[row][col] = up[row - 1][col] + 1 if row > 0 else 1

        res = 0
        for row in range(ROW):
            # 注意到如果遍历到左侧第一个小于等于他的高度的元素，那么就会`重复计算`
            # !用单调栈处理出每个位置左侧第一个小于等于的元素
            leftSmaller, stack = [-1] * COL, []
            for col in range(COL - 1, -1, -1):
                while stack and up[row][stack[-1]] >= up[row][col]:
                    leftSmaller[stack.pop()] = col
                stack.append(col)

            dp = [0] * COL  # 以每个COL为右下角的矩形个数
            for col in range(COL):
                height, width = up[row][col], left[row][col]
                pre = leftSmaller[col]
                if col - width >= pre:
                    dp[col] = width * height
                else:
                    dp[col] = (col - pre) * height + dp[pre]
                res += dp[col]
        return res


print(Solution().numSubmat(mat=[[1, 0, 1], [1, 1, 0], [1, 1, 0]]))
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
