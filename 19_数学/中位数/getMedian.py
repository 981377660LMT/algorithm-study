from typing import List, Optional


def getMedian(sortedNums: List[int], start: Optional[int] = None, end: Optional[int] = None) -> int:
    """有序数组的中位数(向下取整)."""
    if start is None:
        start = 0
    if end is None:
        end = len(sortedNums)
    if start < 0:
        start = 0
    if end > len(sortedNums):
        end = len(sortedNums)
    if start >= end:
        return 0
    if (end - start) & 1 == 0:
        return (sortedNums[(end + start) // 2 - 1] + sortedNums[(end + start) // 2]) // 2
    return sortedNums[(end + start) // 2]
