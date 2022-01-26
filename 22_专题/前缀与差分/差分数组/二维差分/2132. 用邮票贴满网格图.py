from typing import List
from 二维差分模板 import DiffMatrix, PreSumMatrix


class Solution:
    def possibleToStamp(self, grid: List[List[int]], h: int, w: int) -> bool:
        row, col = len(grid), len(grid[0])
        preSumMatrix = PreSumMatrix(grid)
        diffMatrix = DiffMatrix(grid)
        for r in range(row):
            for c in range(col):
                if (
                    r + h - 1 < row
                    and c + w - 1 < col
                    and preSumMatrix.sumRegion(r, c, r + h - 1, c + w - 1) == 0
                ):
                    diffMatrix.add(r, c, r + h - 1, c + w - 1, 1)
        diffMatrix.update()

        for r in range(row):
            for c in range(col):
                if grid[r][c] == 0 and diffMatrix.matrix[r][c] == 0:
                    return False
        return True


if __name__ == '__main__':

    assert (
        Solution().possibleToStamp(
            grid=[[1, 0, 0, 0], [1, 0, 0, 0], [1, 0, 0, 0], [1, 0, 0, 0], [1, 0, 0, 0]], h=4, w=3
        )
        == True
    )

    assert (
        Solution().possibleToStamp(
            grid=[[1, 0, 0, 0], [0, 1, 0, 0], [0, 0, 1, 0], [0, 0, 0, 1]], h=2, w=2
        )
        == False
    )

