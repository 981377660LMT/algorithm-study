# 静态区间频率查询

from bisect import bisect_left
from collections import defaultdict
from typing import Callable, List, TypeVar

V = TypeVar("V")


def rangeFreq(arr: List[V]) -> Callable[[int, int, V], int]:
    """返回一个函数,用于查询arr[start:end]内value的频率."""
    mp = defaultdict(list)
    for i, v in enumerate(arr):
        mp[v].append(i)

    def query(start: int, end: int, value: V) -> int:
        if start < 0:
            start = 0
        if end > len(arr):
            end = len(arr)
        if start >= end:
            return 0
        if value not in mp:
            return 0
        else:
            return bisect_left(mp[value], end) - bisect_left(mp[value], start)

    return query


if __name__ == "__main__":
    # https://leetcode.cn/problems/number-of-divisible-triplet-sums/description/
    # 1 <= nums.length <= 1000
    # 1 <= nums[i] <= 109
    # 1 <= d <= 109
    # !固定左端点，剩下的问题就是两数之和，O(n^2)

    from bisect import bisect_left
    from collections import defaultdict
    from typing import Callable, List, TypeVar

    class Solution:
        def divisibleTripletCount(self, nums: List[int], d: int) -> int:
            mods = [v % d for v in nums]
            Q = rangeFreq(mods)
            res = 0
            n = len(nums)
            for i in range(n - 2):
                for j in range(i + 1, n - 1):
                    res += Q(j + 1, n, (-mods[i] - mods[j]) % d)
            return res
