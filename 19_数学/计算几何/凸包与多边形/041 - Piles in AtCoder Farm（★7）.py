# !包围已有木桩的墙壁取到最短周长时 求需要的新木桩数目(木桩只能在整数格点上)
# n<=1e5

# !皮克定理 (ピックの定理)
# 给定顶点座标均是整点（或正方形格子点）的简单多边形，
# 皮克定理说明了其面积 S 和内部格点数目 a  边上的格点数目 b  两个数目的关系：
# !S=a+b//2-1
# 作用:求多边形内部格点数


from math import gcd
from typing import List


Point = List[int]


def calCross(A: Point, B: Point, C: Point) -> int:
    """"计算AB与AC的叉乘"""

    AB = [B[0] - A[0], B[1] - A[1]]
    AC = [C[0] - A[0], C[1] - A[1]]
    return AB[0] * AC[1] - AB[1] * AC[0]


def calConvexHull(points: List[Point]) -> List[Point]:
    """Andrew 算法 nlogn 求凸包"""
    if len(points) <= 3:
        return points

    points = sorted(points)
    stack = []

    # 寻找凸壳的下半部分
    for i in range(len(points)):
        while len(stack) >= 2 and calCross(stack[-2], stack[-1], points[i]) < 0:
            stack.pop()
        stack.append(tuple(points[i]))

    # 寻找凸壳的上半部分
    for i in range(len(points) - 1, -1, -1):
        while len(stack) >= 2 and calCross(stack[-2], stack[-1], points[i]) < 0:
            stack.pop()
        stack.append(tuple(points[i]))
    return list(set(stack))


def calPolygenArea2(points: List[Point]) -> int:
    """多边形面积的2倍"""
    assert len(points) >= 3
    res = 0
    for (x1, y1), (x2, y2) in zip(points, points[1:] + [points[0]]):
        res += x1 * y2 - x2 * y1
    return abs(res)


def countLatticePoints(points: List[Point]) -> int:
    """计算多边形边界的格点数目"""
    res = 0
    for (x1, y1), (x2, y2) in zip(points, points[1:] + [points[0]]):
        res += gcd(abs(x1 - x2), abs(y1 - y2))
    return res


n = int(input())
points = []
for _ in range(n):
    x, y = map(int, input().split())
    points.append((x, y))

convexHull = calConvexHull(points)
area2 = calPolygenArea2(convexHull)  # 凸多边形面积的2倍
outer = countLatticePoints(convexHull)  # 凸包边界的格点数目
# pick定理求多边形内部格点数
inner = (area2 - outer + 2) // 2
print(inner + outer - n)  # 新加的木桩数目

