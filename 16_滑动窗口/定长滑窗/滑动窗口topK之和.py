# !滑动窗口topK之和 (这里topK是最小值)

from typing import List
from sortedcontainers import SortedList


def windowTopKSum(nums: List[int], windowSize: int, k: int) -> List[int]:
    def add(x: int) -> None:
        nonlocal topKSum
        pos = sl.bisect_left(x)
        if pos < k:
            topKSum += x
            if k - 1 < len(sl):
                topKSum -= sl[k - 1]  # type: ignore
        sl.add(x)

    def remove(x: int) -> None:
        nonlocal topKSum
        pos = sl.bisect_left(x)
        if pos < k:
            topKSum -= x
            if k < len(sl):
                topKSum += sl[k]  # type: ignore
        sl.remove(x)

    def query() -> int:
        return topKSum

    n = len(nums)
    sl = SortedList()
    res, topKSum = [], 0
    for right in range(n):
        add(nums[right])
        if right >= windowSize:
            remove(nums[right - windowSize])
        if right >= windowSize - 1:
            res.append(query())
    return res


assert windowTopKSum([3, 1, 4, 1, 5, 9], 4, 3) == [5, 6, 10]
assert windowTopKSum([12, 2, 17, 11, 19, 8, 4, 3, 6, 20], 6, 3) == [21, 14, 15, 13, 13]
