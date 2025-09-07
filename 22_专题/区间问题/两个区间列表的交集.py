# EnumerateIntervalsIntersection
# https://atcoder.jp/contests/abc421/tasks/abc421_d D - RLE Moving

from typing import Generator, Optional, Tuple
from typing import List, Tuple


Interval = Tuple[int, int]


def enumerateIntervalsIntersection(
    intervals1: List[Interval], intervals2: List[Interval]
) -> Generator[Tuple[int, int, int, int], None, None]:
    """
    枚举两个已排序区间列表的交集。

    Args:
        intervals1: 第一个已排序的区间列表。
        intervals2: 第二个已排序的区间列表。

    Yields:
        一个元组 (start, end, index1, index2)，
        其中 (start, end) 是交集区间，
        index1 和 index2 是交集区间在原始列表中的索引。
    """
    i, j = 0, 0
    n1, n2 = len(intervals1), len(intervals2)
    while i < n1 and j < n2:
        s1, e1 = intervals1[i]
        s2, e2 = intervals2[j]
        intersectS = max(s1, s2)
        intersectE = min(e1, e2)
        if intersectS <= intersectE:
            yield intersectS, intersectE, i, j
        if e1 < e2:
            i += 1
        else:
            j += 1


def getIntersect(
    interval1: Tuple[int, int], interval2: Tuple[int, int]
) -> Optional[Tuple[int, int]]:
    """获取两个区间的交集。"""
    s1, e1 = interval1
    s2, e2 = interval2
    intersectS = max(s1, s2)
    intersectE = min(e1, e2)
    if intersectS <= intersectE:
        return (intersectS, intersectE)
    return None


def intersection(intervals1: List[Interval], intervals2: List[Interval]) -> List[Interval]:
    """
    双指针求两个已排序的区间列表的交集。

    Returns:
        一个列表，包含所有的交集区间。
    """
    n1, n2 = len(intervals1), len(intervals2)
    res = []
    left, right = 0, 0
    while left < n1 and right < n2:
        s1, e1 = intervals1[left]
        s2, e2 = intervals2[right]
        intersectS = max(s1, s2)
        intersectE = min(e1, e2)
        if intersectS <= intersectE:
            res.append((intersectS, intersectE))
        if e1 < e2:
            left += 1
        else:
            right += 1
    return res


def intersectionLen(intervals1: List[Interval], intervals2: List[Interval]) -> int:
    """
    计算两个已排序的区间列表相交的长度之和。
    注意：这里的长度是包含端点的整数个数，例如 [1, 3] 的长度是 3。
    """
    n1, n2 = len(intervals1), len(intervals2)
    res = 0
    left, right = 0, 0
    while left < n1 and right < n2:
        s1, e1 = intervals1[left]
        s2, e2 = intervals2[right]
        intersectS = max(s1, s2)
        intersectE = min(e1, e2)
        if intersectS <= intersectE:
            res += intersectE - intersectS + 1
        if e1 < e2:
            left += 1
        else:
            right += 1
    return res


def intersect(s1: int, e1: int, s2: int, e2: int) -> int:
    """
    计算两个区间相交的长度。
    注意：这里的长度是区间的几何长度，例如 [1, 3] 的长度是 2。
    """
    return max(0, min(e1, e2) - max(s1, s2))


if __name__ == "__main__":
    intervals1 = [(0, 2), (5, 10), (13, 23), (24, 25)]
    intervals2 = [(1, 5), (8, 12), (15, 24), (25, 26)]

    print("intersection:", intersection(intervals1, intervals2))
    # 预期输出: intersection: [(1, 2), (5, 5), (8, 10), (15, 23), (24, 24), (25, 25)]

    print("intersectionLen:", intersectionLen(intervals1, intervals2))
    # 预期输出: intersectionLen: 20

    print("intersect(0, 2, 1, 5):", intersect(0, 2, 1, 5))
    # 预期输出: intersect(0, 2, 1, 5): 1

    print("enumerateIntervalsIntersection:")
    for s, e, i, j in enumerateIntervalsIntersection(intervals1, intervals2):
        print(f"  intersection: [{s}, {e}], from intervals1[{i}] and intervals2[{j}]")
    # 预期输出:
    # enumerateIntervalsIntersection:
    #   intersection: [1, 2], from intervals1[0] and intervals2[0]
    #   intersection: [5, 5], from intervals1[1] and intervals2[0]
    #   intersection: [8, 10], from intervals1[1] and intervals2[1]
    #   intersection: [15, 23], from intervals1[2] and intervals2[2]
    #   intersection: [24, 24], from intervals1[3] and intervals2[2]
    #   intersection: [25, 25], from intervals1[3] and intervals2[3]
