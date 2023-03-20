# https://leetcode.cn/problems/count-submatrices-with-all-ones/
# 单调栈优化 O(n*m)

from typing import List

# 1 <= rows <= 150
# 1 <= columns <= 150

# !请你返回有多少个 子矩形 的元素全部都是 1 。(子矩形的个数)

# O(n ^ 2 * m)
# 用 left[i][j] 表示矩阵中位置为 (i, j) 的元素前面有多少个连续的 1，
# !然后遍历矩阵中的每一个元素，计算以这个元素`为右下角`的全 1 子矩阵有多少个。
# 这里以每一个遍历到的元素`为右下角`计算全 1 子矩阵，既不会遗漏也不会重复。
# mat = [[1,0,1],
#        [1,1,0],
#        [1,1,0]]
# becomes
# mat = [[1,0,1],
#        [2,1,0],
#        [3,2,0]]


class Solution:
    def numSubmat(self, mat: List[List[int]]) -> int:
        """O(n^2*m)预处理"""
        ROW, COL = len(mat), len(mat[0])
        up, left = [[0] * COL for _ in range(ROW)], [[0] * COL for _ in range(ROW)]
        for row in range(ROW):
            for col in range(COL):
                if mat[row][col] == 1:
                    left[row][col] = left[row][col - 1] + 1 if col > 0 else 1
                    up[row][col] = up[row - 1][col] + 1 if row > 0 else 1

        res = 0
        for row in range(ROW):
            for col in range(COL):
                height, width = up[row][col], left[row][col]
                for w in range(width):  # !这一段能否加速? => 注意到如果遍历到第一个左侧比当前位置小的元素，那么就会重复计算
                    height = min(height, up[row][col - w])
                    res += height
        return res


print(Solution().numSubmat(mat=[[1, 0, 1], [1, 1, 0], [1, 1, 0]]))
assert Solution().numSubmat(mat=[[1, 0, 1], [1, 1, 0], [1, 1, 0]]) == 13
print(Solution().numSubmat(mat=[[0, 1, 1, 0], [0, 1, 1, 1], [1, 1, 1, 0]]))
assert Solution().numSubmat(mat=[[0, 1, 1, 0], [0, 1, 1, 1], [1, 1, 1, 0]]) == 24
# 输出：13
# 解释：
# 有 6 个 1x1 的矩形。
# 有 2 个 1x2 的矩形。
# 有 3 个 2x1 的矩形。
# 有 1 个 2x2 的矩形。
# 有 1 个 3x1 的矩形。
# 矩形数目总共 = 6 + 2 + 3 + 1 + 1 = 13 。
