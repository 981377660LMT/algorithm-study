# 点到直线距离公式
from math import sqrt
from typing import Tuple


Point = Tuple[float, float]
Line = Tuple[Point, Point]


def solve(point: Point, line: Line) -> float:
    """求点A到直线BC距离"""
    ax, ay = point
    (bx, by), (cx, cy) = line

    BAx, BAy = ax - bx, ay - by
    BCx, BCy = cx - bx, cy - by
    CAx, CAy = ax - cx, ay - cy
    CBx, CBy = bx - cx, by - cy

    # BAとBCのなす角が90度以上
    if (BAx * BCx + BAy * BCy) < 0:
        return sqrt(BAx**2 + BAy**2)

    # CAとCBのなす角が90度以上
    if (CAx * CBx + CAy * CBy) < 0:
        return sqrt(CAx**2 + CAy**2)

    area = abs(BAx * BCy - BAy * BCx)  # 平行四边形の面積
    return area / sqrt(BCx**2 + BCy**2)


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = sys.stdin.readline
    MOD = int(1e9 + 7)

    ax, ay = map(int, input().split())
    bx, by = map(int, input().split())
    cx, cy = map(int, input().split())
    print(solve((ax, ay), ((bx, by), (cx, cy))))
