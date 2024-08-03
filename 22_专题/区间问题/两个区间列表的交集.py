from typing import Optional, Tuple
from typing import List, Tuple


Interval = Tuple[int, int]


def getIntersect(
    interval1: Tuple[int, int], interval2: Tuple[int, int]
) -> Optional[Tuple[int, int]]:
    """获取两个区间的交集"""
    if interval1[0] > interval2[1] or interval2[0] > interval1[1]:
        return None
    return (max(interval1[0], interval2[0]), min(interval1[1], interval2[1]))


def solve(intervals1: List[Interval], intervals2: List[Interval]) -> List[Interval]:
    """双指针求两个已排序的区间列表的交集"""
    n1, n2 = len(intervals1), len(intervals2)
    res = []
    left, right = 0, 0
    while left < n1 and right < n2:
        s1, e1, s2, e2 = *intervals1[left], *intervals2[right]
        # !相交
        if s1 <= e2 <= e1 or s2 <= e1 <= e2:
            # !尽量往内缩
            res.append((max(s1, s2), min(e1, e2)))
        if e1 < e2:
            left += 1
        else:
            right += 1
    return res


def solve2(intervals1: List[Interval], intervals2: List[Interval]) -> int:
    """两个已排序的区间列表相交的长度之和"""
    n1, n2 = len(intervals1), len(intervals2)
    res = 0
    left, right = 0, 0
    while left < n1 and right < n2:
        s1, e1, s2, e2 = *intervals1[left], *intervals2[right]
        if s1 <= e2 <= e1 or s2 <= e1 <= e2:  # 相交
            res += min(e1, e2) - max(s1, s2) + 1  # [1,1] 区间长度为1
        if e1 < e2:
            left += 1
        else:
            right += 1
    return res
