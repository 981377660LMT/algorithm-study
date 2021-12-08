from typing import List

# 整数N 的范围： [1, 500].
# mines 的最大长度为 5000.
# 除了在 mines 中给出的单元为 0，其他每个单元都是 1。
# 网格中包含 1 的最大的轴对齐加号标志是多少阶
class Solution:
    def orderOfLargestPlusSign(self, n: int, mines: List[List[int]]) -> int:
        grid = [[min(r, n - 1 - r, c, n - 1 - c) + 1 for c in range(n)] for r in range(n)]
        for (x, y) in mines:
            for i in range(n):
                grid[i][y] = min(grid[i][y], abs(i - x))
                grid[x][i] = min(grid[x][i], abs(i - y))

        return max([max(row) for row in grid])


print(Solution().orderOfLargestPlusSign(n=5, mines=[[4, 2]]))
# 阶 1:
# 000
# 010
# 000

# 阶 2:
# 00000
# 00100
# 01110
# 00100
# 00000

# 阶 3:
# 0000000
# 0001000
# 0001000
# 0111110
# 0001000
# 0001000
# 0000000


# 总结：没有炸弹时只受边缘限制
# 有炸弹需要考虑和炸弹同行列的情形
# g[x][y] is the largest plus sign allowed centered at position (x, y). When no mines are presented, it is only limited by the boundary and should be something similar to
# 1 1 1 1 1
# 1 2 2 2 1
# 1 2 3 2 1
# 1 2 2 2 1
# 1 1 1 1 1
