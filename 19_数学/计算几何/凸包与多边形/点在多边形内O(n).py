# Ray casting algorithm
# !判断点是否在多边形内,如果点不是逆时针给出的多边形的顶点，则需要先极角排序.


# !我们取一条从多边形外部开始，以给定目标坐标为终点的射线，
# 并计算该射线与多边形边之间的交点数。每次光线与边相交时，
# 我们要么进入多边形，要么离开它。
# 因此，奇数交集计数表示我们在多边形内部，偶数表示我们在外部。

from typing import List, Tuple

Point = Tuple[int, int]


def inside_polygon2(p0: Point, poly: List[Point]) -> int:
    """O(n) 判断点是否在多边形内,其中多边形的顶点按逆时针方向给出.

    Returns:
        0: 在多边形外部
        1: 在多边形边上
        2: 在多边形内部
    """
    cnt = 0
    n = len(poly)
    x, y = p0
    for i in range(n):
        x0, y0 = poly[i - 1]
        x1, y1 = poly[i]
        x0 -= x
        y0 -= y
        x1 -= x
        y1 -= y

        cv = x0 * x1 + y0 * y1
        sv = x0 * y1 - x1 * y0
        if sv == 0 and cv <= 0:
            # a point is on a segment
            return 1

        if not y0 < y1:
            x0, x1 = x1, x0
            y0, y1 = y1, y0

        if y0 <= 0 < y1 and x0 * (y1 - y0) > y0 * (x1 - x0):
            cnt += 1
    return 2 if cnt % 2 else 0


if __name__ == "__main__":

    n = int(input())
    polygen = list(tuple(map(int, input().split())) for _ in range(n))
    q = int(input())
    for _ in range(q):
        x, y = map(int, input().split())
        print(inside_polygon2((x, y), polygen))
