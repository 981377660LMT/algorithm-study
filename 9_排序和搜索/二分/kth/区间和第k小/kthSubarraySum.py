from itertools import accumulate
from typing import List


def kthSubarraySum(nums: List[int], k: int) -> int:
    """第k小的子数组和.k从1开始."""
    k -= 1
    n = len(nums)
    preSum = [0] + list(accumulate(nums))

    def check(target: int) -> bool:
        """countNGT:和小于等于target的子数组个数<=k."""
        res, left, curSum = 0, 0, 0
        for right in range(n):
            curSum += nums[right]
            while left <= right and curSum > target:
                curSum -= nums[left]
                left += 1
            res += right - left + 1
        return res <= k

    left, right = 0, preSum[-1]
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1
    return left
