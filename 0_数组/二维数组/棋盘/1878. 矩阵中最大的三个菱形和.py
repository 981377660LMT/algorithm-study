from functools import lru_cache
from heapq import nlargest
from typing import List, Tuple

# 菱形和 指的是 grid 中一个正菱形 边界 上的元素之和。本题中的菱形必须为正方形旋转45度，且四个角都在一个格子当中。
# 请你按照 降序 返回 grid 中三个最大的 `互不相同的菱形和`.如果不同的和少于三个，则将它们全部返回。
# 枚举中心点，枚举边长
# 1 <= m, n <= 100

# 比较好的方法是对角线前缀和
# 对角线的坐标的key表示为x-y x+y

DIR4 = [[-1, 0], [0, -1], [1, 0], [0, 1]]
Point = Tuple[int, int]


class Solution:
    def getBiggestThree(self, grid: List[List[int]]) -> List[int]:
        """对角线前缀和"""

        def calSum(up: Point, left: Point, down: Point, right: Point) -> int:
            r1, c1, r2, c2, r3, c3, r4, c4 = *up, *left, *down, *right
            sum1 = preSum2[r2 + 1][c2 + 1] - preSum2[r1][c1 + 2]
            sum2 = preSum1[r3 + 1][c3 + 1] - preSum1[r2][c2]
            sum3 = preSum2[r3 + 1][c3 + 1] - preSum2[r4][c4 + 2]
            sum4 = preSum1[r4 + 1][c4 + 1] - preSum1[r1][c1]
            return (
                sum1
                + sum2
                + sum3
                + sum4
                - grid[r1][c1]
                - grid[r2][c2]
                - grid[r3][c3]
                - grid[r4][c4]
            )

        ROW, COL = len(grid), len(grid[0])
        res = set()

        preSum1 = [[0] * (COL + 5) for _ in range(ROW + 5)]  # 正对角线前缀和
        preSum2 = [[0] * (COL + 5) for _ in range(ROW + 5)]  # 反对角线前缀和
        for r in range(1, ROW + 1):
            for c in range(1, COL + 1):
                preSum1[r][c] = preSum1[r - 1][c - 1] + grid[r - 1][c - 1]
                preSum2[r][c] = preSum2[r - 1][c + 1] + grid[r - 1][c - 1]

        for r in range(ROW):
            for c in range(COL):
                res.add(grid[r][c])
                for side in range(1, int(1e20)):
                    up, left, down, right = [(r + side * dr, c + side * dc) for dr, dc in DIR4]
                    if not all(
                        0 <= nr < ROW and 0 <= nc < COL for nr, nc in [up, left, down, right]
                    ):
                        break
                    res.add(calSum(up, left, down, right))

        return nlargest(3, res)


print(
    Solution().getBiggestThree(
        grid=[
            [3, 4, 5, 1, 3],
            [3, 3, 4, 2, 3],
            [20, 30, 200, 40, 10],
            [1, 5, 5, 4, 1],
            [4, 3, 2, 2, 5],
        ]
    )
)
