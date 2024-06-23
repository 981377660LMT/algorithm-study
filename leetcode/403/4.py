# 给你一个二维 二进制 数组 grid。你需要找到 3 个 不重叠、面积 非零 、边在水平方向和竖直方向上的矩形，并且满足 grid 中所有的 1 都在这些矩形的内部。
# 返回这些矩形面积之和的 最小 可能值。
# 注意，这些矩形可以相接。


# TODO: fromBounding
# TODO: enumerateBoundingOfThreeSquares

import pprint
from typing import List


INF = int(1e20)


def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minimumSum(self, grid: List[List[int]]) -> int:
        res = INF
        for _ in range(4):
            grid = self._rotate(grid)
            res = min2(res, self._solve1(grid))
            res = min2(res, self._solve2(grid))
        return res

    def _solve1(self, grid: List[List[int]]) -> int:
        """竖着切两刀(左、中、右)."""
        res = INF
        row, col = len(grid), len(grid[0])
        cols = self._collectCols(grid)
        for split1 in range(len(cols)):
            area1 = self._minimumBoundingArea(0, row - 1, 0, cols[split1], grid)
            for split2 in range(split1 + 1, len(cols)):
                area2 = self._minimumBoundingArea(0, row - 1, cols[split1] + 1, cols[split2], grid)
                area3 = self._minimumBoundingArea(0, row - 1, cols[split2] + 1, col - 1, grid)
                res = min2(res, area1 + area2 + area3)
        return res

    def _solve2(self, grid: List[List[int]]) -> int:
        """横一刀竖一刀(上、左下、右下)."""
        res = INF
        row, col = len(grid), len(grid[0])
        cols = self._collectCols(grid)
        rows = self._collectRows(grid)
        for rSplit in range(len(rows)):
            area1 = self._minimumBoundingArea(0, rows[rSplit], 0, col - 1, grid)
            for cSplit in range(len(cols)):
                area2 = self._minimumBoundingArea(rows[rSplit] + 1, row - 1, 0, cols[cSplit], grid)
                area3 = self._minimumBoundingArea(
                    rows[rSplit] + 1, row - 1, cols[cSplit] + 1, col - 1, grid
                )
                res = min2(res, area1 + area2 + area3)
        return res

    def _minimumBoundingArea(
        self, top: int, bottom: int, left: int, right: int, grid: List[List[int]]
    ) -> int:
        minTop, maxBottom, minLeft, maxRight = INF, -INF, INF, -INF
        for r in range(top, bottom + 1):
            for c in range(left, right + 1):
                if grid[r][c] == 1:
                    minTop = min2(minTop, r)
                    maxBottom = max2(maxBottom, r)
                    minLeft = min2(minLeft, c)
                    maxRight = max2(maxRight, c)
        if minTop == INF or maxBottom == -INF or minLeft == INF or maxRight == -INF:
            return 0
        return (maxRight - minLeft + 1) * (maxBottom - minTop + 1)

    def _collectCols(self, grid: List[List[int]]) -> List[int]:
        res = []
        for c, col in enumerate(zip(*grid)):
            for v in col:
                if v == 1:
                    res.append(c)
                    break
        return res

    def _collectRows(self, grid: List[List[int]]) -> List[int]:
        res = []
        for r, row in enumerate(grid):
            for v in row:
                if v == 1:
                    res.append(r)
                    break
        return res

    def _rotate(self, grid: List[List[int]]) -> List[List[int]]:
        return [list(col[::-1]) for col in zip(*grid)]


if __name__ == "__main__":
    S = Solution()
    grid = [[0, 0, 0, 1, 0], [0, 0, 0, 0, 0], [0, 1, 0, 0, 1], [0, 0, 0, 0, 0], [0, 0, 1, 0, 0]]
    print(S.minimumSum(grid))
    pprint.pprint(grid)
