# the crossing points of two circles
from typing import List, Tuple


def is_circle_cross(x1: int, y1: int, r1: int, x2: int, y2: int, r2: int) -> bool:
    """圆与圆是否有交点"""
    a = (r1 - r2) * (r1 - r2)
    b = (r1 + r2) * (r1 + r2)
    dist = (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2)
    return a <= dist <= b


def circles_cross_points(
    x1: int, y1: int, r1: int, x2: int, y2: int, r2: int
) -> List[Tuple[float, float]]:
    """两圆的交点"""
    if not is_circle_cross(x1, y1, r1, x2, y2, r2):
        return []
    rr0 = (x2 - x1) * (x2 - x1) + (y2 - y1) * (y2 - y1)
    xd = x2 - x1
    yd = y2 - y1
    rr1 = r1 * r1
    rr2 = r2 * r2
    cv = rr0 + rr1 - rr2
    sv = (4 * rr0 * rr1 - cv * cv) ** 0.5
    return sorted(
        set(
            (
                (x1 + (cv * xd - sv * yd) / (2.0 * rr0), y1 + (cv * yd + sv * xd) / (2.0 * rr0)),
                (x1 + (cv * xd + sv * yd) / (2.0 * rr0), y1 + (cv * yd - sv * xd) / (2.0 * rr0)),
            )
        )
    )


# !圆 (x1, y1, r1) 与 点 (x2, y2) 的两条切线的切点坐标
def circle_tangent_points(x1, y1, r1, x2, y2) -> List[Tuple[float, float]]:
    dd = (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2)
    r2 = (dd - r1 * r1) ** 0.5
    return circles_cross_points(x1, y1, r1, x2, y2, r2)


if __name__ == "__main__":
    p0, p1 = circle_tangent_points(0, 0, 10, 20, 15)
    print("(%.4f, %.4f) (%.4f, %.4f)" % (p0 + p1))
    # => "(-2.2991, 9.7321) (8.6991, -4.9321)"
