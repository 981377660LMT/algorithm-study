# 切匹萨 披萨是一个凸n边形 披萨面积的1/4记作b
# 希望能沿着两个顶点连线切 切出接近1/4面积的一块 面积记作b
# 求8*abs(a-b)的最小值
# !即求 abs(总面积-4*b)的最小值 (计算三角形面积时不除以2)

# 单调性,双指针 => 每个顶点作为left 寻找对应的right
# 4<=n<=1e5

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

# region
from fractions import Fraction
from typing import List


Point = List[int]


def calCross(p1: Point, p2: Point, p3: Point) -> int:
    ab = [p2[0] - p1[0], p2[1] - p1[1]]
    ac = [p3[0] - p1[0], p3[1] - p1[1]]
    return ab[0] * ac[1] - ab[1] * ac[0]


def calTriangleArea(p1: Point, p2: Point, p3: Point) -> Fraction:
    """三角形面积公式"""
    return Fraction(abs(calCross(p1, p2, p3)), 2)


def calTriangleArea2(p1: Point, p2: Point, p3: Point) -> int:
    """三角形面积公式2倍"""
    return abs(calCross(p1, p2, p3))


def calPolygenArea(points: List[Point]) -> Fraction:
    """多边形面积公式 鞋带定理"""
    assert len(points) >= 3
    res = Fraction(0)
    for (x1, y1), (x2, y2) in zip(points, points[1:] + [points[0]]):
        res += x1 * y2 - x2 * y1
    return Fraction(abs(res), 2)


def calPolygenArea2(points: List[Point]) -> int:
    """多边形面积的2倍 鞋带定理"""
    assert len(points) >= 3
    res = 0
    for (x1, y1), (x2, y2) in zip(points, points[1:] + [points[0]]):
        res += x1 * y2 - x2 * y1
    return abs(res)


# endregion


def main() -> None:
    n = int(input())
    points = []
    for _ in range(n):
        x, y = map(int, input().split())
        points.append((x, y))

    allSum = calPolygenArea2(points)
    res = allSum

    curSum = 0  # 8*b
    right = 0
    for left in range(n):
        while curSum < allSum:
            i, j, k = left, right % n, (right + 1) % n
            curSum += 4 * calTriangleArea2(points[i], points[j], points[k])
            right = (right + 1) % n

        cand1 = curSum
        cand2 = curSum - 4 * calTriangleArea2(
            points[left], points[(right - 1) % n], points[right % n]
        )

        res = min(res, abs(cand1 - allSum), abs(cand2 - allSum))
        i, j, k = left, (left + 1) % n, right % n  # !注意这里
        curSum -= 4 * calTriangleArea2(points[i], points[j], points[k])

    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
