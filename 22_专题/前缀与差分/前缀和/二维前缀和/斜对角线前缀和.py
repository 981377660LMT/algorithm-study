from heapq import nlargest
from typing import List, Tuple


DIR4 = [[-1, 0], [0, -1], [1, 0], [0, 1]]
Point = Tuple[int, int]


class DiagonalPresum:
    """二维矩阵对角线前缀和"""

    def __init__(self, matrix: List[List[int]]):
        row, col = len(matrix), len(matrix[0])

        # 前缀和数组
        preSum1 = [[0] * (col + 2) for _ in range(row + 2)]
        preSum2 = [[0] * (col + 2) for _ in range(row + 2)]
        for r in range(1, row + 1):
            for c in range(1, col + 1):
                preSum1[r][c] = preSum1[r - 1][c - 1] + matrix[r - 1][c - 1]
                preSum2[r][c] = preSum2[r - 1][c + 1] + matrix[r - 1][c - 1]

        self._preSum1 = preSum1
        self._preSum2 = preSum2

    def cal(self, leftUp: Point, rightDown: Point) -> int:
        """正对角线右下角到左上角的前缀和"""
        r1, c1, r2, c2 = *leftUp, *rightDown
        return self._preSum1[r2 + 1][c2 + 1] - self._preSum1[r1][c1]

    def rCal(self, leftDown: Point, rightUp: Point) -> int:
        """副对角线左下角到右上角的前缀和"""
        r1, c1, r2, c2 = *leftDown, *rightUp
        return self._preSum2[r1 + 1][c1 + 1] - self._preSum2[r2][c2 + 2]


if __name__ == '__main__':
    # 1878. 矩阵中最大的三个菱形和
    class Solution:
        def getBiggestThree(self, grid: List[List[int]]) -> List[int]:
            """对角线前缀和"""

            def calSum(up: Point, left: Point, down: Point, right: Point) -> int:
                return (
                    (
                        dp.cal(up, right)
                        + dp.rCal(left, up)
                        + dp.cal(left, down)
                        + dp.rCal(down, right)
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

