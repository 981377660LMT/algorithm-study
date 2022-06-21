# 判断点是否在多边形上/内
# Ray casting algorithm
from typing import List


def isInPoly(polygon: List[List[int]], x: int, y: int) -> bool:
    if len(polygon) < 3:
        return False
    n = len(polygon)
    isInside = False

    for i in range(n):
        x0, y0 = polygon[i]
        x1, y1 = polygon[(i + 1) % n]

        if not min(y0, y1) < y <= max(y0, y1):
            continue

        slope = (x1 - x0) / (y1 - y0)
        x2 = x0 + (y - y0) * slope

        # 如果x2等于x 那么就是在边上
        # if x2 == x:
        #     isInside = True
        #     break

        # 点在边的左侧
        if x2 < x:
            isInside = not isInside

    return isInside


print(isInPoly(polygon=[[-3, -3], [-3, 3], [3, 3], [3, -3]], x=0, y=0))
# 我们取一条从多边形外部开始，以给定目标坐标为终点的射线，
# 并计算该射线与多边形边之间的交点数。每次光线与边相交时，
# 我们要么进入多边形，要么离开它。
# 因此，奇数交集计数表示我们在多边形内部，偶数表示我们在外部。


def isInPoly2(p: List[int], polygen: List[List[int]]) -> bool:
    px, py = p
    isInside = False
    for i, vertex in enumerate(polygen):
        ni = i + 1 if i + 1 < len(polygen) else 0
        x1, y1 = vertex
        x2, y2 = polygen[ni]
        if (x1 == px and y1 == py) or (x2 == px and y2 == py):  # if point is on vertex
            isInside = True
            break
        if min(y1, y2) < py <= max(y1, y2):  # find horizontal edges of polygon
            x = x1 + (py - y1) * (x2 - x1) / (y2 - y1)
            if x == px:  # if point is on edge
                isInside = True
                break
            elif x > px:  # if point is on left-side of line
                isInside = not isInside
    return isInside

