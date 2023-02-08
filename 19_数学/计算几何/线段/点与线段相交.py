# 点可以视为特殊的线段,就变成了线段与线段相交
from typing import Tuple

Segment = Tuple[int, int, int, int]


def isPointOnSegment(point: Tuple[int, int], segment: Segment) -> bool:
    """点是否在线段上"""

    def cross(x1: int, y1: int, x2: int, y2: int) -> int:
        return x1 * y2 - y1 * x2

    def isSegCross(segment1: Segment, segment2: Segment) -> bool:
        x1, y1, x2, y2 = segment1
        x3, y3, x4, y4 = segment2
        res1 = cross(x2 - x1, y2 - y1, x3 - x1, y3 - y1)  # 2 1 3
        res2 = cross(x2 - x1, y2 - y1, x4 - x1, y4 - y1)  # 2 1 4
        res3 = cross(x4 - x3, y4 - y3, x1 - x3, y1 - y3)  # 4 3 1
        res4 = cross(x4 - x3, y4 - y3, x2 - x3, y2 - y3)  # 4 3 2

        if res1 == 0 and res2 == 0 and res3 == 0 and res4 == 0:
            A, B, C, D = (x1, y1), (x2, y2), (x3, y3), (x4, y4)
            A, B = sorted((A, B))
            C, D = sorted((C, D))
            return max(A, C) <= min(B, D)

        canAB = (res1 >= 0 and res2 <= 0) or (res1 <= 0 and res2 >= 0)  # 線分 AB が点 C, D を分けるか？
        canCD = (res3 >= 0 and res4 <= 0) or (res3 <= 0 and res4 >= 0)  # 線分 CD が点 A, B を分けるか？
        return canAB and canCD

    x, y = point
    return isSegCross((x, y, x, y), segment)


# 点と線分の交差判定: 点 (1, 1)
print(isPointOnSegment((1, 1), (0, 0, 3, 3)))
# => "True"
print(isPointOnSegment((1, 1), (-3, -2, 1, 1)))
# => "True"
