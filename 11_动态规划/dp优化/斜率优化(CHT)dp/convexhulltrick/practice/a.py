from typing import *
from math import *
import sys

input = sys.stdin.readline
Point = Tuple[int, int]
INV_PHI = (sqrt(5) - 1) / 2


def ccw(a: Point, b: Point, c: Point) -> int:
    cross = (b[0] - a[0]) * (c[1] - b[1]) - (b[1] - a[1]) * (c[0] - b[0])
    if cross > 0:  # a -> b -> c is counter clockwise
        return 1
    elif cross == 0:
        return 0
    else:  # a -> b -> c is clockwise
        return -1


class ConvexHull:
    def __init__(self, p: List[Point]) -> None:
        self.p = p
        lower = self.lower = []
        upper = self.upper = []
        p.sort()
        lower.extend(p[:2])
        upper.extend(p[:2])
        for x in p[2:]:
            lower.append(x)
            upper.append(x)
            while len(lower) >= 3 and ccw(*lower[-3:]) <= 0:
                lower.pop(-2)
            while len(upper) >= 3 and ccw(*upper[-3:]) >= 0:
                upper.pop(-2)

    def get(self, a: int, b: int) -> int:
        p = self.upper if b > 0 else self.lower

        # golden-section search
        def f(i: int) -> int:
            x, y = p[i]
            return a * x + b * y

        l = 0
        r = len(p) - 1
        r2 = int(round(r * INV_PHI))
        f_r2 = f(r2)
        while abs(r - l) >= 6:
            l2 = r + int(round((l - r) * INV_PHI))
            f_l2 = f(l2)
            if f_l2 < f_r2:
                l, r = r, l2
            else:
                r, r2, f_r2 = r2, l2, f_l2
        if l > r:
            l, r = r, l
        return max(f(i) for i in range(l, r + 1))


Q = int(input())
S = []
for i in range(1, Q + 1):
    X, Y, A, B = map(int, input().split())
    S.append(ConvexHull([(X, Y)]))
    while (i & 1) == 0:
        i >>= 1
        p1 = S.pop().p
        p2 = S.pop().p
        S.append(ConvexHull(p1 + p2))
    print(max(s.get(A, B) for s in S))
