from typing import List

# 整数N 的范围： [1, 500].
# mines 的最大长度为 5000.
# 除了在 mines 中给出的单元为 0，其他每个单元都是 1。
# 网格中包含 1 的最大的轴对齐加号标志是多少阶

# 764. 最大加号标志
# !预处理出每个点上下左右连续1的长度
# 6053. 统计网格图中没有被保卫/守卫的格子数-四个方向扩展


class Solution:
    def orderOfLargestPlusSign(self, n: int, mines: List[List[int]]) -> int:
        bad = set(map(tuple, mines))
        ROW, COL = n, n
        up = [[0] * COL for _ in range(ROW)]
        down = [[0] * COL for _ in range(ROW)]
        left = [[0] * COL for _ in range(ROW)]
        right = [[0] * COL for _ in range(ROW)]
        for r in range(ROW):
            for c in range(COL):
                if (r, c) not in bad:
                    up[r][c] = up[r - 1][c] + 1 if r else 1
                    left[r][c] = left[r][c - 1] + 1 if c else 1
                if (ROW - 1 - r, COL - 1 - c) not in bad:
                    down[ROW - 1 - r][COL - 1 - c] = down[ROW - r][COL - 1 - c] + 1 if r else 1
                    right[ROW - 1 - r][COL - 1 - c] = right[ROW - 1 - r][COL - c] + 1 if c else 1

        return max(
            min(up[r][c], down[r][c], left[r][c], right[r][c])
            for r in range(ROW)
            for c in range(COL)
        )


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
