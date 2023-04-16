from bisect import bisect_left
from typing import Tuple, List


def rectangleUnion(rectangles: List[Tuple[int, int, int, int]]) -> int:
    """矩形面积并."""
    ys = []
    for _, y1, _, y2 in rectangles:
        ys.append(y1)
        ys.append(y2)
    ys = sorted(set(ys))
    n = len(ys)
    C, A = [0] * 8 * n, [0] * 8 * n

    def min(a: int, b: int) -> int:
        return a if a < b else b

    def max(a: int, b: int) -> int:
        return a if a > b else b

    def aux(a: int, b: int, c: int, l: int, r: int, k: int) -> None:
        a, b = max(a, l), min(b, r)
        if a >= b:
            return
        if a == l and b == r:
            C[k] += c
        else:
            aux(a, b, c, l, (l + r) // 2, 2 * k + 1)
            aux(a, b, c, (l + r) // 2, r, 2 * k + 2)
        if C[k] != 0:
            A[k] = ys[r] - ys[l]
        else:
            A[k] = A[2 * k + 1] + A[2 * k + 2]

    events = []
    for x1, y1, x2, y2 in rectangles:
        l = bisect_left(ys, y1)
        h = bisect_left(ys, y2)
        events.append((x1, l, h, 1))
        events.append((x2, l, h, -1))
    events.sort(key=lambda e: (e[0], -e[3]))

    area, prev = 0, 0
    for x, l, h, c in events:
        area += (x - prev) * A[0]
        prev = x
        aux(l, h, c, 0, n, 0)
    return area


if __name__ == "__main__":
    # https://leetcode.cn/problems/rectangle-area-ii/
    # 850. 矩形面积 II
    class Solution:
        def rectangleArea(self, rectangles: List[List[int]]) -> int:
            rec = [(x1, y1, x2, y2) for x1, y1, x2, y2 in rectangles]
            return rectangleUnion(rec) % (10**9 + 7)
