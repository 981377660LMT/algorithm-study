from itertools import accumulate
from typing import Callable, List


def maxRight(left: int, check: Callable[[int], bool], upper: int) -> int:
    """返回最大的 right 使得 [left,right) 内的值满足 check. right<=upper."""
    ok, ng = left, upper + 1
    while ok + 1 < ng:
        mid = (ok + ng) >> 1
        if check(mid):
            ok = mid
        else:
            ng = mid
    return ok


def minLeft(right: int, check: Callable[[int], bool], lower: int) -> int:
    """返回最小的 left 使得 [left,right) 内的值满足 check. left>=lower."""
    ok, ng = right, lower - 1
    while ng + 1 < ok:
        mid = (ok + ng) >> 1
        if check(mid):
            ok = mid
        else:
            ng = mid
    return ok


if __name__ == "__main__":
    INF = int(1e18)

    def min2(a: int, b: int) -> int:
        return a if a < b else b

    def circularPresum(nums: List[int]) -> Callable[[int, int], int]:
        """环形数组前缀和."""
        n = len(nums)
        preSum = list(accumulate(nums, initial=0))

        def _cal(r: int) -> int:
            return preSum[n] * (r // n) + preSum[r % n]

        def query(start: int, end: int) -> int:
            """[start,end)的和.
            0 <= start < end <= n.
            """
            if start >= end:
                return 0
            return _cal(end) - _cal(start)

        return query

    # https://leetcode.cn/problems/minimum-size-subarray-in-infinite-array/description/
    # 2875. 无限数组的最短子数组
    class Solution:
        def minSizeSubarray(self, nums: List[int], target: int) -> int:
            sum = circularPresum(nums)
            res = INF
            for start in range(len(nums)):
                end = maxRight(start, lambda r: sum(start, r) <= target, int(1e9 + 10))
                if sum(start, end) == target:
                    res = min2(res, end - start)
            return res if res != INF else -1
