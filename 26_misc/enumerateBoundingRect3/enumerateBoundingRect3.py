# 切蛋糕问题/矩形切割问题

from typing import Generator, List, Tuple

BoundingRect = Tuple[int, int, int, int]  # (top,bottom,left,right)


def enumerateBoundingRect3(
    row: int, col: int
) -> Generator[Tuple[BoundingRect, BoundingRect, BoundingRect], None, None]:
    """
    给定一个row*col的矩阵,分割成3个不重合的矩形,返回所有可能的分割方法.
    https://leetcode.cn/circle/discuss/6PHVvJ/
    """
    # 三横
    for r1 in range(row - 2):
        for r2 in range(r1 + 1, row - 1):
            yield ((0, r1, 0, col - 1), (r1 + 1, r2, 0, col - 1), (r2 + 1, row - 1, 0, col - 1))

    # 三竖
    for c1 in range(col - 2):
        for c2 in range(c1 + 1, col - 1):
            yield ((0, row - 1, 0, c1), (0, row - 1, c1 + 1, c2), (0, row - 1, c2 + 1, col - 1))

    # 先一横 然后再切一竖
    for r in range(row - 1):
        for c in range(col - 1):
            yield ((0, r, 0, c), (0, r, c + 1, col - 1), (r + 1, row - 1, 0, col - 1))
        for c in range(col - 1):
            yield (
                (0, r, 0, col - 1),
                (r + 1, row - 1, c + 1, col - 1),
                (r + 1, row - 1, 0, c),
            )

    # 先一竖 再切一横
    for c in range(col - 1):
        for r in range(row - 1):
            yield ((0, r, 0, c), (r + 1, row - 1, 0, c), (0, row - 1, c + 1, col - 1))
        for r in range(row - 1):
            yield (
                (0, row - 1, 0, c),
                (0, r, c + 1, col - 1),
                (r + 1, row - 1, c + 1, col - 1),
            )


if __name__ == "__main__":
    INF = int(1e18)

    def min2(a: int, b: int) -> int:
        return a if a < b else b

    def max2(a: int, b: int) -> int:
        return a if a > b else b

    class Solution:
        def minimumSum(self, grid: List[List[int]]) -> int:
            def f(boundingRect: BoundingRect) -> int:
                top, bottom, left, right = boundingRect
                minTop, maxBottom, minLeft, maxRight = INF, -INF, INF, -INF
                for r in range(top, bottom + 1):
                    for c in range(left, right + 1):
                        minTop = min2(minTop, r)
                        maxBottom = max2(maxBottom, r)
                        minLeft = min2(minLeft, c)
                        maxRight = max2(maxRight, c)
                if minTop == INF or maxBottom == -INF or minLeft == INF or maxRight == -INF:
                    return 0
                return (maxBottom - minTop + 1) * (maxRight - minLeft + 1)

            row, col = len(grid), len(grid[0])
            return min(f(b1) + f(b2) + f(b3) for b1, b2, b3 in enumerateBoundingRect3(row, col))
