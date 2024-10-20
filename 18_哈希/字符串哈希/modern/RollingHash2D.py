from typing import List, Optional
from random import randint

MOD61 = (1 << 61) - 1


class RollingHash2D:
    __slots__ = ("_base1", "_base2", "_power1", "_power2")

    def __init__(self, base1: Optional[int] = None, base2: Optional[int] = None):
        if base1 is None:
            base1 = randint(1, MOD61 - 1)
        if base2 is None:
            base2 = randint(1, MOD61 - 1)
        self._base1 = base1
        self._base2 = base2
        self._power1 = [1]
        self._power2 = [1]

    def build(self, grid: List[List[int]]) -> List[List[int]]:
        row = len(grid)
        col = len(grid[0]) if row else 0
        table = [[0] * (col + 1) for _ in range(row + 1)]
        for i, r in enumerate(grid):
            table1, table2 = table[i + 1], table[i]
            for j, v in enumerate(r):
                table1[j + 1] = (table1[j] * self._base2 + v + 1) % MOD61
            for j in range(col + 1):
                table1[j] = (table1[j] + table2[j] * self._base1) % MOD61
        while len(self._power1) <= row:
            self._power1.append(self._power1[-1] * self._base1 % MOD61)
        while len(self._power2) <= col:
            self._power2.append(self._power2[-1] * self._base2 % MOD61)
        return table

    def query(self, table: List[List[int]], xStart: int, xEnd: int, yStart: int, yEnd: int) -> int:
        a = (table[xEnd][yEnd] - table[xEnd][yStart] * self._power2[yEnd - yStart]) % MOD61
        b = (table[xStart][yEnd] - table[xStart][yStart] * self._power2[yEnd - yStart]) % MOD61
        return (a - b * self._power1[xEnd - xStart]) % MOD61


if __name__ == "__main__":
    # https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=ALDS1_14_C
    # !检测二维矩阵中是否存在子矩阵与给定的特征矩阵相同,输出左上角坐标
    row1, col1 = list(map(int, input().split()))
    grid1 = [list(map(ord, input())) for _ in range(row1)]
    winRow, winCol = list(map(int, input().split()))
    grid2 = [list(map(ord, input())) for _ in range(winRow)]

    H = RollingHash2D()
    table1 = H.build(grid1)
    table2 = H.build(grid2)

    target = H.query(table2, 0, winRow, 0, winCol)

    for i in range(row1 - winRow + 1):
        for j in range(col1 - winCol + 1):
            if H.query(table1, i, i + winRow, j, j + winCol) == target:
                print(i, j)
