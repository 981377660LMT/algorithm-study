# 100332. 包含所有 1 的最小矩形面积 II
# https://leetcode.cn/problems/find-the-minimum-area-to-cover-all-ones-ii/description/
# 找到 3 个 不重叠、面积 非零 、边在水平方向和竖直方向上的矩形，并且满足 grid 中所有的 1 都在这些矩形的内部。
# 返回这些矩形面积之和的 最小 可能值。

from typing import List, Optional
from enumerateBoundingRect3 import BoundingRect, enumerateBoundingRect3


INF = int(1e18)


class Solution:
    def minimumSum(self, grid: List[List[int]]) -> int:
        res = INF
        row, col = len(grid), len(grid[0])
        rows = [FastSet(col) for _ in range(row)]
        cols = [FastSet(row) for _ in range(col)]
        for r in range(row):
            for c in range(col):
                if grid[r][c] == 1:
                    rows[r].insert(c)
                    cols[c].insert(r)

        def calc(boundingRect: BoundingRect) -> int:
            top, bottom, left, right = boundingRect
            minTop, maxBottom, minLeft, maxRight = INF, -INF, INF, -INF
            for c in range(left, right + 1):
                topCand = cols[c].next(top)
                if topCand is not None and topCand < minTop:
                    minTop = topCand
                bottomCand = cols[c].prev(bottom)
                if bottomCand is not None and bottomCand > maxBottom:
                    maxBottom = bottomCand
            for r in range(top, bottom + 1):
                leftCand = rows[r].next(left)
                if leftCand is not None and leftCand < minLeft:
                    minLeft = leftCand
                rightCand = rows[r].prev(right)
                if rightCand is not None and rightCand > maxRight:
                    maxRight = rightCand
            if minTop == INF or maxBottom == -INF or minLeft == INF or maxRight == -INF:
                return 0
            return (maxRight - minLeft + 1) * (maxBottom - minTop + 1)

        for b1, b2, b3 in enumerateBoundingRect3(row, col):
            res = min2(res, calc(b1) + calc(b2) + calc(b3))
        return res


class FastSet:
    """利用位运算寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
    初始时,所有位置都未被访问过.
    """

    __slots__ = "_n", "_lg", "_seg"

    @staticmethod
    def _trailingZeros1024(x: int) -> int:
        if x == 0:
            return 1024
        return (x & -x).bit_length() - 1

    def __init__(self, n: int) -> None:
        self._n = n
        seg = []
        while True:
            seg.append([0] * ((n + 1023) >> 10))
            n = (n + 1023) >> 10
            if n <= 1:
                break
        self._seg = seg
        self._lg = len(seg)

    def insert(self, i: int) -> None:
        for h in range(self._lg):
            self._seg[h][i >> 10] |= 1 << (i & 1023)
            i >>= 10

    def erase(self, i: int) -> None:
        for h in range(self._lg):
            self._seg[h][i >> 10] &= ~(1 << (i & 1023))
            if self._seg[h][i >> 10]:
                break
            i >>= 10

    def next(self, i: int) -> Optional[int]:
        """返回x右侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        if i < 0:
            i = 0
        if i >= self._n:
            return
        seg = self._seg
        for h in range(self._lg):
            if i >> 10 == len(seg[h]):
                break
            d = seg[h][i >> 10] >> (i & 1023)
            if d == 0:
                i = (i >> 10) + 1
                continue
            i += self._trailingZeros1024(d)
            for g in range(h - 1, -1, -1):
                i <<= 10
                i += self._trailingZeros1024(seg[g][i >> 10])
            return i

    def prev(self, i: int) -> Optional[int]:
        """返回x左侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        if i < 0:
            return
        if i >= self._n:
            i = self._n - 1
        seg = self._seg
        for h in range(self._lg):
            if i == -1:
                break
            d = seg[h][i >> 10] << (1023 - (i & 1023)) & ((1 << 1024) - 1)
            if d == 0:
                i = (i >> 10) - 1
                continue
            i += d.bit_length() - 1024
            for g in range(h - 1, -1, -1):
                i <<= 10
                i += (seg[g][i >> 10]).bit_length() - 1
            return i

    def islice(self, begin: int, end: int):
        """遍历[start,end)区间内的元素."""
        x = begin - 1
        while True:
            x = self.next(x + 1)
            if x is None or x >= end:
                break
            yield x

    def __contains__(self, i: int) -> bool:
        return self._seg[0][i >> 10] & (1 << (i & 1023)) != 0

    def __iter__(self):
        yield from self.islice(0, self._n)

    def __repr__(self):
        return f"FastSet({list(self)})"


def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


if __name__ == "__main__":
    S = Solution()
    grid = [[1, 0, 1], [1, 1, 1]]
    print(S.minimumSum(grid))  # 5
