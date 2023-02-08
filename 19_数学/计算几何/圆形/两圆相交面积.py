# 两圆相交面积
# https://tjkendev.github.io/procon-library/python/geometry/circles_intersection_area.html


from math import pi, atan2

# 两圆相交面积


def circles_intersection_area(x1: int, y1: int, r1: int, x2: int, y2: int, r2: int) -> float:
    """两圆相交面积"""
    dd = (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2)

    if (r1 + r2) * (r1 + r2) <= dd:
        return 0.0

    if dd <= (r1 - r2) * (r1 - r2):
        min_ = r1 if r1 < r2 else r2
        return pi * min_ * min_

    p1 = r1 * r1 - r2 * r2 + dd
    p2 = r2 * r2 - r1 * r1 + dd

    S1 = r1 * r1 * atan2((4 * dd * r1 * r1 - p1 * p1) ** 0.5, p1)
    S2 = r2 * r2 * atan2((4 * dd * r2 * r2 - p2 * p2) ** 0.5, p2)
    S0 = (4 * dd * r1 * r1 - p1 * p1) ** 0.5 / 2

    return S1 + S2 - S0


if __name__ == "__main__":
    print(circles_intersection_area(0, 1, 1, 1, 0, 1))
    # => "0.5707963267948966"
    # https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_7_E&lang=ja
    def solve():
        x1, y1, r1 = map(int, input().split())
        x2, y2, r2 = map(int, input().split())
        print("%.16f\n" % circles_intersection_area(x1, y1, r1, x2, y2, r2))

    solve()
