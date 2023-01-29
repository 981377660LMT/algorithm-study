# 区间内的点到最近点的距离


from bisect import bisect_left, bisect_right
from typing import Sequence

INF = int(1e18)


def findNearest(sortedPoints: "Sequence[int]", left: int, right: int) -> int:
    """求区间[left,right]内的整点到sortedPoints内的点的最近距离"""
    count = bisect_right(sortedPoints, right) - bisect_left(sortedPoints, left)
    if count > 0:
        return 0

    res = INF
    for num in (left, right):
        pos1 = bisect_left(sortedPoints, num)
        if pos1 < len(sortedPoints):
            res = min(res, abs(sortedPoints[pos1] - num))
        pos2 = bisect_right(sortedPoints, num) - 1
        if pos2 >= 0:
            res = min(res, abs(sortedPoints[pos2] - num))
    return res
