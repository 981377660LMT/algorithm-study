# 矩阵斜向前缀和/对角线前缀和/斜对角线前缀和

from typing import List, Tuple


DIR4 = [[-1, 0], [0, -1], [1, 0], [0, 1]]
Point = Tuple[int, int]


class DiagonalPresum:
    """二维矩阵对角线前缀和"""

    __slots__ = "_preSum1", "_preSum2"

    def __init__(self, matrix: List[List[int]]):
        row, col = len(matrix), len(matrix[0])
        preSum1 = [[0] * (col + 1) for _ in range(row + 1)]  # 主对角线前缀和 ↘
        preSum2 = [[0] * (col + 1) for _ in range(row + 1)]  # 反对角线前缀和 ↙
        for i, r in enumerate(matrix):
            tmp1 = preSum1[i]
            tmp2 = preSum1[i + 1]
            tmp3 = preSum2[i]
            tmp4 = preSum2[i + 1]
            for j, v in enumerate(r):
                tmp2[j + 1] = tmp1[j] + v
                tmp4[j] = tmp3[j + 1] + v
        print(preSum1, preSum2)
        self._preSum1 = preSum1
        self._preSum2 = preSum2

    def queryDiagonal(self, leftUp: Point, rightDown: Point) -> int:
        """正对角线左上角到右下角的前缀和."""
        r1, c1, r2, c2 = *leftUp, *rightDown
        return self._preSum1[r2 + 1][c2 + 1] - self._preSum1[r1][c1]

    def queryAntiDiagonal(self, leftDown: Point, rightUp: Point) -> int:
        """副对角线左下角到右上角的前缀和"""
        r1, c1, r2, c2 = *leftDown, *rightUp
        print(self._preSum2[r1 + 1][c1], self._preSum2[r2][c2 + 1])
        return self._preSum2[r1 + 1][c1] - self._preSum2[r2][c2 + 1]


if __name__ == "__main__":
    grid = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
    dp = DiagonalPresum(grid)
    print(dp.queryDiagonal((0, 0), (2, 2)))
    print(dp.queryAntiDiagonal((2, 0), (0, 2)))

    from heapq import nlargest

    # 1878. 矩阵中最大的三个菱形和
    # LC1878 https://leetcode-cn.com/problems/get-biggest-three-rhombus-sums-in-a-grid/
    # 按照 降序 返回 grid 中三个最大的 互不相同的菱形和.
    # 如果不同的和少于三个，则将它们全部返回。
    class Solution:
        def getBiggestThree(self, grid: List[List[int]]) -> List[int]:
            """对角线前缀和"""

            def calSum(up: Point, left: Point, down: Point, right: Point) -> int:
                return (
                    (
                        dp.queryDiagonal(up, right)
                        + dp.queryAntiDiagonal(left, up)
                        + dp.queryDiagonal(left, down)
                        + dp.queryAntiDiagonal(down, right)
                    )
                    - grid[up[0]][up[1]]
                    - grid[left[0]][left[1]]
                    - grid[down[0]][down[1]]
                    - grid[right[0]][right[1]]
                )

            ROW, COL = len(grid), len(grid[0])
            res = set()

            # 对角线前缀和矩阵
            dp = DiagonalPresum(grid)
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
