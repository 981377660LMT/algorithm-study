# 给定一个n*m的01矩阵，要求每一行的和为r，每一列的和为c（nr=mc且r<=m且c<=n)。

from typing import List
from math import gcd


class ZeroOneMatrixAssignProblem:
    __slots__ = "_n", "_m", "_rowSum", "_colSum", "_g", "_mg", "_rg"

    def __init__(self, n: int, m: int, rowSum: int, colSum: int):
        if not ((n * rowSum == m * colSum) and rowSum <= m and colSum <= n):
            raise ValueError()
        self._n = n
        self._m = m
        self._rowSum = rowSum
        self._colSum = colSum
        self._g = gcd(rowSum, m)
        self._mg = m // self._g
        self._rg = rowSum // self._g

    def isOne(self, i: int, j: int) -> bool:
        return (i + j) % self._mg < self._rg

    def valueOf(self, i: int, j: int) -> int:
        return int(self.isOne(i, j))

    def toMatrix(self) -> List[List[int]]:
        return [[self.valueOf(i, j) for j in range(self._m)] for i in range(self._n)]

    def __repr__(self) -> str:
        import pprint

        return pprint.pformat(self.toMatrix())


if __name__ == "__main__":
    M = ZeroOneMatrixAssignProblem(6, 4, 2, 3)
    print(M.toMatrix())
    print(M)
