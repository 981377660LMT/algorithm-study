from typing import List

# 每个值 v = grid[i][j] 表示 v 个正方体叠放在单元格 (i, j) 上。
# 返回所有三个投影的总面积。
class Solution:
    def projectionArea(self, grid: List[List[int]]) -> int:
        return (
            sum([max(row) for row in grid])
            + sum([max(col) for col in zip(*grid)])
            + sum(h > 0 for h in sum(grid, []))
        )


# 解题思路
# X的投影
# sum([max(a) for a in zip(*grid)])
# Y的投影
# sum([max(b) for b in grid])
# Z的投影(flatten the nested list)
# sum(k>0 for k in sum(grid, []))

# 1 2
# 3 4

# 2+4+3+4+4
