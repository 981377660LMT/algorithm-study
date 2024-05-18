from typing import Callable, Generic, List, TypeVar

T = TypeVar("T")


class Presum2DMonoid(Generic[T]):
    __slot__ = "_presum", "_row", "_col", "_e"

    def __init__(self, matrix: List[List[T]], e: Callable[[], T], op: Callable[[T, T], T]):
        row, col = len(matrix), len(matrix[0])
        presum = [[e() for _ in range(col + 1)] for _ in range(row + 1)]
        for r in range(row):
            tmpsum0, tmpsum1 = presum[r], presum[r + 1]
            tmpm = matrix[r]
            for c in range(col):
                tmpsum1[c + 1] = op(op(tmpsum0[c + 1], tmpsum1[c]), tmpm[c])
        self._presum = presum
        self._row = row
        self._col = col
        self._e = e

    def query(self, x: int, y: int) -> T:
        """查询[0:x, 0:y]的聚合值."""
        if x <= 0 or y <= 0:
            return self._e()
        if x > self._row:
            x = self._row
        if y > self._col:
            y = self._col
        return self._presum[x][y]


def min2(x: int, y: int) -> int:
    return x if x < y else y


def max2(x: int, y: int) -> int:
    return x if x > y else y


if __name__ == "__main__":
    # 3148. 矩阵中的最大得分
    # https://leetcode.cn/problems/maximum-difference-score-in-a-grid/description/
    # 给你一个由 正整数 组成、大小为 m x n 的矩阵 grid。
    # 你可以从矩阵中的任一单元格移动到另一个位于正下方或正右侧的任意单元格（不必相邻）。
    # 从值为 c1 的单元格移动到值为 c2 的单元格的得分为 c2 - c1 。
    # 你可以从 任一 单元格开始，并且必须至少移动一次。
    # 返回你能得到的 最大 总得分。
    class Solution:
        def maxScore(self, grid: List[List[int]]) -> int:
            INF = int(1e18)
            S = Presum2DMonoid(grid, lambda: INF, min2)
            res = -INF
            for i in range(len(grid)):
                for j in range(len(grid[0])):
                    res = max(res, grid[i][j] - min2(S.query(i, j + 1), S.query(i + 1, j)))
            return res
