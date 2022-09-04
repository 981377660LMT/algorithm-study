from typing import Optional, Tuple


def getIntersect(
    interval1: Tuple[int, int], interval2: Tuple[int, int]
) -> Optional[Tuple[int, int]]:
    """获取两个区间的交集"""
    if interval1[0] > interval2[1] or interval2[0] > interval1[1]:
        return None
    return (max(interval1[0], interval2[0]), min(interval1[1], interval2[1]))
