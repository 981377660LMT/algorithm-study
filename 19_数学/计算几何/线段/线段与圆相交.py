"""
情况一、两点都在圆内。不相交
情况二、一个点在圆内,一个点在圆外。相交
情况三、两个点都在圆外

设点p1和p2均在圆外,判断线段p1p2与圆是否相交的方法
1、求出直线p1p2的一般式方程
2、用距离公式判断圆心到直线p1p2的距离是否大于半径:距离大于半径,则不相交；距离小于等于半径,执行3
3、设圆心为o,使用余弦定理判断角op1p2和角op2p1是否都为锐角,都为锐角则相交,否则不相交。
"""

from typing import Tuple


Segment = Tuple[int, int, int, int]
Circle = Tuple[int, int, int]


def isSegCircleCross(segment: Segment, circle: Circle) -> bool:
    """
    线段与圆是否相交
    !线段缩成点视为与圆相交

    https://blog.csdn.net/SongBai1997/article/details/86599879
    """
    sx1, sy1, sx2, sy2 = segment
    cx, cy, r = circle
    if sx1 == sx2 and sy1 == sy2:
        return ((sx1 - cx) * (sx1 - cx) + (sy1 - cy) * (sy1 - cy)) == r * r

    flag1 = (sx1 - cx) * (sx1 - cx) + (sy1 - cy) * (sy1 - cy) <= r * r
    flag2 = (sx2 - cx) * (sx2 - cx) + (sy2 - cy) * (sy2 - cy) <= r * r
    if flag1 and flag2:  # 两点都在圆内 不相交
        return False
    if flag1 or flag2:  # 一点在圆内一点在圆外 相交
        return True

    # 两点都在圆外

    # 将直线p1p2化为直线方程的一般式:Ax+By+C=0的形式。先化为两点式，然后由两点式得出一般式
    A = sy1 - sy2
    B = sx2 - sx1
    C = sx1 * sy2 - sx2 * sy1
    # 使用距离公式判断圆心到直线ax+by+c=0的距离是否大于半径
    dist1 = A * cx + B * cy + C
    dist1 *= dist1
    dist2 = (A * A + B * B) * r * r
    if dist1 > dist2:  # 圆心到直线距离大于半径,不相交
        return False

    # 需要判断角op1p2和角op2p1是否都为锐角,都为锐角则相交,否则不相交
    angle1 = (cx - sx1) * (sx2 - sx1) + (cy - sy1) * (sy2 - sy1)
    angle2 = (cx - sx2) * (sx1 - sx2) + (cy - sy2) * (sy1 - sy2)
    if angle1 > 0 and angle2 > 0:
        return True
    return False


if __name__ == "__main__":
    circle = (0, 0, 4)
    segment1 = (4, 0, 4, 0)  # !线段缩成点视为与圆相交
    segment2 = (4, 0, 4, 1)  # 一个点在圆内,一个点在圆外 相交
    segment3 = (-1, 1, 1, 1)  # 两点都在圆内 不相交
    segment4 = (-100, 1, 100, 1)  # 两个点都在圆外 相交

    assert isSegCircleCross(segment1, circle)
    assert isSegCircleCross(segment2, circle)
    assert not isSegCircleCross(segment3, circle)
    assert isSegCircleCross(segment4, circle)

    circle = (0, 0, 5)
    H = (3, 4, 3, 4)
    print(isSegCircleCross(H, circle))  # False
    KL = (-20, 0, -2, 0)
    print(isSegCircleCross(KL, circle))  # True
    IJ = (-10, 0, 10, 0)
    print(isSegCircleCross(IJ, circle))  # True
