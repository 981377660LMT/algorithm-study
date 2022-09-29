# 线段相交/线段交点
# !判断两线段是否相交(是否存在坐标相同的点)
# 也可参考
# https://leetcode.cn/circle/discuss/fC4N4x/

from typing import Tuple, Union

Segment = Tuple[int, int, int, int]


def cross(x1: int, y1: int, x2: int, y2: int) -> int:
    """内积"""
    return x1 * y2 - y1 * x2


def isSegCross(segment1: Segment, segment2: Segment) -> bool:
    """线段 (x1,y1,x2,y2) 与 (x3,y3,x4,y4) 是否相交"""
    x1, y1, x2, y2 = segment1
    x3, y3, x4, y4 = segment2
    res1 = cross(x2 - x1, y2 - y1, x3 - x1, y3 - y1)  # 2 1 3
    res2 = cross(x2 - x1, y2 - y1, x4 - x1, y4 - y1)  # 2 1 4
    res3 = cross(x4 - x3, y4 - y3, x1 - x3, y1 - y3)  # 4 3 1
    res4 = cross(x4 - x3, y4 - y3, x2 - x3, y2 - y3)  # 4 3 2

    # 线段共线
    if res1 == 0 and res2 == 0 and res3 == 0 and res4 == 0:
        A, B, C, D = (x1, y1), (x2, y2), (x3, y3), (x4, y4)
        A, B = sorted((A, B))
        C, D = sorted((C, D))
        return max(A, C) <= min(B, D)

    # 不共线
    canAB = (res1 >= 0 and res2 <= 0) or (res1 <= 0 and res2 >= 0)  # 線分 AB が点 C, D を分けるか？
    canCD = (res3 >= 0 and res4 <= 0) or (res3 <= 0 and res4 >= 0)  # 線分 CD が点 A, B を分けるか？
    return canAB and canCD


def intersection(
    segment1: Segment, segment2: Segment
) -> Union[Tuple[float, float], Tuple[None, None]]:
    """线段 (x1,y1,x2,y2) 与 (x3,y3,x4,y4) 的交点

    浮点型误差不超过10^-6。若有多个交点(线段重叠)则返回 X 值最小的点,X 坐标相同则返回 Y 值最小的点。
    """
    x1, y1, x2, y2 = segment1
    x3, y3, x4, y4 = segment2

    up1 = cross(x4 - x2, x4 - x3, y4 - y2, y4 - y3)
    up2 = cross(x1 - x2, x4 - x2, y1 - y2, y4 - y2)
    down = cross(x1 - x2, x4 - x3, y1 - y2, y4 - y3)

    # 唯一解
    if down != 0:
        a, b = up1 / down, up2 / down
        isValid = 0 <= a <= 1 and 0 <= b <= 1
        return (a * x1 + (1 - a) * x2, a * y1 + (1 - a) * y2) if isValid else (None, None)

    # 无解(平行)
    if up1 or up2:
        return (None, None)

    # 无穷解(重叠)
    A, B, C, D = (x1, y1), (x2, y2), (x3, y3), (x4, y4)
    A, B = sorted((A, B))
    C, D = sorted((C, D))
    if B < C or D < A:
        return (None, None)
    return C if A < C else A


if __name__ == "__main__":
    x1, y1 = map(int, input().split())
    x2, y2 = map(int, input().split())
    x3, y3 = map(int, input().split())
    x4, y4 = map(int, input().split())
    if isSegCross(x1, y1, x2, y2, x3, y3, x4, y4):
        print("Yes")
    else:
        print("No")
