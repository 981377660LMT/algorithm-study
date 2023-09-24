# 2852. 所有单元格的远离程度之和
# https://leetcode.cn/problems/sum-of-remoteness-of-all-cells/description/
# 序号0起始的二维数组grid，每个正整数数表示一个区域，-1表示阻挡。
# 你可以从一个区域自由移动到任何相邻的区域。
# 对于任意点(i,j)，其遥远度定义
# 1）若为区域，则为整个矩阵内所有无法从该点到达的区域之和；
# 2）若为阻挡则为0。
# 返回整个矩阵的遥远度总和。
# n<=300

# 计算每个联通分量的和


from typing import List


class Solution:
    def sumRemoteness(self, grid: List[List[int]]) -> int:
        row, col = len(grid), len(grid[0])
        starts = []
        for i in range(row):
            for j in range(col):
                if grid[i][j] == 1:
                    starts.append((i, j))
