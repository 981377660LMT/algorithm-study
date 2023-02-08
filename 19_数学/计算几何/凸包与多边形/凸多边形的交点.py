from typing import List, Tuple


Point = Tuple[int, int]

# ps, qs: a polygon (counter-clockwise)
def is_convex_polygons_intersection(ps: List[Point], qs: List[Point]) -> bool:
    """两个凸多边形是否相交,其中多边形的顶点按逆时针方向给出."""
    pl = len(ps)
    ql = len(qs)
    i = j = 0
    while (i < pl or j < ql) and (i < 2 * pl) and (j < 2 * ql):
        px0, py0 = ps0 = ps[(i - 1) % pl]
        px1, py1 = ps1 = ps[i % pl]
        qx0, qy0 = qs0 = qs[(j - 1) % ql]
        qx1, qy1 = qs1 = qs[j % ql]

        if is_intersection(ps0, ps1, qs0, qs1):
            return True

        ax = px1 - px0
        ay = py1 - py0
        bx = qx1 - qx0
        by = qy1 - qy0

        v = ax * by - bx * ay
        va = cross3(qs0, qs1, ps1)
        vb = cross3(ps0, ps1, qs1)

        if v == 0 and va < 0 and vb < 0:
            return False
        if v == 0 and va == 0 and vb == 0:
            i += 1
        elif v >= 0:
            if vb > 0:
                i += 1
            else:
                j += 1
        else:
            if va > 0:
                j += 1
            else:
                i += 1
    return False


# ps, qs: a polygon (counter-clockwise)
def convex_polygons_intersection_points(
    ps: List[Point], qs: List[Point]
) -> List[Tuple[float, float]]:
    pl = len(ps)
    ql = len(qs)
    R = []
    i = j = 0
    i2 = j2 = 0
    mode = 0
    while (i2 < pl or j2 < ql) and (i2 < 2 * pl) and (j2 < 2 * ql):
        px0, py0 = ps0 = ps[(i - 1) % pl]
        px1, py1 = ps1 = ps[i % pl]
        qx0, qy0 = qs0 = qs[(j - 1) % ql]
        qx1, qy1 = qs1 = qs[j % ql]

        ax = px1 - px0
        ay = py1 - py0
        bx = qx1 - qx0
        by = qy1 - qy0

        v = ax * by - bx * ay
        va = cross3(qs0, qs1, ps1)
        vb = cross3(ps0, ps1, qs1)

        if is_intersection(ps0, ps1, qs0, qs1):
            if not R:
                i2 = 0
                j2 = 0
            if va > 0:
                mode = 1
            elif vb > 0:
                mode = 2
            R.append(line_cross_point(ps0, ps1, qs0, qs1))

        if v == 0 and va < 0 and vb < 0:
            break
        if v == 0 and va == 0 and vb == 0:
            if mode == 1:
                j += 1
                j2 += 1
            else:
                i += 1
                i2 += 1
        elif (ax * by - bx * ay) >= 0:
            if vb > 0:
                i += 1
                i2 += 1
            else:
                j += 1
                j2 += 1
        else:
            if va > 0:
                j += 1
                j2 += 1
            else:
                i += 1
                i2 += 1
    return R


def dot3(O, A, B):
    ox, oy = O
    ax, ay = A
    bx, by = B
    return (ax - ox) * (bx - ox) + (ay - oy) * (by - oy)


def cross3(O, A, B):
    ox, oy = O
    ax, ay = A
    bx, by = B
    return (ax - ox) * (by - oy) - (bx - ox) * (ay - oy)


def dist2(A, B):
    ax, ay = A
    bx, by = B
    return (ax - bx) * (ax - bx) + (ay - by) * (ay - by)


def is_intersection(P0, P1, Q0, Q1):
    C0 = cross3(P0, P1, Q0)
    C1 = cross3(P0, P1, Q1)
    if C0 == C1 == 0:
        E0 = dot3(P0, P1, Q0)
        E1 = dot3(P0, P1, Q1)
        if not E0 < E1:
            E0, E1 = E1, E0
        return 0 <= E1 and E0 <= dist2(P0, P1)
    D0 = cross3(Q0, Q1, P0)
    D1 = cross3(Q0, Q1, P1)
    return C0 * C1 <= 0 and D0 * D1 <= 0


def line_cross_point(P0, P1, Q0, Q1):
    x0, y0 = P0
    x1, y1 = P1
    x2, y2 = Q0
    x3, y3 = Q1
    dx0 = x1 - x0
    dy0 = y1 - y0
    dx1 = x3 - x2
    dy1 = y3 - y2

    s = (y0 - y2) * dx1 - (x0 - x2) * dy1
    sm = dx0 * dy1 - dy0 * dx1
    return (x0 + s * dx0 / sm, y0 + s * dy0 / sm) if s != 0 else (x0, y0)


if __name__ == "__main__":

    ps = [(0, 0), (1, -1), (4, 2), (3, 3)]
    qs0 = [(0, 1), (1, 1), (1, 2), (0, 2)]
    qs1 = [(0, 2), (-2, 0), (-1, -1), (1, 1)]
    qs2 = [(10, 10), (-10, 10), (-10, -10), (10, -10)]
    print(is_convex_polygons_intersection(ps, qs0))
    # => "True"
    print(is_convex_polygons_intersection(ps, qs1))
    # => "True"
    print(is_convex_polygons_intersection(ps, qs2))
    # => "False"

    ps = [(10, 0), (0, 10), (-10, 0), (0, -10)]
    qs = [(8, 8), (-10, -3), (8, -8)]
    for x, y in convex_polygons_intersection_points(ps, qs):
        print("%.03f %.03f" % (x, y))
    # =>
    # 8.000 -2.000
    # 8.000 2.000
    # 4.276 5.724
    # -8.138 -1.862
    # -5.846 -4.154
    # 3.304 -6.696
